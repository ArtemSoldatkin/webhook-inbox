package routev1

import (
	"encoding/json"
	"net/http"
	"time"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

// listSources handles GET requests to list all sources.
func listSources(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sources, err := svc.ListSources(r.Context())
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}
		sourceDTOs := make([]dtov1.SourceDTO, len(sources))
		for i, source := range sources {
			var staticHeaders = make(map[string]string)
			if err := json.Unmarshal(source.StaticHeaders, &staticHeaders); err != nil {
				logrus.WithError(err).Error("Failed to unmarshal static headers")
				http.Error(w, "Failed to list sources", http.StatusInternalServerError)
				return
			}
			if staticHeaders == nil {
				staticHeaders = make(map[string]string)
			}
			var disbaleAt *time.Time
			if source.DisableAt.Valid {
				disbaleAt = &source.DisableAt.Time
			} else {
				disbaleAt = nil
			}
			sourceDTOs[i] = dtov1.SourceDTO{
				ID:             source.ID,
				IngressUrl:     source.IngressUrl,
				EgressUrl:      source.EgressUrl,
				StaticHeaders:  staticHeaders,
				Status:         source.Status,
				StatusReason:   source.StatusReason.String,
				Description:    source.Description.String,
				CreatedAt:      source.CreatedAt.Time,
				UpdatedAt:      source.UpdatedAt.Time,
				DisableAt:      disbaleAt,
			}
		}
		response, err := json.Marshal(sourceDTOs)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// CreateSourceData defines the parameters required to create a new source.
type CreateSourceData struct {
	IngressUrl           string            `json:"IngressUrl"`
	EgressUrl            string            `json:"EgressUrl"`
	StaticHeaders map[string]string `json:"staticHeaders,omitempty"`
	Description   string            `json:"description,omitempty"`
}

// createSource handles POST requests to create a new source.
func createSource(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data CreateSourceData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logrus.WithError(err).Error("Failed to decode create source request")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		staticHeaders, staticHeadersErr := json.Marshal(data.StaticHeaders)
		if staticHeadersErr != nil {
			logrus.WithError(staticHeadersErr).Error("Failed to marshal static headers")
			http.Error(w, "Invalid static headers", http.StatusBadRequest)
			return
		}
		source, err := svc.CreateSource(r.Context(), db.CreateSourceParams{
			IngressUrl:		data.IngressUrl,
			EgressUrl:		data.EgressUrl,
			StaticHeaders:  staticHeaders,
			Description: 	pgtype.Text{String: data.Description},
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			http.Error(w, "Failed to create source", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(source)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal created source")
			http.Error(w, "Failed to create source", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

// sourcesRouter sets up the router for sources-related endpoints.
func sourcesRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listSources(svc))
	r.Post("/", createSource(svc))
	return r
}
