// Package requestsv1 defines the request structures and validation logic for API version 1.
package requestsv1

import (
	"fmt"
	"slices"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
)

const (
	maxDescriptionLen = 500
	maxHeaders        = 20
	maxHeaderKeyLen   = 100
	maxHeaderValueLen = 500
)

// ListSourcesInput defines the expected input parameters for listing sources.
type ListSourcesInput struct {
	Filter        string       `query_param:"filter_status,allowed:active|paused|quarantined|disabled|*,default:*"`
	SortDirection string       `query_param:"sort_direction,allowed:ASC|DESC,default:DESC"`
	Search        string       `query_param:"search,max_length:512"`
	PageSize      int          `query_param:"limit,min:1,max:100,default:20"`
	Cursor        types.Cursor `query_param:"cursor"`
}

// GetSourceByIDInput defines the expected input parameters for retrieving a source by its ID.
type GetSourceByIDInput struct {
	SourceID int64 `url_param:"source_id,required,min:1"`
}

// CreateSourceData defines the expected input parameters for creating a new source.
type CreateSourceInput struct {
	EgressUrl     string            `json:"egress_url"`
	StaticHeaders map[string]string `json:"static_headers,omitempty"`
	Description   string            `json:"description,omitempty"`
}

// ValidateCreateSourceInput validates the input parameters for creating a new source.
func ValidateCreateSourceInput(input *CreateSourceInput) error {
	if len(input.Description) > maxDescriptionLen {
		return fmt.Errorf("description exceeds maximum length of %d characters", maxDescriptionLen)
	}

	if len(input.StaticHeaders) > maxHeaders {
		return fmt.Errorf("too many headers: maximum allowed is %d", maxHeaders)
	}
	for k, v := range input.StaticHeaders {
		if len(k) > maxHeaderKeyLen || len(v) > maxHeaderValueLen {
			return fmt.Errorf("header key or value exceeds maximum length (key: %d, value: %d)", maxHeaderKeyLen, maxHeaderValueLen)
		}
	}

	return nil
}

// allowedStatuses defines the valid statuses for sources.
var allowedStatuses = []string{"active", "paused", "quarantined", "disabled"}

// UpdateSourceStatusInput defines the expected input parameters for updating the status of a source.
type UpdateSourceStatusInput struct {
	SourceID     int64  `url_param:"source_id"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason,omitempty"`
}

// ValidateUpdateSourceStatusInput validates the input parameters for updating the status of a source.
func ValidateUpdateSourceStatusInput(input *UpdateSourceStatusInput) error {
	if input.SourceID <= 0 {
		return fmt.Errorf("invalid source ID: must be a positive integer")
	}

	if !slices.Contains(allowedStatuses, input.Status) {
		return fmt.Errorf("invalid status: must be one of %v", allowedStatuses)
	}

	if len(input.StatusReason) > maxDescriptionLen {
		return fmt.Errorf("status reason exceeds maximum length of %d characters", maxDescriptionLen)
	}

	return nil
}
