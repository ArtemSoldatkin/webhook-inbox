package requestsv1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestParseListEventsInput(t *testing.T) {
	t.Parallel()

	cursorTime := time.Date(2026, 4, 1, 11, 22, 33, 0, time.UTC)
	cursorID := int64(88)
	cursor := newTestCursor(cursorTime, cursorID)

	req := httptest.NewRequest(
		"GET",
		"/sources/21/events?search=invoice&sort_direction=ASC&limit=50&cursor="+cursor,
		nil,
	)
	req = withRouteParam(req, "source_id", "21")

	input, err := api.ParseRequestInput[ListEventsInput](req)

	require.NoError(t, err)
	require.NotNil(t, input)
	assert.Equal(t, int64(21), input.SourceID)
	assert.Equal(t, "invoice", input.Search)
	assert.Equal(t, "ASC", input.SortDirection)
	assert.Equal(t, 50, input.PageSize)
	require.NotNil(t, input.Cursor.Timestamp)
	require.NotNil(t, input.Cursor.ID)
	assert.Equal(t, cursorTime, *input.Cursor.Timestamp)
	assert.Equal(t, cursorID, *input.Cursor.ID)
}

func TestParseListEventsInput_AppliesDefaults(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/sources/7/events", nil)
	req = withRouteParam(req, "source_id", "7")

	input, err := api.ParseRequestInput[ListEventsInput](req)

	require.NoError(t, err)
	assert.Equal(t, int64(7), input.SourceID)
	assert.Equal(t, "", input.Search)
	assert.Equal(t, "DESC", input.SortDirection)
	assert.Equal(t, 20, input.PageSize)
	assert.Nil(t, input.Cursor.Timestamp)
	assert.Nil(t, input.Cursor.ID)
}

func TestParseListEventsInput_InvalidPageSize(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/sources/7/events?limit=101", nil)
	req = withRouteParam(req, "source_id", "7")

	input, err := api.ParseRequestInput[ListEventsInput](req)

	assert.Nil(t, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "limit")
}

func TestParseGetEventInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/events/42", nil)
	req = withRouteParam(req, "event_id", "42")

	input, err := api.ParseRequestInput[GetEventInput](req)

	require.NoError(t, err)
	require.NotNil(t, input)
	assert.Equal(t, int64(42), input.EventID)
}

func withRouteParam(req *http.Request, key, value string) *http.Request {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(key, value)
	return req.WithContext(contextWithRouteContext(req, routeContext))
}

func contextWithRouteContext(req *http.Request, routeContext *chi.Context) context.Context {
	return context.WithValue(req.Context(), chi.RouteCtxKey, routeContext)
}

func newTestCursor(ts time.Time, id int64) string {
	timestamp := ts
	cursorID := id
	cursor := types.NewCursor(&timestamp, &cursorID)
	return cursor.ToString()
}
