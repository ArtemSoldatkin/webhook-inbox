package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listDevliveryAttempts handles GET requests to list all delivery attempts.
func listDevliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIDRaw := r.URL.Query().Get("event_id")
		if eventIDRaw == "" {
			logrus.Error("Missing event_id query parameter")
			http.Error(w, "event_id query parameter is required", http.StatusBadRequest)
			return
		}
		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event_id query parameter")
			http.Error(w, "Invalid event_id query parameter", http.StatusBadRequest)
			return
		}
		deliveryAttempts, err := svc.ListDeliveryAttempts(r.Context(), eventID)
		if err != nil {
			logrus.WithError(err).Error("Failed to list delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(deliveryAttempts)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal delivery attempts")
			http.Error(w, "Failed to list delivery attempts", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// deliveryAttemptsRouter sets up the router for delivery attempts-related endpoints.
func deliveryAttemptsRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listDevliveryAttempts(svc))
	return r
}
