package dtov1

import "time"

// EventDTO represents the data transfer object for an event in the API.
type EventDTO struct {
	ID       int64 `json:"id"`
	SourceID int64 `json:"source_id"`
	// DedupHash       *string
	Method          string              `json:"method"`
	IngressPath     string              `json:"ingress_path"`
	RemoteAddress   *string             `json:"remote_address,omitempty"`
	QueryParams     map[string][]string `json:"query_params,omitempty"`
	RawHeaders      map[string][]string `json:"raw_headers,omitempty"`
	Body            []byte              `json:"body,omitempty"`
	BodyContentType string              `json:"body_content_type,omitempty"`
	ReceivedAt      time.Time           `json:"received_at"`
}
