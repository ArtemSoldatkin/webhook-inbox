package routev1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestHealth_ReturnsOKStatus(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	health().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[map[string]string](t, recorder)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
}

func TestHealthRouter_MountsHealthEndpoint(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	systemRouter().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	response := decodeJSONResponse[map[string]string](t, recorder)
	assert.Equal(t, "ok", response["status"])
}
