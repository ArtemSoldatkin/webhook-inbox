package deliveryengine

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/dgraph-io/ristretto"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestInterpretDeliveryResponse(t *testing.T) {
	t.Parallel()

	success := interpretDeliveryResponse(&http.Response{StatusCode: http.StatusNoContent})
	assert.Equal(t, "succeeded", success.DeliveryState)
	assert.Equal(t, "", success.ErrorType)
	assert.Equal(t, "", success.ErrorMessage)

	clientError := interpretDeliveryResponse(&http.Response{StatusCode: http.StatusNotFound})
	assert.Equal(t, "failed", clientError.DeliveryState)
	assert.Equal(t, "http_4xx", clientError.ErrorType)
	assert.Equal(t, "Not Found", clientError.ErrorMessage)

	serverError := interpretDeliveryResponse(&http.Response{StatusCode: http.StatusBadGateway})
	assert.Equal(t, "failed", serverError.DeliveryState)
	assert.Equal(t, "http_5xx", serverError.ErrorType)
	assert.Equal(t, "Bad Gateway", serverError.ErrorMessage)
}

func TestSendDeliveryRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/deliver", r.URL.Path)
		assert.ElementsMatch(t, []string{"bar", "baz"}, r.URL.Query()["foo"])
		assert.Equal(t, "token", r.Header.Get("X-Token"))
		assert.Equal(t, "payload", mustReadBody(t, r))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	res, err := sendDeliveryRequest(context.Background(), server.Client(), &DeliveryPayload{
		URL:    server.URL + "/deliver",
		Method: http.MethodPost,
		Headers: map[string][]string{
			"X-Token": {"token"},
		},
		QueryParams: map[string][]string{
			"foo": {"bar", "baz"},
		},
		Body: []byte("payload"),
	})

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	require.NoError(t, res.Body.Close())
}

