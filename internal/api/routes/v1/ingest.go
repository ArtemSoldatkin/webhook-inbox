package routev1

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"regexp"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

var uuidRegexp = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)

// IngestEventInput defines the expected input parameters for ingesting a new event.
type IngestEventInput struct {
	PublicID string `url_param:"public_id,required"`
}

// ingestEvent handles ANY requests to ingest a new event.
func ingestEvent(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[IngestEventInput](r)
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

		if !validatePublicID(input.PublicID) {
			logrus.WithField("public_id", input.PublicID).Error("Invalid public_id")
			http.Error(w, "Invalid public_id", http.StatusBadRequest)
			return
		}

		source, err := svc.GetSourceByPublicID(r.Context(), input.PublicID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("public_id", input.PublicID).Info("Source not found")
				http.Error(w, "Source not found", http.StatusNotFound)
				return
			}
			logrus.WithField("public_id", input.PublicID).WithError(err).Error("Failed to get source")
			http.Error(w, "Failed to get source", http.StatusInternalServerError)
			return
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logrus.
				WithError(err).
				Errorf("Failed to split host and port from remote address: %s", r.RemoteAddr)
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

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read request body")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		method := r.Method
		ingressPath := r.URL.Path

		dedupHash, err := generateDedupHash(DedupPayload{
			Method:      method,
			IngressPath: ingressPath,
			QueryParams: queryParams,
			RawHeaders:  headerBytes,
			Body:        bodyBytes,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to generate deduplication hash")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		eventID, err := svc.CreateEvent(r.Context(), db.CreateEventParams{
			SourceID:        source.ID,
			DedupHash:       pgtype.Text{String: dedupHash, Valid: true},
			Method:          method,
			IngressPath:     ingressPath,
			RemoteAddress:   &remoteAddress,
			QueryParams:     queryParams,
			RawHeaders:      headerBytes,
			Body:            bodyBytes,
			BodyContentType: r.Header.Get("Content-Type"),
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create event")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"event_id":   eventID,
			"body_bytes": len(bodyBytes),
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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	}
}

// DedupPayload represents the data used to generate a deduplication hash for an event.
type DedupPayload struct {
	Method      string `json:"method"`
	IngressPath string `json:"ingress_path"`
	QueryParams []byte `json:"query_params"`
	RawHeaders  []byte `json:"headers"`
	Body        []byte `json:"body"`
}

// generateDedupHash generates a SHA-256 hash of the input data for deduplication purposes.
func generateDedupHash(dedupPayload DedupPayload) (string, error) {
	dedupData, err := json.Marshal(dedupPayload)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(dedupData)
	return hex.EncodeToString(hash[:]), nil
}

// validatePublicID checks if the provided public ID matches the expected UUID format.
func validatePublicID(publicID string) bool {
	return uuidRegexp.MatchString(publicID)
}

// ingestRouter sets up the router for event ingestion endpoints.
func ingestRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.HandleFunc("/{public_id}", ingestEvent(svc))
	return r
}
