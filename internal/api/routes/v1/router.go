// Package routev1 contains API version 1 route handlers.
package routev1

import (
	"fmt"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/go-chi/cors"

	"github.com/unrolled/secure"
)

// V1Router sets up and returns the router for API version 1.
func V1Router(svc *service.Service) chi.Router {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
	})

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			fmt.Sprintf("%s://%s:%d", svc.Config.UIProtocol, svc.Config.UIHost, svc.Config.UIPort),
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         svc.Config.APICORSMaxAgeSec,
	}))

	r.Use(httprate.LimitByIP(
		svc.Config.APIRateLimitRequests,
		time.Duration(svc.Config.APIRateLimitWindowSec)*time.Second,
	))
	r.Use(middleware.Throttle(svc.Config.APIThrottleConcurrentLimit))
	r.Use(middleware.RequestSize(svc.Config.APIRequestSizeLimitBytes))

	r.Use(secureMiddleware.Handler)

	r.Mount("/ingest", ingestRouter(svc))
	r.Mount("/sources", sourcesRouter(svc))
	return r
}
