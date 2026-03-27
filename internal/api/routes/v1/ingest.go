package routev1

import (
	"errors"
	"net/http"

	requestsv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/requests/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// ingestEvent handles ANY requests to ingest a new event.
func ingestEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[requestsv1.IngestEventInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"public_id": input.PublicID,
			"method":    r.Method,
			"path":      r.URL.Path,
			"query":     r.URL.RawQuery,
		}).Debug("Received ingestEvent request")

		if err := requestsv1.ValidateIngestEventInput(input); err != nil {
			logrus.WithError(err).WithField("public_id", input.PublicID).Error("Input validation failed")
			http.Error(w, "Invalid input parameters: "+err.Error(), http.StatusBadRequest)
			return
		}

		eventID, err := svc.CreateEvent(r.Context(), r, input.PublicID)
		if err != nil {
			var notFoundErr *service.SourceIsNotFound
			if errors.As(err, &notFoundErr) {
				logrus.WithError(err).Error("Source not found for given public_id")
				http.Error(w, notFoundErr.Message, http.StatusNotFound)
				return
			}
			logrus.WithError(err).Error("Failed to create event")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"event_id": eventID,
		}).Info("Created event")

		deliveryAttemptID, err := svc.CreateDeliveryAttempt(
			r.Context(),
			db.CreateDeliveryAttemptParams{
				EventID:       eventID,
				AttemptNumber: 1,
				State:         "pending",
			},
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to create delivery attempt")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"event_id":            eventID,
			"delivery_attempt_id": deliveryAttemptID,
			"query":               r.URL.RawQuery,
		}).Info("Created initial delivery attempt")

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK"))
	}
}

// ingestRouter sets up the router for event ingestion endpoints.
func ingestRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/{public_id}", ingestEvent(svc))
	return r
}
