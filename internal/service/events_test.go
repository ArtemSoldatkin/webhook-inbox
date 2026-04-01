package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestListEvents(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	ts := testNow()
	cursorID := int64(33)
	cursor := types.NewCursor(&ts, &cursorID)
	dbtx.queryHandlers["-- name: ListEventsBySource :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 6)
		assert.Equal(t, int64(12), args[0])
		assert.Equal(t, ts, args[1].(pgtype.Timestamptz).Time)
		assert.Equal(t, "DESC", args[2])
		assert.Equal(t, int64(33), args[3])
		assert.Equal(t, "billing", args[4])
		assert.Equal(t, int32(5), args[5])
		return &serviceTestRows{rows: [][]any{
			serviceEventRowValues(db.Event{ID: 1, SourceID: 12, Method: http.MethodPost}),
		}}, nil
	}

	events, err := svc.ListEvents(context.Background(), 12, cursor, 5, "billing", "DESC")

	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, int64(1), events[0].ID)
}

func TestGetEventByID_UsesCache(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	event := db.Event{ID: 7, SourceID: 3, Method: http.MethodPost}
	calls := 0
	dbtx.queryRowHandlers["-- name: GetEventByID :one"] = func(args ...any) pgx.Row {
		calls++
		return &serviceTestRow{values: serviceEventRowValues(event)}
	}

	first, err := svc.GetEventByID(context.Background(), 7)
	require.NoError(t, err)
	svc.Cache.Wait()
	second, err := svc.GetEventByID(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, event, first)
	assert.Equal(t, event, second)
	assert.Equal(t, 1, calls)
}

func TestCreateEvent_SourceNotFound(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	dbtx.queryRowHandlers["-- name: GetSourceByPublicID :one"] = func(args ...any) pgx.Row {
		return &serviceTestRow{err: pgx.ErrNoRows}
	}

	req := httptest.NewRequest(http.MethodPost, "/ingest/source?foo=bar", nil)
	req.RemoteAddr = "192.0.2.10:1234"

	eventID, err := svc.CreateEvent(context.Background(), req, "11111111-1111-1111-1111-111111111111")

	assert.Equal(t, int64(0), eventID)
	var notFoundErr *SourceIsNotFound
	require.Error(t, err)
	require.True(t, errors.As(err, &notFoundErr))
	assert.Contains(t, notFoundErr.Error(), "source with public_id")
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCreateEvent(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	dbtx.queryRowHandlers["-- name: GetSourceByPublicID :one"] = func(args ...any) pgx.Row {
		return &serviceTestRow{values: serviceSourceRowValues(db.Source{
			ID:            9,
			PublicID:      mustServiceUUID(t, "11111111-1111-1111-1111-111111111111"),
			EgressUrl:     "https://example.com",
			StaticHeaders: []byte(`{}`),
			Status:        "active",
		})}
	}
	dbtx.queryRowHandlers["-- name: CreateEvent :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 9)
		assert.Equal(t, int64(9), args[0])
		assert.Equal(t, http.MethodPost, args[2])
		assert.Equal(t, "/ingest/source", args[3])
		assert.JSONEq(t, `{"foo":["bar"]}`, string(args[5].([]byte)))
		assert.JSONEq(t, `{"Content-Type":["application/json"],"X-Test":["value"]}`, string(args[6].([]byte)))
		assert.Equal(t, []byte(`payload`), args[7].([]byte))
		assert.Equal(t, "application/json", args[8])
		return &serviceTestRow{values: []any{int64(101)}}
	}

	req := httptest.NewRequest(http.MethodPost, "/ingest/source?foo=bar", strings.NewReader(`payload`))
	req.RemoteAddr = "192.0.2.10:1234"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test", "value")

	eventID, err := svc.CreateEvent(context.Background(), req, "11111111-1111-1111-1111-111111111111")

	require.NoError(t, err)
	assert.Equal(t, int64(101), eventID)
	body, readErr := io.ReadAll(req.Body)
	require.NoError(t, readErr)
	assert.Equal(t, "payload", string(body))
}

func TestSourceIsNotFoundUnwrap(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base")
	err := &SourceIsNotFound{Message: "missing", Err: baseErr}

	assert.Equal(t, "missing", err.Error())
	assert.ErrorIs(t, err, baseErr)
}

func TestGenerateDedupHash(t *testing.T) {
	t.Parallel()

	payload := DedupPayload{
		Method:      http.MethodPost,
		IngressPath: "/ingest/source",
		QueryParams: []byte(`{"foo":["bar"]}`),
		RawHeaders:  []byte(`{"X-Test":["value"]}`),
		Body:        []byte(`payload`),
	}

	hashA, err := generateDedupHash(payload)
	require.NoError(t, err)
	hashB, err := generateDedupHash(payload)

	require.NoError(t, err)
	assert.Equal(t, hashA, hashB)
	assert.Len(t, hashA, 64)
}
