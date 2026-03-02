// Package api provides utilities for handling API requests, including pagination.
package api

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// PaginatedResponse represents a standard structure for paginated API responses.
type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	NextCursor *time.Time `json:"nextCursor,omitempty"`
	HasNext    bool       `json:"hasNext"`
}

// ToPaginatedResponse converts a slice of data into a PaginatedResponse, determining if there is a next page based on the provided page size and cursor function.
func ToPaginatedResponse[T any](data []T, pageSize int, getCursor func(T) *time.Time) PaginatedResponse[T] {
	hasNext := len(data) > pageSize
	var nextCursor *time.Time
	if hasNext {
		nextCursor = getCursor(data[pageSize])
		data = data[:pageSize]
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
	cursor *time.Time,
	err error,
) {
	pageSize, err = getIntQueryParam(query, "limit", defaultPageSize)
	if err != nil {
		return 0, nil, err
	}
	if pageSize < minPageSize || pageSize > maxPageSize {
		return 0, nil, fmt.Errorf(
			"limit parameter must be between %d and %d",
			minPageSize,
			maxPageSize,
		)
	}

	cursor, err = getTimeQueryParam(query, "cursor")
	if err != nil {
		return 0, nil, err
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

// getTimeQueryParam retrieves a time query parameter in RFC3339 format.
func getTimeQueryParam(query url.Values, key string) (value *time.Time, err error) {
	valueStr := query.Get(key)
	if valueStr == "" {
		return nil, nil
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, valueStr)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, valueStr)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"invalid %s parameter",
			key,
		)
	}
	return &parsedTime, nil
}
