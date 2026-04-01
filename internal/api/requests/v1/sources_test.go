package requestsv1

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestValidateCreateSourceInput(t *testing.T) {
	t.Parallel()

	err := ValidateCreateSourceInput(&CreateSourceInput{
		EgressUrl: "https://example.com/webhook",
		StaticHeaders: map[string]string{
			"Authorization": "Bearer token",
			"X-Trace":       "abc123",
		},
		Description: "customer billing events",
	})

	require.NoError(t, err)
}

func TestValidateCreateSourceInput_DescriptionTooLong(t *testing.T) {
	t.Parallel()

	err := ValidateCreateSourceInput(&CreateSourceInput{
		Description: strings.Repeat("a", maxDescriptionLen+1),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "description exceeds maximum length")
}

func TestValidateCreateSourceInput_TooManyHeaders(t *testing.T) {
	t.Parallel()

	headers := make(map[string]string, maxHeaders+1)
	for i := range maxHeaders + 1 {
		headers[string(rune('a'+i))] = "value"
	}

	err := ValidateCreateSourceInput(&CreateSourceInput{
		StaticHeaders: headers,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "too many headers")
}

func TestValidateCreateSourceInput_HeaderLengthExceeded(t *testing.T) {
	t.Parallel()

	err := ValidateCreateSourceInput(&CreateSourceInput{
		StaticHeaders: map[string]string{
			strings.Repeat("k", maxHeaderKeyLen+1): "value",
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "header key or value exceeds maximum length")
}

func TestValidateUpdateSourceStatusInput(t *testing.T) {
	t.Parallel()

	err := ValidateUpdateSourceStatusInput(&UpdateSourceStatusInput{
		SourceID:     14,
		Status:       "paused",
		StatusReason: "manual pause",
	})

	require.NoError(t, err)
}

func TestValidateUpdateSourceStatusInput_InvalidSourceID(t *testing.T) {
	t.Parallel()

	err := ValidateUpdateSourceStatusInput(&UpdateSourceStatusInput{
		SourceID: 0,
		Status:   "active",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source ID")
}

func TestValidateUpdateSourceStatusInput_InvalidStatus(t *testing.T) {
	t.Parallel()

	err := ValidateUpdateSourceStatusInput(&UpdateSourceStatusInput{
		SourceID: 1,
		Status:   "deleted",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestValidateUpdateSourceStatusInput_StatusReasonTooLong(t *testing.T) {
	t.Parallel()

	err := ValidateUpdateSourceStatusInput(&UpdateSourceStatusInput{
		SourceID:     1,
		Status:       "disabled",
		StatusReason: strings.Repeat("a", maxDescriptionLen+1),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status reason exceeds maximum length")
}
