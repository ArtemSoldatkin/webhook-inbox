package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

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
	if r.Body == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// Restore the body for downstream handlers, even if it's empty.
	defer func() {
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}()

	if len(bodyBytes) == 0 {
		// Treat empty body as "no body", not an error.
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	if err := decoder.Decode(params); err != nil {
		if err == io.EOF {
			// Treat EOF as "no body".
			return nil
		}
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

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := ParseJsonBodyParams(r, &params); err != nil {
			return nil, err
		}
	}

	return &params, nil
}
