package api

import "net/url"

// SortDirection represents the direction of sorting (ascending or descending).
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

// ParseSortDirection parses the order direction from the query parameters.
func ParseSortDirection(
	query url.Values, defaultSortDirection SortDirection,
) string {
	order := query.Get("sort_direction")
	if order == "ASC" || order == "DESC" {
		return order
	}
	return string(defaultSortDirection)
}
