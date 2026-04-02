package routev1

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestListSources_ReturnsPaginatedResponse(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	updatedAt := testTime()
	dbtx.queryHandlers["-- name: ListSources :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 6)
		assert.Equal(t, "DESC", args[1])
		assert.Equal(t, "*", args[4])
		assert.Equal(t, int32(1), args[5])
		return &testRows{
			rows: [][]any{
				sourceRowValues(t, db.Source{
					ID:            1,
					PublicID:      mustPGUUID(t, "11111111-1111-1111-1111-111111111111"),
					EgressUrl:     "https://consumer.example.com/one",
					StaticHeaders: []byte(`{"Authorization":"Bearer one"}`),
					Status:        "active",
					CreatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(-time.Hour), Valid: true},
					UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
				}),
				sourceRowValues(t, db.Source{
					ID:            2,
					PublicID:      mustPGUUID(t, "22222222-2222-2222-2222-222222222222"),
					EgressUrl:     "https://consumer.example.com/two",
					StaticHeaders: []byte(`{"Authorization":"Bearer two"}`),
					Status:        "paused",
					CreatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(-2 * time.Hour), Valid: true},
					UpdatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(time.Minute), Valid: true},
				}),
			},
		}, nil
	}

	req := httptest.NewRequest("GET", "/sources?limit=1", nil)
	recorder := httptest.NewRecorder()

	listSources(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[paginatedSourcesResponse](t, recorder)
	require.Len(t, response.Data, 1)
	assert.True(t, response.HasNext)
	require.NotNil(t, response.NextCursor)
	assert.Equal(t, int64(1), response.Data[0].ID)
	assert.Equal(t, "https://api.example.com:8443/ingest/11111111-1111-1111-1111-111111111111", response.Data[0].IngressUrl)
}

func TestGetSourceByID_NotFoundReturns404(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &testRow{err: pgx.ErrNoRows}
	}

	req := withURLParam(httptest.NewRequest("GET", "/sources/7", nil), "source_id", "7")
	recorder := httptest.NewRecorder()

	getSourceByID(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	response := decodeJSONResponse[errorResponse](t, recorder)
	assert.Equal(t, "Source not found", response.Error)
}

func TestCreateSource_ValidationErrorReturns400(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		"POST",
		"/sources",
		bytes.NewBufferString(`{"EgressUrl":"https://example.com/webhook","Description":"`+string(bytes.Repeat([]byte("a"), 501))+`"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	createSource(newTestService(t, newTestDB(), newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	response := decodeJSONResponse[errorResponse](t, recorder)
	assert.Contains(t, response.Error, "Invalid input parameters")
}

func TestCreateSource_ReturnsCreatedSource(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	createdAt := testTime()
	dbtx.queryRowHandlers["-- name: CreateSource :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 3)
		assert.Equal(t, "https://example.com/webhook", args[0])
		return &testRow{values: sourceRowValues(t, db.Source{
			ID:            55,
			PublicID:      mustPGUUID(t, "33333333-3333-3333-3333-333333333333"),
			EgressUrl:     "https://example.com/webhook",
			StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
			Status:        "active",
			Description:   pgtype.Text{String: "billing events", Valid: true},
			CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
			UpdatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
		})}
	}

	req := httptest.NewRequest(
		"POST",
		"/sources",
		bytes.NewBufferString(`{"EgressUrl":"https://example.com/webhook","StaticHeaders":{"Authorization":"Bearer token"},"Description":"billing events"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	createSource(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusCreated, recorder.Code)
	response := decodeJSONResponse[map[string]any](t, recorder)
	assert.Equal(t, float64(55), response["id"])
	assert.Equal(t, "33333333-3333-3333-3333-333333333333", response["public_id"])
}

func TestUpdateSourceStatus_ReturnsUpdatedSource(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	updatedAt := testTime()
	getCalls := 0
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		getCalls++
		switch getCalls {
		case 1:
			return &testRow{values: sourceRowValues(t, db.Source{
				ID:            44,
				PublicID:      mustPGUUID(t, "44444444-4444-4444-4444-444444444444"),
				EgressUrl:     "https://example.com/webhook",
				StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
				Status:        "paused",
				CreatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(-time.Hour), Valid: true},
				UpdatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(-time.Hour), Valid: true},
			})}
		default:
			return &testRow{values: sourceRowValues(t, db.Source{
				ID:            44,
				PublicID:      mustPGUUID(t, "44444444-4444-4444-4444-444444444444"),
				EgressUrl:     "https://example.com/webhook",
				StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
				Status:        "active",
				StatusReason:  pgtype.Text{String: "resume delivery", Valid: true},
				CreatedAt:     pgtype.Timestamptz{Time: updatedAt.Add(-time.Hour), Valid: true},
				UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
			})}
		}
	}
	dbtx.execHandlers["-- name: UpdateSourceStatus :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		require.Len(t, args, 3)
		assert.Equal(t, "active", args[0])
		assert.Equal(t, int64(44), args[2])
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	req := withURLParam(httptest.NewRequest("PUT", "/sources/44/status", bytes.NewBufferString(`{"status":"active","status_reason":"resume delivery"}`)), "source_id", "44")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	updateSourceStatus(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[map[string]any](t, recorder)
	assert.Equal(t, float64(44), response["id"])
	assert.Equal(t, "active", response["status"])
}

func TestUpdateSourceStatus_ServiceErrorReturns500(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	dbtx.queryRowHandlers["-- name: GetSourceByID :one"] = func(args ...any) pgx.Row {
		return &testRow{values: sourceRowValues(t, db.Source{
			ID:            44,
			PublicID:      mustPGUUID(t, "44444444-4444-4444-4444-444444444444"),
			EgressUrl:     "https://example.com/webhook",
			StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
			Status:        "paused",
			CreatedAt:     pgtype.Timestamptz{Time: testTime().Add(-time.Hour), Valid: true},
			UpdatedAt:     pgtype.Timestamptz{Time: testTime().Add(-time.Hour), Valid: true},
		})}
	}
	dbtx.execHandlers["-- name: UpdateSourceStatus :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		return pgconn.NewCommandTag(""), errors.New("update failed")
	}

	req := withURLParam(httptest.NewRequest("PUT", "/sources/44/status", bytes.NewBufferString(`{"status":"active","status_reason":"resume delivery"}`)), "source_id", "44")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	updateSourceStatus(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	response := decodeJSONResponse[errorResponse](t, recorder)
	assert.Equal(t, "Failed to update source status", response.Error)
}

func mustPGUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()
	var id pgtype.UUID
	require.NoError(t, id.Scan(value))
	return id
}

func sourceRowValues(t *testing.T, source db.Source) []any {
	t.Helper()
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
