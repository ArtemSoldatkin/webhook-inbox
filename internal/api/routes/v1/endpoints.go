package routev1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// postEndpointRequest represents the expected structure of a POST endpoint request.
type postEndpointRequest struct {
	Name string `json:"name"`
	Description *string `json:"description"`
}

// postEndpoints sets up the POST endpoints for the API.
func postEndpoints(r chi.Router, service *service.Service) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req postEndpointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		endpoint, err := service.CreateEndpoint(req.Name, req.Description)
		if err != nil {
			log.Println("Error creating endpoint:", err)
			http.Error(w, "Failed to create endpoint", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(endpoint)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

// getEvents sets up the GET events endpoint for the API.
func getEvents(r chi.Router) {
	r.Get("/{endpointId}/events", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET events reached"))
	})
}


// endpointsRouter sets up the router for endpoint routes.
func endpointsRouter(service *service.Service) chi.Router {
	r := chi.NewRouter()
	postEndpoints(r, service)
	getEvents(r)
	return r
}
