package routev1

import (
	"errors"
	"net/http"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// ListDeliveryAttemptsInput defines the expected input parameters for listing delivery attempts.
type ListDeliveryAttemptsInput struct {
	EventID       int64        `url_param:"event_id,required,min=1"`
	Search        string       `query_param:"search,max_length=512"`
	Filter        string       `query_param:"filter_state,allowed=pending|in_flight|succeeded|failed|aborted,default=*"`
	SortDirection string       `query_param:"sort,allowed=ASC|DESC,default=DESC"`
	PageSize      int          `query_param:"limit,min=1,max=100,default=20"`
	Cursor        types.Cursor `query_param:"cursor"`
}

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[ListDeliveryAttemptsInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"event_id":       input.EventID,
			"search":         input.Search,
			"filter_state":   input.Filter,
			"sort_direction": input.SortDirection,
			"page_size":      input.PageSize,
			"cursor":         input.Cursor,
			"query":          r.URL.RawQuery,
		}).Debug("Received listDeliveryAttempts request")

		deliveryAttempts, err := svc.ListDeliveryAttempts(
			r.Context(),
			input.EventID,
			input.Cursor,
			input.PageSize,
			input.Search,
			input.Filter,
			input.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}

		deliveryAttemptDTOs := mapperv1.ToDeliveryAttemptDTOs(deliveryAttempts)

		logrus.WithFields(logrus.Fields{
			"event_id":       input.EventID,
			"returned_count": len(deliveryAttemptDTOs),
		}).Debug("Returning delivery attempts")

		var nextCursor types.Cursor
		if len(deliveryAttemptDTOs) > input.PageSize {
			lastAttempt := deliveryAttemptDTOs[len(deliveryAttemptDTOs)-1]
			nextCursor = types.NewCursor(
				&lastAttempt.CreatedAt,
				&lastAttempt.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			deliveryAttemptDTOs,
			input.PageSize,
			nextCursor,
		)

		if err := api.JSON(w, http.StatusOK, paginatedResponse); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			}
		}
	}
}

// deliveryAttemptsRouter sets up the router for delivery attempts-related endpoints.
func deliveryAttemptsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listDeliveryAttempts(svc))
	return r
}
