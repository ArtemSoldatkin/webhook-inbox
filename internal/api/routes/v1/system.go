package routev1

import (
	"net/http"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// health handles GET requests to check the health of the API.
func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := api.JSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
			logrus.WithError(err).Error("Failed to write health check response")
		}
	}
}

// systemRouter sets up the router for health check endpoint.
func systemRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", health())
	return r
}
