package routev1

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestListEvents_ReturnsPaginatedResponse(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	receivedAt := testTime()
	remoteAddress := netip.MustParseAddr("192.0.2.10")
	dbtx.queryHandlers["-- name: ListEventsBySource :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 6)
		assert.Equal(t, int64(21), args[0])
		assert.Equal(t, "ASC", args[2])
		assert.Equal(t, int32(1), args[5])
		return &testRows{
			rows: [][]any{
				{int64(1), int64(21), pgtype.Text{String: "dedup-1", Valid: true}, "POST", "/ingest/source-1", &remoteAddress, []byte(`{"q":["1"]}`), []byte(`{"H":["v"]}`), []byte(`payload-1`), "text/plain", pgtype.Timestamptz{Time: receivedAt, Valid: true}},
				{int64(2), int64(21), pgtype.Text{String: "dedup-2", Valid: true}, "POST", "/ingest/source-1", &remoteAddress, []byte(`{"q":["2"]}`), []byte(`{"H":["w"]}`), []byte(`payload-2`), "text/plain", pgtype.Timestamptz{Time: receivedAt.Add(time.Minute), Valid: true}},
			},
		}, nil
	}

	req := withURLParam(httptest.NewRequest("GET", "/sources/21/events?limit=1&sort_direction=ASC", nil), "source_id", "21")
	recorder := httptest.NewRecorder()

	listEvents(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[paginatedEventsResponse](t, recorder)
	require.Len(t, response.Data, 1)
	assert.True(t, response.HasNext)
	require.NotNil(t, response.NextCursor)
	assert.Equal(t, int64(1), response.Data[0].ID)
	assert.Equal(t, "192.0.2.10", *response.Data[0].RemoteAddress)
}

func TestGetEvent_NotFoundReturns404(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	dbtx.queryRowHandlers["-- name: GetEventByID :one"] = func(args ...any) pgx.Row {
		assert.Equal(t, int64(42), args[0])
		return &testRow{err: pgx.ErrNoRows}
	}

	req := withURLParam(httptest.NewRequest("GET", "/events/42", nil), "event_id", "42")
	recorder := httptest.NewRecorder()

	getEvent(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	response := decodeJSONResponse[errorResponse](t, recorder)
	assert.Equal(t, "Event not found", response.Error)
}
