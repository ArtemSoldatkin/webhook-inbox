package routev1

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// createUserRequest represents the expected payload for creating a new user.
type createUserRequest struct {
	Email    string `json:"email"`
}

// createUser handles the creation of a new user.
func createUser(r chi.Router, svc *service.Service) {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logrus.WithError(err).Error("Invalid request body")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		user, err := svc.CreateUser(r.Context(), req.Email)
		if err != nil {
			logrus.WithError(err).Error("Failed to create user")
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(user)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal user")
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

// usersRouter sets up the router for users-related endpoints.
func usersRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	createUser(router, svc)
	return router
}
