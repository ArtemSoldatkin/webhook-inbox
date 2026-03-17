package requestsv1

import "github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"

// ListEventsInput defines the expected input parameters for listing events.
type ListEventsInput struct {
	SourceID      int64        `url_param:"source_id,required,min:1"`
	Search        string       `query_param:"search,max_length:512"`
	SortDirection string       `query_param:"sort,allowed:ASC|DESC,default:DESC"`
	PageSize      int          `query_param:"limit,min:1,max:100,default:20"`
	Cursor        types.Cursor `query_param:"cursor"`
}

// GetEventInput defines the expected input parameters for retrieving a specific event.
type GetEventInput struct {
	EventID int64 `url_param:"event_id,required,min:1"`
}
