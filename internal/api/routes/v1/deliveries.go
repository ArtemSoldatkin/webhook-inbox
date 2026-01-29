package routev1

import (
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
)

// deliveriesRouter sets up the router for deliveries-related endpoints.
func deliveriesRouter(svc *service.Service) chi.Router {
	router := chi.NewRouter()
	return router
}