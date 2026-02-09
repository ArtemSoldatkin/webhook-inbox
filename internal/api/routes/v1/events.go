package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listEvents handles GET requests to list all events.
func listEvents(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceIDRaw := r.URL.Query().Get("source_id")
		if sourceIDRaw == "" {
			logrus.Error("Missing source_id query parameter")
			http.Error(w, "source_id query parameter is required", http.StatusBadRequest)
			return
		}
		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
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
		response, err := json.Marshal(events)
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
