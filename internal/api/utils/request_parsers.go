package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	structparser "github.com/ArtemSoldatkin/webhook-inbox/internal/struct_parser"
	"github.com/go-chi/chi/v5"
)

// ParseUrlParams is a generic function that parses URL parameters
// from an HTTP request into a struct based on struct tags.
func ParseUrlParams[T any](r *http.Request, params *T) error {
	if err := structparser.ParseStruct(
		params,
		"url_param",
		func(varName string) string {
			return chi.URLParam(r, varName)
		},
		true,
	); err != nil {
		return err
	}

	return nil
}

// ParseQueryParams is a generic function that parses query parameters
// from an HTTP request into a struct based on struct tags.
func ParseQueryParams[T any](queryParams url.Values, params *T) error {
	if err := structparser.ParseStruct(
		params,
		"query_param",
		func(varName string) string {
			return queryParams.Get(varName)
		},
		true,
	); err != nil {
		return err
	}

	return nil
}

// ParseJsonBodyParams is a generic function that parses JSON body parameters
// from an HTTP request into a struct based on struct tags.
func ParseJsonBodyParams[T any](r *http.Request, params *T) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		return err
	}

	return nil
}

// ParseRequestInput is a generic function that parses URL parameters, query parameters,
// and JSON body parameters from an HTTP request into a struct based on struct tags.
func ParseRequestInput[T any](r *http.Request) (*T, error) {
	var params T

	if err := ParseUrlParams(r, &params); err != nil {
		return nil, err
	}

	if err := ParseQueryParams(r.URL.Query(), &params); err != nil {
		return nil, err
	}

	if err := ParseJsonBodyParams(r, &params); err != nil {
		return nil, err
	}

	return &params, nil
}
