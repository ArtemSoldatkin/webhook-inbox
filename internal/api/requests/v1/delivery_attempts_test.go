package requestsv1

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestParseListDeliveryAttemptsInput(t *testing.T) {
	t.Parallel()

	cursorTime := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	cursorID := int64(77)
	cursor := newTestCursor(cursorTime, cursorID)

	req := httptest.NewRequest(
		"GET",
		"/events/15/delivery-attempts?search=timeout&filter_state=failed&sort_direction=ASC&limit=25&cursor="+cursor,
		nil,
	)
	req = withURLParam(req, "event_id", "15")

	input, err := api.ParseRequestInput[ListDeliveryAttemptsInput](req)

	require.NoError(t, err)
	require.NotNil(t, input)
	assert.Equal(t, int64(15), input.EventID)
	assert.Equal(t, "timeout", input.Search)
	assert.Equal(t, "failed", input.Filter)
	assert.Equal(t, "ASC", input.SortDirection)
	assert.Equal(t, 25, input.PageSize)
	require.NotNil(t, input.Cursor.Timestamp)
	require.NotNil(t, input.Cursor.ID)
	assert.Equal(t, cursorTime, *input.Cursor.Timestamp)
	assert.Equal(t, cursorID, *input.Cursor.ID)
}

func TestParseListDeliveryAttemptsInput_AppliesDefaults(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/events/9/delivery-attempts", nil)
	req = withURLParam(req, "event_id", "9")

	input, err := api.ParseRequestInput[ListDeliveryAttemptsInput](req)

	require.NoError(t, err)
	assert.Equal(t, int64(9), input.EventID)
	assert.Equal(t, "", input.Search)
	assert.Equal(t, "*", input.Filter)
	assert.Equal(t, "DESC", input.SortDirection)
	assert.Equal(t, 20, input.PageSize)
	assert.Nil(t, input.Cursor.Timestamp)
	assert.Nil(t, input.Cursor.ID)
}

func TestParseListDeliveryAttemptsInput_InvalidFilter(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		"GET",
		"/events/9/delivery-attempts?filter_state=unknown",
		nil,
	)
	req = withURLParam(req, "event_id", "9")

	input, err := api.ParseRequestInput[ListDeliveryAttemptsInput](req)

	assert.Nil(t, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "filter_state")
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	return withRouteParam(req, key, value)
}
