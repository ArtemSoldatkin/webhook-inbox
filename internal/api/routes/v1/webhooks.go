package routev1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listWebhooks handles listing all webhooks for a given endpoint.
func listWebhooks(r chi.Router, svc *service.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		endpointIDRaw := r.URL.Query().Get("endpointID")
		endpointID, err := strconv.ParseInt(endpointIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid endpoint ID")
			http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
			return
		}
		webhooks, err := svc.ListWebhooks(r.Context(), endpointID)
		if err != nil {
			logrus.WithError(err).Error("Failed to list webhooks")
			http.Error(w, "Failed to list webhooks", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(webhooks)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal webhooks")
			http.Error(w, "Failed to list webhooks", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// createWebhookRequest represents the expected payload for creating a new webhook.
type createWebhookRequest struct {
	EndpointID int64 `json:"endpoint_id"`
	Name   string `json:"name"`
	Description string `json:"description"`
}

// createWebhook sets up the route for creating a new webhook.
func createWebhook(r chi.Router, svc *service.Service) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req createWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logrus.WithError(err).Error("Invalid request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// TODO add public key generation
		user, err := svc.CreateWebhook(r.Context(), req.EndpointID, req.Name, req.Description)
		if err != nil {
			logrus.WithError(err).Error("Failed to create webhook")
			http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(user)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal webhook")
			http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

// toggleWebhook handles toggling the active status of a webhook.
func toggleWebhook(r chi.Router, svc *service.Service) {
	r.Put("/{webhookID}/toggle", func(w http.ResponseWriter, r *http.Request) {
		webhookIDRaw := chi.URLParam(r, "webhookID")
		webhookID, err := strconv.ParseInt(webhookIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid webhook ID")
			http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
			return
		}
		webhook, err := svc.ToggleWebhook(r.Context(), webhookID)
		if err != nil {
			logrus.WithError(err).Error("Failed to toggle webhook")
			http.Error(w, "Failed to toggle webhook", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(webhook)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal webhook")
			http.Error(w, "Failed to toggle webhook", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

// webhooksRouter sets up the router for webhooks-related endpoints.
func webhooksRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	listWebhooks(router, svc)
	createWebhook(router, svc)
	toggleWebhook(router, svc)
	return router
}
