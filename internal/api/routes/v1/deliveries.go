package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// listDeliveries handles GET requests to list deliveries for a specific endpoint.
func listDeliveries(r chi.Router, svc *service.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		endpointIDRaw := r.URL.Query().Get("endpoint_id")
		if endpointIDRaw == "" {
			http.Error(w, "endpoint_id query parameter is required", http.StatusBadRequest)
			return
		}
		endpointID, err := strconv.ParseInt(endpointIDRaw, 10, 64)
		if err != nil {
			http.Error(w, "Invalid endpoint_id query parameter", http.StatusBadRequest)
			return
		}
		deliveries, err := svc.Deliveries.ListDeliveries(endpointID)
		if err != nil {
			http.Error(w, "Failed to list deliveries", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(deliveries)
		if err != nil {
			http.Error(w, "Failed to list deliveries", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}


// deliveriesRouter sets up the router for deliveries-related endpoints.
func deliveriesRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	listDeliveries(router, svc)
	return router
}