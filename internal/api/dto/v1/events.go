package dtov1

import "time"

// EventDTO represents the data transfer object for an event in the API.
type EventDTO struct {
	ID              int64
	SourceID        int64
	DedupHash       string
	Method          string
	IngressPath     string
	RemoteAddress   string
	QueryParams     map[string][]string
	RawHeaders      map[string][]string
	Body           	map[string]string
	BodyContentType string
	ReceivedAt      time.Time
}