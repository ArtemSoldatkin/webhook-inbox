package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// listEvents handles listing all events for a given webhook.
func listEvents(r chi.Router, svc *service.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		webhookIDRaw := r.URL.Query().Get("webhook_id")
		if webhookIDRaw == "" {
			http.Error(w, "webhook_id query parameter is required", http.StatusBadRequest)
			return
		}
		webhookID, err := strconv.ParseInt(webhookIDRaw, 10, 64)
		if err != nil {
			http.Error(w, "Invalid webhook_id query parameter", http.StatusBadRequest)
			return
		}
		events, err := svc.Events.ListEvents(webhookID)
		if err != nil {
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(events)
		if err != nil {
			http.Error(w, "Failed to list events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}


// deliveriesRouter sets up the router for deliveries-related endpoints.
func eventsRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	listEvents(router, svc)
	return router
}