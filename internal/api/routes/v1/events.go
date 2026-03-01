package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/go-chi/chi/v5"
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
		logrus.
			WithFields(logrus.Fields{
				"limit":  pageSize,
				"cursor": cursor,
			}).
			Info("Listing events with pagination")
		sourceIDRaw := chi.URLParam(r, "sourceID")
		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid source ID")
			http.Error(w, "Invalid source ID", http.StatusBadRequest)
			return
		}
		events, err := svc.ListEvents(r.Context(), sourceID, cursor, pageSize)
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
				continue
			}
			rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal raw headers")
				continue
			}
			var remoteAddress *string
			if event.RemoteAddress != nil {
				str := event.RemoteAddress.String()
				remoteAddress = &str
			}
			eventDTOs = append(eventDTOs, dtov1.EventDTO{
				ID:       event.ID,
				SourceID: event.SourceID,
				// DedupHash:       event.DedupHash.String,
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
		paginatedResponse := api.ToPaginatedResponse(eventDTOs, pageSize, func(e dtov1.EventDTO) *time.Time {
			return &e.ReceivedAt
		})
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
		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event ID")
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}
		event, err := svc.GetEventByID(r.Context(), eventID)
		if err != nil {
			logrus.WithError(err).Error("Failed to get event")
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
			ID:       event.ID,
			SourceID: event.SourceID,
			// DedupHash:       event.DedupHash.String,
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
