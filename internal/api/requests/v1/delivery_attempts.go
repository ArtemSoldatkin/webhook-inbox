package requestsv1

import "github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"

// ListDeliveryAttemptsInput defines the expected input parameters for listing delivery attempts.
type ListDeliveryAttemptsInput struct {
	EventID       int64        `url_param:"event_id,required,min:1"`
	Search        string       `query_param:"search,max_length:512"`
	Filter        string       `query_param:"filter_state,allowed:pending|in_flight|succeeded|failed|aborted,default:*"`
	SortDirection string       `query_param:"sort,allowed:ASC|DESC,default:DESC"`
	PageSize      int          `query_param:"limit,min:1,max:100,default:20"`
	Cursor        types.Cursor `query_param:"cursor"`
}
