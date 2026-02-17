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

// MergeMaps merges two maps of string keys and values, with values from the second map taking precedence in case of key conflicts.
func MergeHeaders(staticHeaders map[string]string, rawHeaders map[string][]string) map[string][]string {
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

// GenerateIngressURL generates the ingress URL for a given source ID.
func GenerateIngressURL(sourceID string) string {
	return fmt.Sprintf("http://localhost:3001/ingest/%s", sourceID) // TODO: make this dynamic based on config
}
