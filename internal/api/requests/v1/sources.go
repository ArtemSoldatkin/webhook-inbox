// Package requestsv1 defines the request structures and validation logic for API version 1.
package requestsv1

import (
	"fmt"

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
	Filter        string       `query_param:"filter_state,allowed:active|paused|quarantined|disabled|*,default:*"`
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
	EgressUrl     string            `json:"EgressUrl"`
	StaticHeaders map[string]string `json:"StaticHeaders,omitempty"`
	Description   string            `json:"Description,omitempty"`
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
