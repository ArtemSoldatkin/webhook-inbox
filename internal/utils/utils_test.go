package utils

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestJSONBtoType(t *testing.T) {
	t.Parallel()

	value, err := JSONBtoType[map[string][]string]([]byte(`{"foo":["bar","baz"]}`))

	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"foo": {"bar", "baz"}}, value)
}

func TestJSONBtoType_InvalidJSON(t *testing.T) {
	t.Parallel()

	value, err := JSONBtoType[map[string]string]([]byte(`{`))

	assert.Nil(t, value)
	require.Error(t, err)
}

func TestMergeHeaders(t *testing.T) {
	t.Parallel()

	headers := MergeHeaders(
		map[string]string{
			"Authorization": "Bearer token",
			"X-Static":      "yes",
		},
		map[string][]string{
			"Authorization": {"dynamic-token"},
			"X-Raw":         {"raw-value"},
		},
	)

	assert.Equal(t, []string{"Bearer token", "dynamic-token"}, headers["Authorization"])
	assert.Equal(t, []string{"yes"}, headers["X-Static"])
	assert.Equal(t, []string{"raw-value"}, headers["X-Raw"])
}

func TestGenerateIngressURL(t *testing.T) {
	t.Parallel()

	url := GenerateIngressURL("https", "api.example.com", 8443, "source-123")

	assert.Equal(t, "https://api.example.com:8443/ingest/source-123", url)
}

func TestPtrIfValid(t *testing.T) {
	t.Parallel()

	value := PtrIfValid("test", true)
	require.NotNil(t, value)
	assert.Equal(t, "test", *value)

	assert.Nil(t, PtrIfValid("test", false))
}
