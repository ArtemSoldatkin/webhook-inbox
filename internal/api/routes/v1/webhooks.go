package routev1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// postWebhook sets up the POST webhook endpoint for the API.
func postWebhook(r chi.Router) {
	r.Post("/{endpointId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POST webhook received"))
	})
}

// webhooksRouter sets up the router for webhook endpoints.
func webhooksRouter() chi.Router {
	r := chi.NewRouter()
	postWebhook(r)
	return r
}