package api

import (
	"net/http"
	"net/url"

	structparser "github.com/ArtemSoldatkin/webhook-inbox/internal/struct_parser"
	"github.com/go-chi/chi/v5"
)

// ParseUrlParams is a generic function that parses URL parameters
// from an HTTP request into a struct based on struct tags.
func ParseUrlParams[T any](r *http.Request) (*T, error) {
	var params T

	if err := structparser.ParseStruct(&params, "param", func(varName string) string {
		return chi.URLParam(r, varName)
	}); err != nil {
		return nil, err
	}

	return &params, nil
}

// ParseQueryParams is a generic function that parses query parameters
// from an HTTP request into a struct based on struct tags.
func ParseQueryParams[T any](queryParams url.Values) (*T, error) {
	var params T

	if err := structparser.ParseStruct(&params, "param", func(varName string) string {
		return queryParams.Get(varName)
	}); err != nil {
		return nil, err
	}

	return &params, nil
}
