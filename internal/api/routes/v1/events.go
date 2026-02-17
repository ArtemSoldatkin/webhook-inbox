package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listEvents handles GET requests to list all events.
func listEvents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceIDRaw := chi.URLParam(r, "sourceID")
		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid source ID")
			http.Error(w, "Invalid source ID", http.StatusBadRequest)
			return
		}
		events, err := svc.ListEvents(r.Context(), sourceID)
		if err != nil {
			logrus.WithError(err).Error("Failed to list events")
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}
		eventDTOs := make([]dtov1.EventDTO, len(events))
		for i, event := range events {
			queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams); if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal query params")
				continue
			}
			rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders); if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal query params")
				continue
			}
			body, err := utils.JSONBtoType[map[string]string](event.Body); if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal query params")
				continue
			}
			remoteAddress := ""
			if event.RemoteAddress != nil {
				remoteAddress = event.RemoteAddress.String()
			}
			eventDTOs[i] = dtov1.EventDTO{
				ID:              event.ID,
				SourceID:        event.SourceID,
				DedupHash:       event.DedupHash.String,
				Method:          event.Method,
				IngressPath:     event.IngressPath,
				RemoteAddress:   remoteAddress,
				QueryParams:     queryParams,
				RawHeaders:      rawHeaders,
				Body:            body,
				BodyContentType: event.BodyContentType,
				ReceivedAt:      event.ReceivedAt.Time,
			}
		}
		response, err := json.Marshal(eventDTOs)
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
		queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams); if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal query params")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}
		rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders); if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal query params")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}
		body, err := utils.JSONBtoType[map[string]string](event.Body); if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal query params")
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}
		remoteAddress := ""
		if event.RemoteAddress != nil {
			remoteAddress = event.RemoteAddress.String()
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
			Body:            body,
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
	r.Mount("/{eventID}/delivery_attempts", deliveryAttemptsRouter(svc))
	r.Get("/{eventID}", getEvent(svc))
	r.Get("/", listEvents(svc))
	return r
}
