// Package routev1 contains API version 1 route handlers.
package routev1

import (
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// V1Router sets up and returns the router for API version 1.
func V1Router(service *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Mount("/users", usersRouter(service))
	r.Mount("/endpoints", endpointsRouter(service))
	r.Mount("/webhooks", webhooksRouter(service))
	r.Mount("/events", eventsRouter(service))
	r.Mount("/deliveries", deliveriesRouter(service))
	return r
}