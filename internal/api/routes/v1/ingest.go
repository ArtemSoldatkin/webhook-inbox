package routev1

import (
	"encoding/json"
	"net"
	"net/http"
	"net/netip"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

// ingestEvent handles ANY requests to ingest a new event.
func ingestEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		publicID := chi.URLParam(r, "public_id")
		if publicID == "" {
			http.Error(w, "public_id is required", http.StatusBadRequest)
			return
		}
		logrus.Infof("Received event for public_id: %s", publicID)
		source, err := svc.GetSourceByPublicID(r.Context(), publicID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to retrieve source for public_id: %s", publicID)
			http.Error(w, "Source not found", http.StatusNotFound)
			return
		}
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to split host and port from remote address: %s", r.RemoteAddr)
			http.Error(w, "Invalid remote address", http.StatusBadRequest)
			return
		}
		remoteAddress, err := netip.ParseAddr(host)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to parse remote address: %s", r.RemoteAddr)
			http.Error(w, "Invalid remote address", http.StatusBadRequest)
			return
		}
		headerBytes, err := json.Marshal(r.Header)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal headers")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		queryParams, err := json.Marshal(r.URL.Query())
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal query parameters")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		bodyBytes, err := json.Marshal(r.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read request body")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		eventID, err := svc.CreateEvent(r.Context(), db.CreateEventParams{
			SourceID: source.ID,
			DedupHash: pgtype.Text{String: "", Valid: false}, // TODO replace with actual deduplication logic
			Method: r.Method,
			IngressPath: r.URL.Path,
			RemoteAddress: &remoteAddress,
			QueryParams: queryParams,
			RawHeaders: headerBytes,
			Body: bodyBytes,
			BodyContentType: r.Header.Get("Content-Type"),
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create event")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logrus.Infof("Successfully ingested event with ID: %d", eventID)
		deliveryAttemptID, err := svc.CreateDeliveryAttempt(r.Context(), db.CreateDeliveryAttemptParams{
			EventID: eventID,
			AttemptNumber: 1,
			State: "pending",
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create delivery attempt")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logrus.Infof("Created initial delivery attempt with ID: %d for event ID: %d", deliveryAttemptID, eventID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	}}

// ingestRouter sets up the router for event ingestion endpoints.
func ingestRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/{public_id}", ingestEvent(svc))
	return r
}
