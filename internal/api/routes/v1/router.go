// Package routev1 contains API version 1 route handlers.
package routev1

import (
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// V1Router sets up and returns the router for API version 1.
func V1Router(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Mount("/ingest", ingestRouter(svc))
	r.Mount("/sources", sourcesRouter(svc))
	return r
}
