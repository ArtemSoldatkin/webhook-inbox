package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()

	err := JSON(recorder, http.StatusCreated, map[string]string{"status": "ok"})

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"status":"ok"}`, recorder.Body.String())
}

func TestJSON_ReturnsMarshalError(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	err := JSON(recorder, http.StatusOK, map[string]any{
		"invalid": make(chan int),
	})

	var marshalErr *JSONMarshalError
	require.Error(t, err)
	require.True(t, errors.As(err, &marshalErr))
	assert.Equal(t, "failed to marshal response data", marshalErr.Error())
	require.Error(t, marshalErr.Unwrap())
}

func TestJSON_ReturnsWriteError(t *testing.T) {
	t.Parallel()

	writer := &failingResponseWriter{header: make(http.Header)}

	err := JSON(writer, http.StatusOK, map[string]string{"status": "ok"})

	var writeErr *JSONWriteError
	require.Error(t, err)
	require.True(t, errors.As(err, &writeErr))
	assert.Equal(t, "failed to write response", writeErr.Error())
	require.Error(t, writeErr.Unwrap())
	assert.Equal(t, "application/json", writer.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, writer.statusCode)
}

type failingResponseWriter struct {
	header     http.Header
	statusCode int
}

func (w *failingResponseWriter) Header() http.Header {
	return w.header
}

func (w *failingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *failingResponseWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failure")
}