func TestMarkInPending(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	dbtx.execHandlers["-- name: UpdateDeliveryAttempt :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		require.Len(t, args, 7)
		assert.Equal(t, "pending", args[0])
		assert.False(t, args[1].(pgtype.Int4).Valid)
		assert.False(t, args[2].(pgtype.Text).Valid)
		assert.False(t, args[3].(pgtype.Text).Valid)
		assert.False(t, args[4].(pgtype.Timestamptz).Valid)
		assert.False(t, args[5].(pgtype.Timestamptz).Valid)
		assert.Equal(t, int64(99), args[6])
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	err := markInPending(context.Background(), newEngineTestService(t, dbtx), 99)

	require.NoError(t, err)
}

func TestFinalizeDeliveryAttempt(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	dbtx.execHandlers["-- name: UpdateDeliveryAttempt :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		require.Len(t, args, 7)
		assert.Equal(t, "failed", args[0])
		assert.Equal(t, int32(502), args[1].(pgtype.Int4).Int32)
		assert.True(t, args[1].(pgtype.Int4).Valid)
		assert.Equal(t, "http_5xx", args[2].(pgtype.Text).String)
		assert.Equal(t, "Bad Gateway", args[3].(pgtype.Text).String)
		assert.True(t, args[5].(pgtype.Timestamptz).Valid)
		assert.Equal(t, int64(55), args[6])
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	err := finalizeDeliveryAttempt(context.Background(), newEngineTestService(t, dbtx), 55, &DeliveryResult{
		StatusCode:    502,
		DeliveryState: "failed",
		ErrorType:     "http_5xx",
		ErrorMessage:  "Bad Gateway",
	})

	require.NoError(t, err)
}

func TestScheduleRetry(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	svc.Config.APIDeliveryRetryBackoffBaseSec = 2
	svc.Config.APIDeliveryRetryBackoffMaxSec = 60

	before := time.Now()
	dbtx.queryRowHandlers["-- name: CreateDeliveryAttempt :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 9)
		assert.Equal(t, int64(42), args[0])
		assert.Equal(t, int32(3), args[1])
		assert.Equal(t, "pending", args[2])
		nextAttempt := args[8].(pgtype.Timestamptz)
		assert.True(t, nextAttempt.Valid)
		assert.WithinDuration(t, before.Add(10*time.Second), nextAttempt.Time, 2*time.Second)
		return &engineTestRow{values: []any{int64(1234)}}
	}

	id, err := scheduleRetry(context.Background(), svc, service.PendingDeliveryAttempt{
		ID:            10,
		EventID:       42,
		AttemptNumber: 2,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1234), id)
}

func TestLoadDeliveryPayload(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	dbtx.queryRowHandlers["-- name: GetEventByID :one"] = func(args ...any) pgx.Row {
		return &engineTestRow{values: eventRowValues(db.Event{
			ID:              11,
			SourceID:        22,
			Method:          http.MethodPost,
			RawHeaders:      []byte(`{"X-Trace":["abc"],"Authorization":["event-token"]}`),
			QueryParams:     []byte(`{"foo":["bar"]}`),
			Body:            []byte(`payload`),
			BodyContentType: "application/json",
		})}
	}
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &engineTestRow{values: sourceRowValues(db.Source{
			ID:            22,
			PublicID:      mustPGUUID(t, "11111111-1111-1111-1111-111111111111"),
			EgressUrl:     "https://example.com/deliver",
			StaticHeaders: []byte(`{"Authorization":"source-token","X-Static":"yes"}`),
			Status:        "active",
		})}
	}

	payload, err := loadDeliveryPayload(context.Background(), svc, service.PendingDeliveryAttempt{
		ID:      1,
		EventID: 11,
	})

	require.NoError(t, err)
	require.NotNil(t, payload)
	assert.Equal(t, "https://example.com/deliver", payload.URL)
	assert.Equal(t, http.MethodPost, payload.Method)
	assert.Equal(t, []byte("payload"), payload.Body)
	assert.Equal(t, map[string][]string{"foo": {"bar"}}, payload.QueryParams)
	assert.Equal(t, []string{"source-token", "event-token"}, payload.Headers["Authorization"])
	assert.Equal(t, []string{"yes"}, payload.Headers["X-Static"])
	assert.Equal(t, []string{"abc"}, payload.Headers["X-Trace"])
}

func TestLoadDeliveryPayload_InvalidStaticHeaders(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	dbtx.queryRowHandlers["-- name: GetEventByID :one"] = func(args ...any) pgx.Row {
		return &engineTestRow{values: eventRowValues(db.Event{
			ID:          11,
			SourceID:    22,
			Method:      http.MethodPost,
			RawHeaders:  []byte(`{"X-Test":["value"]}`),
			QueryParams: []byte(`{"foo":["bar"]}`),
		})}
	}
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &engineTestRow{values: sourceRowValues(db.Source{
			ID:            22,
			PublicID:      mustPGUUID(t, "11111111-1111-1111-1111-111111111111"),
			EgressUrl:     "https://example.com/deliver",
			StaticHeaders: []byte(`{`),
			Status:        "active",
		})}
	}

	payload, err := loadDeliveryPayload(context.Background(), svc, service.PendingDeliveryAttempt{EventID: 11})

	assert.Nil(t, payload)
	require.Error(t, err)
}

func TestHandleDeliveryFinalizationAndRetry_SchedulesRetry(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	svc.Config.APIDeliveryMaxRetries = 3
	svc.Config.APIDeliveryRetryBackoffBaseSec = 1
	svc.Config.APIDeliveryRetryBackoffMaxSec = 60
	finalized := false
	scheduled := false
	dbtx.execHandlers["-- name: UpdateDeliveryAttempt :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		finalized = true
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}
	dbtx.queryRowHandlers["-- name: CreateDeliveryAttempt :one"] = func(args ...any) pgx.Row {
		scheduled = true
		return &engineTestRow{values: []any{int64(999)}}
	}

	handleDeliveryFinalizationAndRetry(context.Background(), svc, service.PendingDeliveryAttempt{
		ID:            5,
		EventID:       6,
		AttemptNumber: 1,
	}, &DeliveryResult{
		StatusCode:    500,
		DeliveryState: "failed",
		ErrorType:     "http_5xx",
		ErrorMessage:  "Internal Server Error",
	})

	assert.True(t, finalized)
	assert.True(t, scheduled)
}

func TestHandleDeliveryFinalizationAndRetry_DoesNotScheduleRetryForSuccess(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	svc.Config.APIDeliveryMaxRetries = 3
	finalized := false
	dbtx.execHandlers["-- name: UpdateDeliveryAttempt :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		finalized = true
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}
	dbtx.queryRowHandlers["-- name: CreateDeliveryAttempt :one"] = func(args ...any) pgx.Row {
		t.Fatalf("retry should not be scheduled for successful delivery")
		return nil
	}

	handleDeliveryFinalizationAndRetry(context.Background(), svc, service.PendingDeliveryAttempt{
		ID:            5,
		EventID:       6,
		AttemptNumber: 1,
	}, &DeliveryResult{
		StatusCode:    204,
		DeliveryState: "succeeded",
	})

	assert.True(t, finalized)
}

