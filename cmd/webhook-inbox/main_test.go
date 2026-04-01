package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/sirupsen/logrus"
)

func TestLogrusLogger_CallsNextHandlerAndPreservesResponse(t *testing.T) {
	var logBuffer bytes.Buffer
	restoreLogger := setTestLoggerOutput(&logBuffer)
	defer restoreLogger()

	called := false
	handler := LogrusLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("created"))
		require.NoError(t, err)
	}))

	req := httptest.NewRequest(http.MethodPost, "/webhooks/ingest", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "created", recorder.Body.String())
	assert.Contains(t, logBuffer.String(), "Handled request")
	assert.Contains(t, logBuffer.String(), "method=POST")
	assert.Contains(t, logBuffer.String(), "path=/webhooks/ingest")
	assert.Contains(t, logBuffer.String(), "duration=")
}

func TestLogrusLogger_LogsForRequestsWithoutBodyWrites(t *testing.T) {
	var logBuffer bytes.Buffer
	restoreLogger := setTestLoggerOutput(&logBuffer)
	defer restoreLogger()

	handler := LogrusLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Empty(t, recorder.Body.String())
	assert.Contains(t, logBuffer.String(), "method=GET")
	assert.Contains(t, logBuffer.String(), "path=/health")
}

func setTestLoggerOutput(buffer *bytes.Buffer) func() {
	originalOutput := logrus.StandardLogger().Out
	originalFormatter := logrus.StandardLogger().Formatter
	originalLevel := logrus.StandardLogger().Level

	logrus.SetOutput(buffer)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	return func() {
		logrus.SetOutput(originalOutput)
		logrus.SetFormatter(originalFormatter)
		logrus.SetLevel(originalLevel)
	}
}
