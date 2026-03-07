// Package routev1 contains API version 1 route handlers.
package routev1

import (
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// V1Router sets up and returns the router for API version 1.
func V1Router(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Use(httprate.LimitByIP(
		svc.Config.APIRateLimitRequests,
		time.Duration(svc.Config.APIRateLimitWindowSec)*time.Second,
	))
	r.Use(middleware.Throttle(svc.Config.APIThrottleConcurrentLimit))

	r.Mount("/ingest", ingestRouter(svc))
	r.Mount("/sources", sourcesRouter(svc))
	return r
}
