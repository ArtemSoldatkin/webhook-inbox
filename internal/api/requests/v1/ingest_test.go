package requestsv1

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestValidatePublicID(t *testing.T) {
	t.Parallel()

	assert.True(t, validatePublicID("123e4567-e89b-12d3-a456-426614174000"))
	assert.True(t, validatePublicID("123E4567-E89B-12D3-A456-426614174000"))
	assert.False(t, validatePublicID("123e4567e89b12d3a456426614174000"))
	assert.False(t, validatePublicID("not-a-uuid"))
}

func TestValidateIngestEventInput(t *testing.T) {
	t.Parallel()

	err := ValidateIngestEventInput(&IngestEventInput{
		PublicID: "123e4567-e89b-12d3-a456-426614174000",
	})

	require.NoError(t, err)
}

func TestValidateIngestEventInput_InvalidPublicID(t *testing.T) {
	t.Parallel()

	err := ValidateIngestEventInput(&IngestEventInput{
		PublicID: "bad-id",
	})

	require.Error(t, err)
	assert.Equal(t, "invalid public_id format", err.Error())
}
