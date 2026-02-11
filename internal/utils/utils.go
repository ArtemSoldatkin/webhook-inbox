// Package utils provides utility functions for the application.
package utils

import "fmt"

// GenerateIngressURL generates the ingress URL for a given source ID.
func GenerateIngressURL(sourceID string) string {
	return fmt.Sprintf("http://localhost:3001/ingest/%s", sourceID) // TODO: make this dynamic based on config
}
