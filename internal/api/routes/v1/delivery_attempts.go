package routev1

import (
	"net/http"
	"strconv"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

var filterStateOptions = map[string]bool{
	"pending":   true,
	"in_flight": true,
	"succeeded": true,
	"failed":    true,
	"aborted":   true,
}

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIDRaw := chi.URLParam(r, "eventID")

		searchQuery := r.URL.Query().Get("search")
		if len(searchQuery) > svc.Config.APIMaxSearchQueryLength {
			logrus.WithField("search_query_length", len(searchQuery)).Error("Search query is too long")
			http.Error(w, "Search query is too long", http.StatusBadRequest)
			return
		}

		filterState := api.ParseFilter(r.URL.Query(), "state", filterStateOptions)
		sortDirection := api.ParseSortDirection(r.URL.Query(), api.SortDirectionDesc)

		logrus.WithFields(logrus.Fields{
			"event_id":      eventIDRaw,
			"search":        searchQuery,
			"filterState":   filterState,
			"sortDirection": sortDirection,
			"query":         r.URL.RawQuery,
		}).Debug("Received listDeliveryAttempts request")

		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event ID")
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		if eventID <= 0 {
			logrus.Error("Event ID must be a positive integer")
			http.Error(w, "Event ID must be a positive integer", http.StatusBadRequest)
			return
		}

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
			eventID,
			cursor,
			pageSize,
			searchQuery,
			filterState,
			sortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}

		deliveryAttemptDTOs := mapperv1.ToDeliveryAttemptDTOs(deliveryAttempts)

		logrus.WithFields(logrus.Fields{
			"event_id":       eventID,
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
			logrus.WithError(err).Error("Failed to write response")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}
	}
}

// deliveryAttemptsRouter sets up the router for delivery attempts-related endpoints.
func deliveryAttemptsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listDeliveryAttempts(svc))
	return r
}
