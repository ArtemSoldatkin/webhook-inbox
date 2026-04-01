package service

import (
	"context"
	"testing"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestListSources(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	ts := testNow()
	cursorID := int64(77)
	cursor := types.NewCursor(&ts, &cursorID)
	dbtx.queryHandlers["-- name: ListSources :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 6)
		assert.Equal(t, ts, args[0].(pgtype.Timestamptz).Time)
		assert.Equal(t, "DESC", args[1])
		assert.Equal(t, int64(77), args[2])
		assert.Equal(t, "billing", args[3])
		assert.Equal(t, "active", args[4])
		assert.Equal(t, int32(20), args[5])
		return &serviceTestRows{rows: [][]any{
			serviceSourceRowValues(db.Source{ID: 1, PublicID: mustServiceUUID(t, "11111111-1111-1111-1111-111111111111"), EgressUrl: "https://example.com", StaticHeaders: []byte(`{}`), Status: "active"}),
		}}, nil
	}

	sources, err := svc.ListSources(context.Background(), cursor, 20, "billing", "active", "DESC")

	require.NoError(t, err)
	require.Len(t, sources, 1)
	assert.Equal(t, int64(1), sources[0].ID)
}

func TestGetSourceByID_UsesCache(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	source := db.Source{
		ID:            5,
		PublicID:      mustServiceUUID(t, "11111111-1111-1111-1111-111111111111"),
		EgressUrl:     "https://example.com",
		StaticHeaders: []byte(`{}`),
		Status:        "active",
	}
	calls := 0
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		calls++
		return &serviceTestRow{values: serviceSourceRowValues(source)}
	}

	first, err := svc.GetSourceByID(context.Background(), 5)
	require.NoError(t, err)
	svc.Cache.Wait()
	second, err := svc.GetSourceByID(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, source, first)
	assert.Equal(t, source, second)
	assert.Equal(t, 1, calls)
}

func TestGetSourceByPublicID_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc := newServiceUnderTest(t, newServiceTestDB())

	source, err := svc.GetSourceByPublicID(context.Background(), "bad-uuid")

	assert.Equal(t, db.Source{}, source)
	require.Error(t, err)
}

func TestCreateSource(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	dbtx.queryRowHandlers["-- name: CreateSource :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 3)
		assert.Equal(t, "https://example.com/webhook", args[0])
		assert.JSONEq(t, `{"Authorization":"Bearer token"}`, string(args[1].([]byte)))
		assert.Equal(t, "billing events", args[2].(pgtype.Text).String)
		return &serviceTestRow{values: serviceSourceRowValues(db.Source{
			ID:            7,
			PublicID:      mustServiceUUID(t, "22222222-2222-2222-2222-222222222222"),
			EgressUrl:     "https://example.com/webhook",
			StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
			Status:        "active",
			Description:   pgtype.Text{String: "billing events", Valid: true},
		})}
	}

	source, err := svc.CreateSource(context.Background(), CreateSourceInput{
		EgressUrl:     "https://example.com/webhook",
		StaticHeaders: map[string]string{"Authorization": "Bearer token"},
		Description:   "billing events",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(7), source.ID)
}

func TestCreateSource_InvalidEgressURL(t *testing.T) {
	t.Parallel()

	svc := newServiceUnderTest(t, newServiceTestDB())

	source, err := svc.CreateSource(context.Background(), CreateSourceInput{
		EgressUrl: "ftp://example.com",
	})

	assert.Equal(t, db.Source{}, source)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid egress URL")
}

func TestIsValidStatusTransition(t *testing.T) {
	t.Parallel()

	assert.True(t, isValidStatusTransition("paused", "active"))
	assert.False(t, isValidStatusTransition("active", "disabled"))
	assert.False(t, isValidStatusTransition("unknown", "active"))
}

func TestUpdateSourceStatus_NoChange(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	source := db.Source{
		ID:            8,
		PublicID:      mustServiceUUID(t, "33333333-3333-3333-3333-333333333333"),
		EgressUrl:     "https://example.com",
		StaticHeaders: []byte(`{}`),
		Status:        "active",
	}
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &serviceTestRow{values: serviceSourceRowValues(source)}
	}

	updated, err := svc.UpdateSourceStatus(context.Background(), UpdateSourceStatusInput{
		SourceID: 8,
		Status:   "active",
	})

	require.NoError(t, err)
	assert.Equal(t, source, updated)
}

func TestUpdateSourceStatus_InvalidTransition(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	source := db.Source{
		ID:            8,
		PublicID:      mustServiceUUID(t, "33333333-3333-3333-3333-333333333333"),
		EgressUrl:     "https://example.com",
		StaticHeaders: []byte(`{}`),
		Status:        "active",
	}
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &serviceTestRow{values: serviceSourceRowValues(source)}
	}

	updated, err := svc.UpdateSourceStatus(context.Background(), UpdateSourceStatusInput{
		SourceID: 8,
		Status:   "disabled",
	})

	assert.Equal(t, source, updated)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status transition")
}

func TestUpdateSourceStatus_UpdatesAndRefetches(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	original := db.Source{
		ID:            9,
		PublicID:      mustServiceUUID(t, "44444444-4444-4444-4444-444444444444"),
		EgressUrl:     "https://example.com",
		StaticHeaders: []byte(`{}`),
		Status:        "paused",
	}
	updatedSource := original
	updatedSource.Status = "active"
	getCalls := 0
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		getCalls++
		if getCalls == 1 {
			return &serviceTestRow{values: serviceSourceRowValues(original)}
		}
		return &serviceTestRow{values: serviceSourceRowValues(updatedSource)}
	}
	dbtx.execHandlers["-- name: UpdateSourceStatus :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		require.Len(t, args, 3)
		assert.Equal(t, "active", args[0])
		assert.Equal(t, "resume", args[1].(pgtype.Text).String)
		assert.Equal(t, int64(9), args[2])
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	source, err := svc.UpdateSourceStatus(context.Background(), UpdateSourceStatusInput{
		SourceID:     9,
		Status:       "active",
		StatusReason: "resume",
	})

	require.NoError(t, err)
	assert.Equal(t, updatedSource.Status, source.Status)
	assert.Equal(t, 2, getCalls)
}

func TestValidateEgressURL(t *testing.T) {
	t.Parallel()

	assert.True(t, validateEgressUrl("https://example.com/webhook", "dev"))
	assert.False(t, validateEgressUrl("", "dev"))
	assert.False(t, validateEgressUrl("ftp://example.com", "dev"))
	assert.False(t, validateEgressUrl("http://localhost/webhook", "prod"))
}
