// Package utils provides utility functions for the application.
package utils

import (
	"encoding/json"
	"fmt"
)

// JSONBtoType converts a JSONB byte array to a Go map of string keys and values.
func JSONBtoType[T any](jsonb []byte) (T, error) {
	var result T
	err := json.Unmarshal(jsonb, &result)
	return result, err
}

// MergeHeaders merges static headers and raw headers into a single map of headers.
func MergeHeaders(
	staticHeaders map[string]string,
	rawHeaders map[string][]string,
) map[string][]string {
	headers := make(map[string][]string)
	for key, value := range staticHeaders {
		headers[key] = []string{value}
	}

	for key, values := range rawHeaders {
		if existingValues, exists := headers[key]; exists {
			headers[key] = append(existingValues, values...)
		} else {
			headers[key] = values
		}
	}

	return headers
}

// GenerateIngressURL generates an ingress URL for a source based on the API protocol,
// host, port, and source ID.
func GenerateIngressURL(apiProtocol, apiHost string, apiPort int, sourceID string) string {
	return fmt.Sprintf(
		"%s://%s:%d/ingest/%s",
		apiProtocol,
		apiHost,
		apiPort,
		sourceID,
	)
}

// PtrIfValid returns a pointer to the value if the valid flag is true, otherwise it returns nil.
func PtrIfValid[T any](value T, valid bool) *T {
	var result *T
	if valid {
		result = &value
	}

	return result
}
