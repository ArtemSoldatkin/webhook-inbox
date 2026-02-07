package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listEndpoints handles listing all endpoints for a given user.
func listEndpoints(r chi.Router, svc *service.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userIDRaw := r.URL.Query().Get("userID")
		if userIDRaw == "" {
			logrus.Error("Missing userID query parameter")
			http.Error(w, "userID query parameter is required", http.StatusBadRequest)
			return
		}
		userID, err := strconv.ParseInt(userIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid userID query parameter")
			http.Error(w, "Invalid userID query parameter", http.StatusBadRequest)
			return
		}
		endpoints, err := svc.ListEndpoints(r.Context(), userID)
		if err != nil {
			logrus.WithError(err).Error("Failed to list endpoints")
			http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
			return
		}
		endpointDTOs := make([]dtov1.Endpoint, len(endpoints))
		for i, endpoint := range endpoints {
			var data map[string]any
			err := json.Unmarshal(endpoint.Headers, &data)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal endpoint headers")
			}
			endpointDTOs[i] = dtov1.Endpoint{
				ID:          endpoint.ID,
				UserID:      endpoint.UserID,
				Url:         endpoint.Url,
				Name:        endpoint.Name,
				Description: endpoint.Description,
				Headers:     data,
				IsActive:    endpoint.IsActive,
				CreatedAt:   endpoint.CreatedAt,
			}
		}
		response, err := json.Marshal(endpointDTOs)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal endpoints")
			http.Error(w, "Failed to list endpoints", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// getEndpointByID handles retrieving a specific endpoint by its ID.
func getEndpointByID(r chi.Router, svc *service.Service) {
	r.Get("/{endpointID}", func(w http.ResponseWriter, r *http.Request) {
		endpointIDRaw := chi.URLParam(r, "endpointID")
		endpointID, err := strconv.ParseInt(endpointIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid endpoint ID")
			http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
			return
		}
		endpoint, err := svc.GetEndpointByID(r.Context(), endpointID)
		if err != nil {
			logrus.WithError(err).Error("Failed to get endpoint")
			http.Error(w, "Failed to get endpoint", http.StatusInternalServerError)
			return
		}
		var data map[string]any
		errHeaders := json.Unmarshal(endpoint.Headers, &data)
		if errHeaders != nil {
			logrus.WithError(errHeaders).Error("Failed to unmarshal endpoint headers")
			http.Error(w, "Failed to get endpoint", http.StatusInternalServerError)
			return
		}
		endpointDTO := dtov1.Endpoint{
			ID:          endpoint.ID,
			UserID:      endpoint.UserID,
			Url:         endpoint.Url,
			Name:        endpoint.Name,
			Description: endpoint.Description,
			Headers:     data,
			IsActive:    endpoint.IsActive,
			CreatedAt:   endpoint.CreatedAt,
		}
		response, err := json.Marshal(endpointDTO)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal endpoint")
			http.Error(w, "Failed to get endpoint", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// registerEndpointRequest represents the request payload for registering a new endpoint.
type registerEndpointRequest struct {
	UserID int64 `json:"userID"`
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
			logrus.WithError(err).Error("Invalid request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		endpoint, err := svc.RegisterEndpoint(r.Context(), req.UserID, req.Url, req.Name, req.Description, req.Headers)
		if err != nil {
			logrus.WithError(err).Error("Failed to register endpoint")
			http.Error(w, "Failed to register endpoint", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(endpoint)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal endpoint")
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
			logrus.WithError(err).Error("Invalid endpoint ID")
			http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
			return
		}
		endpoint, err := svc.ToggleEndpoint(r.Context(), endpointID)
		if err != nil {
			logrus.WithError(err).Error("Failed to toggle endpoint")
			http.Error(w, "Failed to toggle endpoint", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(endpoint)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal endpoint")
			http.Error(w, "Failed to toggle endpoint", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// endpointsRouter sets up the router for endpoints-related endpoints.
func endpointsRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	listEndpoints(router, svc)
	getEndpointByID(router, svc)
	registerEndpoint(router, svc)
	toggleEndpoint(router, svc)
	return router
}
