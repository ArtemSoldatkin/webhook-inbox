// Package api provides utilities for handling API requests, including pagination.
package api

import (
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
)

// PaginatedResponse represents a standard structure for paginated API responses.
type PaginatedResponse[T any] struct {
	Data       []T     `json:"data"`
	NextCursor *string `json:"next_cursor,omitempty"`
	HasNext    bool    `json:"has_next"`
}

// ToPaginatedResponse converts a slice of data into a PaginatedResponse,
// determining if there is a next page based on the provided page size and cursor function.
func ToPaginatedResponse[T any](
	data []T,
	pageSize int,
	cursor types.Cursor,
) PaginatedResponse[T] {
	hasNext := len(data) > pageSize
	var nextCursor *string
	if hasNext {
		data = data[:pageSize]
		cursorStr := cursor.ToString()
		nextCursor = &cursorStr
	}

	return PaginatedResponse[T]{
		Data:       data,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}
}
