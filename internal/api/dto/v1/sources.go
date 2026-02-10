// Package dtov1 contains data transfer objects for version 1 of the API.
package dtov1

import (
	"time"
)

// SourceDTO represents the data transfer object for a source in the API.
type SourceDTO struct {
	ID            int64
	IngressUrl    string
	EgressUrl     string
	StaticHeaders map[string]string
	Status        string
	StatusReason  string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DisableAt     *time.Time
}
