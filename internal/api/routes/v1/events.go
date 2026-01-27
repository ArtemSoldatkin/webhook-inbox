package routev1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// getEventDetails sets up the GET event details endpoint for the API.
func getEventDetails(r chi.Router) {
	r.Get("/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET event details reached"))
	})
}

// eventsRouter sets up the router for event routes.
func eventsRouter() chi.Router {
	r := chi.NewRouter()
	getEventDetails(r)
	return r
}
