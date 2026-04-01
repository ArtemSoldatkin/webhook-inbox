package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type parserInput struct {
	SourceID      int64        `url_param:"source_id,required,min:1"`
	Search        string       `query_param:"search,max_length:20"`
	SortDirection string       `query_param:"sort_direction,allowed:ASC|DESC,default:DESC"`
	PageSize      int          `query_param:"limit,min:1,max:100,default:20"`
	Cursor        types.Cursor `query_param:"cursor"`
	Name          string       `json:"name"`
}

func TestParseURLParams(t *testing.T) {
	t.Parallel()

	req := withTestURLParam(httptest.NewRequest("GET", "/", nil), "source_id", "15")
	var input parserInput

	err := ParseUrlParams(req, &input)

	require.NoError(t, err)
	assert.Equal(t, int64(15), input.SourceID)
}

func TestParseQueryParams(t *testing.T) {
	t.Parallel()

	var input parserInput
	err := ParseQueryParams(
		reqQuery("search=invoice&sort_direction=ASC&limit=10"),
		&input,
	)

	require.NoError(t, err)
	assert.Equal(t, "invoice", input.Search)
	assert.Equal(t, "ASC", input.SortDirection)
	assert.Equal(t, 10, input.PageSize)
}

func TestParseJsonBodyParams(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"webhook"}`))
	var input parserInput

	err := ParseJsonBodyParams(req, &input)

	require.NoError(t, err)
	assert.Equal(t, "webhook", input.Name)

	body, readErr := io.ReadAll(req.Body)
	require.NoError(t, readErr)
	assert.Equal(t, `{"name":"webhook"}`, string(body))
}

func TestParseJsonBodyParams_NilBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/", nil)
	req.Body = nil
	var input parserInput

	require.NoError(t, ParseJsonBodyParams(req, &input))
	assert.Equal(t, "", input.Name)
}

func TestParseJsonBodyParams_EmptyBody(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/", http.NoBody)
	var input parserInput

	require.NoError(t, ParseJsonBodyParams(req, &input))
	assert.Equal(t, "", input.Name)
}

func TestParseJsonBodyParams_InvalidJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":`))
	var input parserInput

	err := ParseJsonBodyParams(req, &input)

	require.Error(t, err)
}

func TestParseRequestInput(t *testing.T) {
	t.Parallel()

	req := withTestURLParam(
		httptest.NewRequest(
			"POST",
			"/sources/12?search=invoice&limit=5",
			strings.NewReader(`{"name":"created"}`),
		),
		"source_id",
		"12",
	)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	input, err := ParseRequestInput[parserInput](req)

	require.NoError(t, err)
	require.NotNil(t, input)
	assert.Equal(t, int64(12), input.SourceID)
	assert.Equal(t, "invoice", input.Search)
	assert.Equal(t, "DESC", input.SortDirection)
	assert.Equal(t, 5, input.PageSize)
	assert.Equal(t, "created", input.Name)
	assert.Nil(t, input.Cursor.Timestamp)
	assert.Nil(t, input.Cursor.ID)
}

func TestParseRequestInput_InvalidURLParam(t *testing.T) {
	t.Parallel()

	req := withTestURLParam(httptest.NewRequest("GET", "/", nil), "source_id", "bad")

	input, err := ParseRequestInput[parserInput](req)

	assert.Nil(t, input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source_id")
}

func TestParseRequestInput_SkipsJSONWhenContentTypeIsNotJSON(t *testing.T) {
	t.Parallel()

	req := withTestURLParam(
		httptest.NewRequest("POST", "/sources/12", strings.NewReader(`{"name":"ignored"}`)),
		"source_id",
		"12",
	)
	req.Header.Set("Content-Type", "text/plain")

	input, err := ParseRequestInput[parserInput](req)

	require.NoError(t, err)
	assert.Equal(t, "", input.Name)
}

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) { return 0, errors.New("read failure") }

func TestParseJsonBodyParams_ReadError(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/", nil)
	req.Body = io.NopCloser(errReader{})
	var input parserInput

	err := ParseJsonBodyParams(req, &input)

	require.Error(t, err)
	assert.Equal(t, "read failure", err.Error())
}

func withTestURLParam(req *http.Request, key, value string) *http.Request {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
}

func reqQuery(raw string) map[string][]string {
	req := httptest.NewRequest("GET", "/?"+raw, nil)
	return req.URL.Query()
}
