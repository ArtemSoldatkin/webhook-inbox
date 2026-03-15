package routev1

import (
	"errors"
	"net/http"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// ListDeliveryAttemptsURLParams defines the URL parameters for the list delivery attempts endpoint.
type ListDeliveryAttemptsURLParams struct {
	EventID int64 `param:"event_id,required,min=1"`
}

// ListDeliveryAttemptsQueryParams defines the query parameters for the list delivery attempts endpoint.
type ListDeliveryAttemptsQueryParams struct {
	Search        string `param:"search,max_length=512"`
	Filter        string `param:"filter_state,allowed=pending|in_flight|succeeded|failed|aborted,default=*"`
	SortDirection string `param:"sort,allowed=ASC|DESC,default=DESC"`
	PageSize      int    `param:"limit,min=1,max=100,default=20"`
	Cursor        string `param:"cursor"`
}

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParams, err := api.ParseUrlParams[ListDeliveryAttemptsURLParams](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse URL parameters")
			http.Error(w, "Invalid URL parameters", http.StatusBadRequest)
			return
		}

		queryParams, err := api.ParseQueryParams[ListDeliveryAttemptsQueryParams](r.URL.Query())
		if err != nil {
			logrus.WithError(err).Error("Failed to parse query parameters")
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"event_id":       urlParams.EventID,
			"search":         queryParams.Search,
			"filter_state":   queryParams.Filter,
			"sort_direction": queryParams.SortDirection,
			"page_size":      queryParams.PageSize,
			"query":          r.URL.RawQuery,
		}).Debug("Received listDeliveryAttempts request")

		// TODO: Move logic into ParseQueryParams and validate parameters there
		pageSize, cursor, err := api.ParsePaginationParams(
			r.URL.Query(),
			svc.Config.APIDefaultPageSize,
			svc.Config.APIMinPageSize,
			svc.Config.APIMaxPageSize,
		)
		if err != nil {
			logrus.
				WithError(err).
				Error("Invalid pagination parameters")
			http.Error(w, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		deliveryAttempts, err := svc.ListDeliveryAttempts(
			r.Context(),
			urlParams.EventID,
			cursor,
			pageSize,
			queryParams.Search,
			queryParams.Filter,
			queryParams.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}

		deliveryAttemptDTOs := mapperv1.ToDeliveryAttemptDTOs(deliveryAttempts)

		logrus.WithFields(logrus.Fields{
			"event_id":       urlParams.EventID,
			"returned_count": len(deliveryAttemptDTOs),
		}).Debug("Returning delivery attempts")

		var nextCursor api.Cursor
		if len(deliveryAttemptDTOs) > pageSize {
			lastAttempt := deliveryAttemptDTOs[len(deliveryAttemptDTOs)-1]
			nextCursor = api.NewCursor(
				&lastAttempt.CreatedAt,
				&lastAttempt.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			deliveryAttemptDTOs,
			pageSize,
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
