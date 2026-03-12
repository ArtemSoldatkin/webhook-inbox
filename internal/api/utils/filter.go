package api

import (
	"fmt"
	"net/url"
)

// ParseFilter extracts a filter value from the query parameters based on the provided filter name and options.
// If the filter value is not present or not valid, it returns a default value of "*".
func ParseFilter(
	query url.Values,
	filterName string,
	filterOptions map[string]bool,
) string {
	filter := query.Get(fmt.Sprintf("filter_%s", filterName))
	if filterOptions[filter] {
		return filter
	}
	return "*"
}
