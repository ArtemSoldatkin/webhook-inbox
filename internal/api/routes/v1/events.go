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
		if err != nil {
			logrus.WithError(err).Error("Invalid source_id query parameter")
			http.Error(w, "Invalid source_id query parameter", http.StatusBadRequest)
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
			queryParams, queryParamsErr := utils.JSONBtoMap(event.QueryParams); if queryParamsErr != nil {
				logrus.WithError(queryParamsErr).Error("Failed to unmarshal query params")
				continue
			}
			// rawHeaders, rawHeadersErr := utils.JSONBtoMap(event.RawHeaders); if rawHeadersErr != nil {
			// 	logrus.WithError(rawHeadersErr).Error("Failed to unmarshal query params")
			// 	continue
			// }
			body, bodyErr := utils.JSONBtoMap(event.Body); if bodyErr != nil {
				logrus.WithError(bodyErr).Error("Failed to unmarshal query params")
				continue
			}
			eventDTOs[i] = dtov1.EventDTO{
				ID:              event.ID,
				SourceID:        event.SourceID,
				DedupHash:       event.DedupHash.String,
				Method:          event.Method,
				IngressPath:     event.IngressPath,
				RemoteAddress:   event.RemoteAddress.String(),
				QueryParams:     queryParams,
				RawHeaders:      map[string]string{}, // TODO: unmarshal raw headers
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

// eventsRouter sets up the router for events-related endpoints.
func eventsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listEvents(svc))
	return r
}
