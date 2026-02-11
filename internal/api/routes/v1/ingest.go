package routev1

import (
	"net/http"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// ingestEvent handles POST requests to ingest a new event.
func ingestEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		publicID := chi.URLParam(r, "public_id")
		if publicID == "" {
			http.Error(w, "public_id is required", http.StatusBadRequest)
			return
		}
		logrus.Infof("Received event for public_id: %s", publicID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	}}

// ingestRouter sets up the router for event ingestion endpoints.
func ingestRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Post("/{public_id}", ingestEvent(svc))
	return r
}
