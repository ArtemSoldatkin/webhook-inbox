package api

import (
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestToPaginatedResponse_WithNextPage(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	id := int64(42)
	cursor := types.NewCursor(&timestamp, &id)

	response := ToPaginatedResponse([]int{1, 2, 3}, 2, cursor)

	assert.Equal(t, []int{1, 2}, response.Data)
	assert.True(t, response.HasNext)
	require.NotNil(t, response.NextCursor)
	assert.Equal(t, cursor.ToString(), *response.NextCursor)
}

func TestToPaginatedResponse_WithoutNextPage(t *testing.T) {
	t.Parallel()

	response := ToPaginatedResponse([]int{1, 2}, 2, types.Cursor{})

	assert.Equal(t, []int{1, 2}, response.Data)
	assert.False(t, response.HasNext)
	assert.Nil(t, response.NextCursor)
}