type engineTestDB struct {
	queryRowHandlers map[string]func(args ...any) pgx.Row
	execHandlers     map[string]func(args ...any) (pgconn.CommandTag, error)
}

func newEngineTestDB() *engineTestDB {
	return &engineTestDB{
		queryRowHandlers: make(map[string]func(args ...any) pgx.Row),
		execHandlers:     make(map[string]func(args ...any) (pgconn.CommandTag, error)),
	}
}

func (dbtx *engineTestDB) Query(_ context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return nil, errors.New("unexpected query: " + engineQuerySignature(query))
}

func (dbtx *engineTestDB) QueryRow(_ context.Context, query string, args ...interface{}) pgx.Row {
	handler, ok := dbtx.queryRowHandlers[engineQuerySignature(query)]
	if !ok {
		return &engineTestRow{err: errors.New("unexpected query row: " + engineQuerySignature(query))}
	}
	return handler(args...)
}

func (dbtx *engineTestDB) Exec(_ context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	handler, ok := dbtx.execHandlers[engineQuerySignature(query)]
	if !ok {
		return pgconn.NewCommandTag(""), errors.New("unexpected exec: " + engineQuerySignature(query))
	}
	return handler(args...)
}

type engineTestRow struct {
	values []any
	err    error
}

func (r *engineTestRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := engineAssignScanValue(dest[i], r.values[i]); err != nil {
			return err
		}
	}
	return nil
}

func engineAssignScanValue(dest any, value any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("scan destination must be a non-nil pointer")
	}
	target := destValue.Elem()
	if value == nil {
		target.SetZero()
		return nil
	}
	source := reflect.ValueOf(value)
	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}
	if source.Type().ConvertibleTo(target.Type()) {
		target.Set(source.Convert(target.Type()))
		return nil
	}
	return errors.New("cannot assign scan value")
}

func newEngineTestService(t *testing.T, dbtx db.DBTX) *service.Service {
	t.Helper()

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e3,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	require.NoError(t, err)
	t.Cleanup(cache.Close)

	svc := &service.Service{
		Config: &config.Config{
			APICacheDefaultCost:            1024,
			APICacheSourceTTLSec:           300,
			APICacheEventTTLSec:            900,
			APIDeliveryRetryBackoffBaseSec: 1,
			APIDeliveryRetryBackoffMaxSec:  60,
			APIDeliveryMaxRetries:          3,
			APIDeliveryTimeoutSec:          15,
			APIDeliveryRequestTimeoutSec:   15,
			APIRecoveryTimeoutSec:          15,
			APIThrottleConcurrentLimit:     1,
			APIDeliveryMaxConcurrency:      1,
			APIRateLimitRequests:           1,
			APIRateLimitWindowSec:          1,
			APIRequestSizeLimitBytes:       1024,
		},
		Cache: cache,
	}
	engineSetUnexportedField(t, svc, "queries", db.New(dbtx))
	return svc
}

func engineSetUnexportedField(t *testing.T, target any, fieldName string, value any) {
	t.Helper()
	v := reflect.ValueOf(target).Elem().FieldByName(fieldName)
	require.True(t, v.IsValid(), "field %s must exist", fieldName)
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}

func engineQuerySignature(query string) string {
	firstLine, _, _ := strings.Cut(query, "\n")
	return firstLine
}

func mustPGUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()
	var id pgtype.UUID
	require.NoError(t, id.Scan(value))
	return id
}

func eventRowValues(event db.Event) []any {
	return []any{
		event.ID,
		event.SourceID,
		event.DedupHash,
		event.Method,
		event.IngressPath,
		event.RemoteAddress,
		event.QueryParams,
		event.RawHeaders,
		event.Body,
		event.BodyContentType,
		event.ReceivedAt,
	}
}

func sourceRowValues(source db.Source) []any {
	return []any{
		source.ID,
		source.PublicID,
		source.EgressUrl,
		source.StaticHeaders,
		source.Status,
		source.StatusReason,
		source.Description,
		source.CreatedAt,
		source.UpdatedAt,
		source.DisableAt,
	}
}

func mustReadBody(t *testing.T, r *http.Request) string {
	t.Helper()
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	return string(body)
}
