package routev1

import (
	"errors"
	"net/http"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	requestsv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/requests/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.ListDeliveryAttemptsInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			if err := api.JSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Invalid input parameters"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
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
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to list delivery attempts"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		deliveryAttemptDTOs := mapperv1.ToDeliveryAttemptDTOs(deliveryAttempts)

		logrus.WithFields(logrus.Fields{
			"event_id":       input.EventID,
			"returned_count": len(deliveryAttemptDTOs),
		}).Debug("Returning delivery attempts")

		var nextCursor types.Cursor
		if len(deliveryAttemptDTOs) > input.PageSize {
			lastAttempt := deliveryAttemptDTOs[input.PageSize-1]
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
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to list delivery attempts"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
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
