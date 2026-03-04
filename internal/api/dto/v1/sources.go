// Package dtov1 contains data transfer objects for version 1 of the API.
package dtov1

import (
	"time"
)

// SourceDTO represents the data transfer object for a source in the API.
type SourceDTO struct {
	ID            int64             `json:"id"`
	PublicID      string            `json:"public_id"`
	IngressUrl    string            `json:"ingress_url"`
	EgressUrl     string            `json:"egress_url"`
	StaticHeaders map[string]string `json:"static_headers,omitempty"`
	Status        string            `json:"status"`
	StatusReason  *string           `json:"status_reason,omitempty"`
	Description   *string           `json:"description,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DisableAt     *time.Time        `json:"disable_at,omitempty"`
}
