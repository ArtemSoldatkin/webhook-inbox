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
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// listEvents handles GET requests to list all events.
func listEvents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.ListEventsInput](r)
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
			"source_id":      input.SourceID,
			"page_size":      input.PageSize,
			"cursor":         input.Cursor,
			"search":         input.Search,
			"sort_direction": input.SortDirection,
			"query":          r.URL.RawQuery,
		}).Debug("Received listEvents request")

		events, err := svc.ListEvents(
			r.Context(),
			input.SourceID,
			input.Cursor,
			input.PageSize,
			input.Search,
			input.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list events")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to list events"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		eventDTOs := mapperv1.ToEventDTOs(events)

		logrus.WithFields(logrus.Fields{
			"source_id":      input.SourceID,
			"returned_count": len(eventDTOs),
		}).Debug("Returning events")

		var nextCursor types.Cursor
		if len(eventDTOs) > input.PageSize {
			lastEvent := eventDTOs[input.PageSize-1]
			nextCursor = types.NewCursor(
				&lastEvent.ReceivedAt,
				&lastEvent.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			eventDTOs,
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
				map[string]string{"error": "Failed to list events"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
		}
	}
}

// getEvent handles GET requests to retrieve an event by its ID.
func getEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.GetEventInput](r)
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
			"event_id": input.EventID,
			"query":    r.URL.RawQuery,
		}).Debug("Received getEvent request")

		event, err := svc.GetEventByID(r.Context(), input.EventID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("event_id", input.EventID).Info("Event not found")
				if err := api.JSON(
					w,
					http.StatusNotFound,
					map[string]string{"error": "Event not found"},
				); err != nil {
					logrus.WithError(err).Error("Failed to write error response")
				}
				return
			}
			logrus.WithField("event_id", input.EventID).WithError(err).Error("Failed to get event")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to get event"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
			return
		}

		eventDTO := mapperv1.ToEventDTO(event)

		if err := api.JSON(w, http.StatusOK, eventDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
				return
			}

			logrus.WithError(err).Error("Failed to marshal response")
			if err := api.JSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "Failed to get event"},
			); err != nil {
				logrus.WithError(err).Error("Failed to write error response")
			}
		}
	}
}

// eventsRouter sets up the router for events-related endpoints.
func eventsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Mount("/{event_id}/delivery-attempts", deliveryAttemptsRouter(svc))
	r.Get("/{event_id}", getEvent(svc))
	r.Get("/", listEvents(svc))
	return r
}
