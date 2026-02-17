package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listDeliveryAttempts handles GET requests to list all delivery attempts.
func listDeliveryAttempts(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIDRaw := chi.URLParam(r, "eventID")
		eventID, err := strconv.ParseInt(eventIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid event ID")
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
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
	r.Get("/", listDeliveryAttempts(svc))
	return r
}
