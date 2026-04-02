package routev1

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestListDeliveryAttempts_InvalidInputReturnsBadRequest(t *testing.T) {
	t.Parallel()

	req := withURLParam(httptest.NewRequest("GET", "/events/x/delivery-attempts", nil), "event_id", "bad")
	recorder := httptest.NewRecorder()

	listDeliveryAttempts(newTestService(t, newTestDB(), newTestConfig())).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	response := decodeJSONResponse[errorResponse](t, recorder)
	assert.Equal(t, "Invalid input parameters", response.Error)
}

func TestListDeliveryAttempts_ReturnsPaginatedResponse(t *testing.T) {
	t.Parallel()

	dbtx := newTestDB()
	createdAt := testTime()
	dbtx.queryHandlers["-- name: ListDeliveryAttemptsByEvent :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 7)
		assert.Equal(t, int64(15), args[0])
		assert.Equal(t, "failed", args[5])
		assert.Equal(t, "ASC", args[2])
		assert.Equal(t, int32(1), args[6])
		return &testRows{
			rows: [][]any{
				{int64(100), int64(15), int32(1), "failed", pgtype.Int4{Int32: 500, Valid: true}, pgtype.Text{}, pgtype.Text{}, pgtype.Timestamptz{}, pgtype.Timestamptz{}, pgtype.Timestamptz{Time: createdAt, Valid: true}, pgtype.Timestamptz{}},
				{int64(101), int64(15), int32(2), "failed", pgtype.Int4{Int32: 502, Valid: true}, pgtype.Text{}, pgtype.Text{}, pgtype.Timestamptz{}, pgtype.Timestamptz{}, pgtype.Timestamptz{Time: createdAt.Add(time.Minute), Valid: true}, pgtype.Timestamptz{}},
			},
		}, nil
	}

	req := withURLParam(httptest.NewRequest("GET", "/events/15/delivery-attempts?limit=1&filter_state=failed&sort_direction=ASC", nil), "event_id", "15")
	recorder := httptest.NewRecorder()

	listDeliveryAttempts(newTestService(t, dbtx, newTestConfig())).ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[paginatedDeliveryAttemptsResponse](t, recorder)
	require.Len(t, response.Data, 1)
	assert.True(t, response.HasNext)
	require.NotNil(t, response.NextCursor)
	assert.NotEmpty(t, *response.NextCursor)
	assert.Equal(t, int64(100), response.Data[0].ID)
}
