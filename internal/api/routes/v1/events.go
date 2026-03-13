package routev1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
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
			"source_id":      sourceIDRaw,
			"pageSize":       pageSize,
			"cursor":         cursor,
			"search":         searchQuery,
			"sortDirection": sortDirection,
			"query":          r.URL.RawQuery,
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

		eventDTOs := make([]dtov1.EventDTO, 0, len(events))
		for _, event := range events {
			queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal query params")
				queryParams = map[string][]string{
					"__error": {"Webhook Inbox Error - Failed to parse"},
				}
			}

			rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal raw headers")
				rawHeaders = map[string][]string{
					"__error": {"Webhook Inbox Error - Failed to parse"},
				}
			}

			var remoteAddress *string
			if event.RemoteAddress != nil {
				str := event.RemoteAddress.String()
				remoteAddress = &str
			}

			eventDTOs = append(eventDTOs, dtov1.EventDTO{
				ID:              event.ID,
				SourceID:        event.SourceID,
				DedupHash:       event.DedupHash.String,
				Method:          event.Method,
				IngressPath:     event.IngressPath,
				RemoteAddress:   remoteAddress,
				QueryParams:     queryParams,
				RawHeaders:      rawHeaders,
				Body:            event.Body,
				BodyContentType: event.BodyContentType,
				ReceivedAt:      event.ReceivedAt.Time,
			})
		}

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

		response, err := json.Marshal(paginatedResponse)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal events")
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
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

		queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal query params")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal raw headers")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		var remoteAddress *string
		if event.RemoteAddress != nil {
			str := event.RemoteAddress.String()
			remoteAddress = &str
		}

		eventDTO := dtov1.EventDTO{
			ID:              event.ID,
			SourceID:        event.SourceID,
			DedupHash:       event.DedupHash.String,
			Method:          event.Method,
			IngressPath:     event.IngressPath,
			RemoteAddress:   remoteAddress,
			QueryParams:     queryParams,
			RawHeaders:      rawHeaders,
			Body:            event.Body,
			BodyContentType: event.BodyContentType,
			ReceivedAt:      event.ReceivedAt.Time,
		}

		response, err := json.Marshal(eventDTO)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal event")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
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
