// Package api provides utilities for handling API requests, including pagination.
package api

import (
	"fmt"
	"net/url"
	"strconv"

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

// ParsePaginationParams extracts pagination parameters from the query string.
func ParsePaginationParams(
	query url.Values,
	defaultPageSize,
	minPageSize,
	maxPageSize int,
) (
	pageSize int,
	cursor types.Cursor,
	err error,
) {
	pageSize, err = getIntQueryParam(query, "limit", defaultPageSize)
	if err != nil {
		return 0, cursor, err
	}

	if pageSize < minPageSize || pageSize > maxPageSize {
		return 0, cursor, fmt.Errorf(
			"limit parameter must be between %d and %d",
			minPageSize,
			maxPageSize,
		)
	}

	cursorStr := query.Get("cursor")
	if err := cursor.FromString(cursorStr); err != nil {
		return 0, cursor, fmt.Errorf("invalid cursor parameter: %w", err)
	}

	return pageSize, cursor, nil
}

// getIntQueryParam retrieves an integer query parameter with a default value.
func getIntQueryParam(query url.Values, key string, defaultValue int) (value int, err error) {
	valueStr := query.Get(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err = strconv.Atoi(valueStr)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf(
			"invalid %s parameter",
			key,
		)
	}

	return value, nil
}
