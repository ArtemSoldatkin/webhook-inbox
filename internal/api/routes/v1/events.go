package routev1

import (
	"errors"
	"net/http"
	"strconv"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// listEvents handles GET requests to list all events.
func listEvents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		sourceIDRaw := chi.URLParam(r, "sourceID")

		searchQuery := r.URL.Query().Get("search")
		if len(searchQuery) > svc.Config.APIMaxSearchQueryLength {
			logrus.WithField("search_query_length", len(searchQuery)).Error("Search query is too long")
			http.Error(w, "Search query is too long", http.StatusBadRequest)
			return
		}

		sortDirection := api.ParseSortDirection(r.URL.Query(), api.SortDirectionDesc)

		logrus.WithFields(logrus.Fields{
			"source_id":     sourceIDRaw,
			"pageSize":      pageSize,
			"cursor":        cursor,
			"search":        searchQuery,
			"sortDirection": sortDirection,
			"query":         r.URL.RawQuery,
		}).Debug("Received listEvents request")

		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid source ID")
			http.Error(w, "Invalid source ID", http.StatusBadRequest)
			return
		}

		if sourceID <= 0 {
			logrus.WithField("source_id", sourceID).Error("Source ID must be a positive integer")
			http.Error(w, "Source ID must be a positive integer", http.StatusBadRequest)
			return
		}

		events, err := svc.ListEvents(
			r.Context(),
			sourceID,
			cursor,
			pageSize,
			searchQuery,
			sortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list events")
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		eventDTOs := mapperv1.ToEventDTOs(events)

		logrus.WithFields(logrus.Fields{
			"source_id":      sourceID,
			"returned_count": len(eventDTOs),
		}).Debug("Returning events")

		var nextCursor api.Cursor
		if len(eventDTOs) > pageSize {
			lastEvent := eventDTOs[len(eventDTOs)-1]
			nextCursor = api.NewCursor(
				&lastEvent.ReceivedAt,
				&lastEvent.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			eventDTOs,
			pageSize,
			nextCursor,
		)

		if err := api.JSON(w, http.StatusOK, paginatedResponse); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to list events", http.StatusInternalServerError)
			}
		}
	}
}

// getEvent handles GET requests to retrieve an event by its ID.
func getEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIDRaw := chi.URLParam(r, "eventID")

		logrus.WithFields(logrus.Fields{
			"event_id": eventIDRaw,
			"query":    r.URL.RawQuery,
		}).Debug("Received getEvent request")

		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event ID")
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		if eventID <= 0 {
			logrus.WithField("event_id", eventID).Error("Event ID must be a positive integer")
			http.Error(w, "Event ID must be a positive integer", http.StatusBadRequest)
			return
		}

		event, err := svc.GetEventByID(r.Context(), eventID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("event_id", eventID).Info("Event not found")
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}
			logrus.WithField("event_id", eventID).WithError(err).Error("Failed to get event")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		eventDTO := mapperv1.ToEventDTO(event)

		if err := api.JSON(w, http.StatusOK, eventDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to get event", http.StatusInternalServerError)
			}
		}
	}
}

// eventsRouter sets up the router for events-related endpoints.
func eventsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Mount("/{eventID}/delivery-attempts", deliveryAttemptsRouter(svc))
	r.Get("/{eventID}", getEvent(svc))
	r.Get("/", listEvents(svc))
	return r
}
