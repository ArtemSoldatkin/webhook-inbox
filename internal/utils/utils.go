// Package utils provides utility functions for the application.
package utils

import (
	"encoding/json"
	"fmt"
)

// JSONBtoMap converts a JSONB byte array to a map[string]string.
func JSONBtoMap(jsonb []byte) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal(jsonb, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateIngressURL generates the ingress URL for a given source ID.
func GenerateIngressURL(sourceID string) string {
	return fmt.Sprintf("http://localhost:3001/ingest/%s", sourceID) // TODO: make this dynamic based on config
}
