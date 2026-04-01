package routev1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
)

func TestIngestEvent_InvalidPublicIDReturnsBadRequest(t *testing.T) {
	t.Parallel()

	req := withURLParam(httptest.NewRequest("POST", "/ingest/not-a-uuid", nil), "public_id", "not-a-uuid")
	recorder := httptest.NewRecorder()

	ingestEvent(newTestService(t, newTestDB(), newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid input parameters: invalid public_id format")
}

func TestIngestEvent_SourceNotFoundReturns404(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	dbtx.queryRowHandlers["-- name: GetSourceByPublicID :one"] = func(args ...any) pgx.Row {
		return &testRow{err: pgx.ErrNoRows}
	}

	req := withURLParam(httptest.NewRequest("POST", "/ingest/123e4567-e89b-12d3-a456-426614174000", nil), "public_id", "123e4567-e89b-12d3-a456-426614174000")
	req.RemoteAddr = "192.0.2.10:1234"
	recorder := httptest.NewRecorder()

	ingestEvent(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "source with public_id '123e4567-e89b-12d3-a456-426614174000' not found")
}

func TestIngestEvent_SuccessReturnsAccepted(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	source := sourceRowValues(t, db.Source{
		ID:            9,
		PublicID:      mustPGUUID(t, "123e4567-e89b-12d3-a456-426614174000"),
		EgressUrl:     "https://consumer.example.com/webhook",
		StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
		Status:        "active",
	})
	dbtx.queryRowHandlers["-- name: GetSourceByPublicID :one"] = func(args ...any) pgx.Row {
		return &testRow{values: source}
	}
	dbtx.queryRowHandlers["-- name: CreateEvent :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 9)
		assert.Equal(t, int64(9), args[0])
		assert.Equal(t, "POST", args[2])
		return &testRow{values: []any{int64(1001)}}
	}
	dbtx.queryRowHandlers["-- name: CreateDeliveryAttempt :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 9)
		assert.Equal(t, int64(1001), args[0])
		assert.Equal(t, int32(1), args[1])
		assert.Equal(t, "pending", args[2])
		return &testRow{values: []any{int64(2002)}}
	}

	req := withURLParam(httptest.NewRequest("POST", "/ingest/123e4567-e89b-12d3-a456-426614174000", nil), "public_id", "123e4567-e89b-12d3-a456-426614174000")
	req.RemoteAddr = "192.0.2.10:1234"
	recorder := httptest.NewRecorder()

	ingestEvent(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusAccepted, recorder.Code)
	assert.Equal(t, "OK", recorder.Body.String())
}
