package routev1

import (
	"encoding/json"
	"net/http"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// listSources handles GET requests to list all sources.
func listSources(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sources, err := svc.ListSources(r.Context())
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(sources)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// sourcesRouter sets up the router for sources-related endpoints.
func sourcesRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listSources(svc))
	return r
}
