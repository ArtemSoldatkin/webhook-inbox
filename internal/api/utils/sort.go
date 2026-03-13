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
) SortDirection {
	order := SortDirection(query.Get("sort_direction"))
	if order == SortDirectionAsc || order == SortDirectionDesc {
		return order
	}
	return defaultSortDirection
}
