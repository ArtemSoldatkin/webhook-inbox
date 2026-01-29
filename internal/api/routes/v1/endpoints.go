package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// listEndpoints handles listing all endpoints for a given user.
func listEndpoints(r chi.Router, svc *service.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userIDRaw := r.URL.Query().Get("user_id")
		if userIDRaw == "" {
			http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
			return
		}
		userID, err := strconv.ParseInt(userIDRaw, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user_id query parameter", http.StatusBadRequest)
			return
		}
		endpoints, err := svc.Endpoints.ListEndpoints(userID)
		if err != nil {
			http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(endpoints)
		if err != nil {
			http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// registerEndpointRequest represents the request payload for registering a new endpoint.
type registerEndpointRequest struct {
	UserID int64 `json:"user_id"`
	Url	string `json:"url"`
	Name   string `json:"name"`
	Description string `json:"description"`
	Headers map[string]string `json:"headers"`
}


// registerEndpoint handles the registration of a new endpoint.
func registerEndpoint(r chi.Router, svc *service.Service) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req registerEndpointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		user, err := svc.Endpoints.RegisterEndpoint(req.UserID, req.Url, req.Name, req.Description, req.Headers)
		if err != nil {
			http.Error(w, "Failed to register endpoint", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Failed to register endpoint", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

// toggleEndpoint handles toggling the active status of an endpoint.
func toggleEndpoint(r chi.Router, svc *service.Service) {
	r.Put("/{endpointID}/toggle", func(w http.ResponseWriter, r *http.Request) {
		endpointIDRaw := chi.URLParam(r, "endpointID")
		endpointID, err := strconv.ParseInt(endpointIDRaw, 10, 64)
		if err != nil {
			http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
			return
		}
		endpoint, err := svc.Endpoints.ToggleEndpoint(endpointID)
		if err != nil {
			http.Error(w, "Failed to toggle endpoint", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(endpoint)
		if err != nil {
			http.Error(w, "Failed to toggle endpoint", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// deliveriesRouter sets up the router for deliveries-related endpoints.
func endpointsRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	listEndpoints(router, svc)
	registerEndpoint(router, svc)
	toggleEndpoint(router, svc)
	return router
}