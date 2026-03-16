package routev1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

const (
	maxDescriptionLen = 500
	maxHeaders        = 20
	maxHeaderKeyLen   = 100
	maxHeaderValueLen = 500
)

var (
	httpRegexp        = regexp.MustCompile(`^https?://`)
	localhostRegexp   = regexp.MustCompile(`^https?://(localhost|127\.0\.0\.1|0\.0\.0\.0|\[?::1\]?)(/|:|$)`)
	private10Regexp   = regexp.MustCompile(`^https?://10\.`)
	private192Regexp  = regexp.MustCompile(`^https?://192\.168\.`)
	private172Regexp  = regexp.MustCompile(`^https?://172\.(1[6-9]|2[0-9]|3[0-1])\.`)
	metadata169Regexp = regexp.MustCompile(`^https?://169\.254\.169\.254(/|:|$)`)
)

// ListSourcesInput defines the expected input parameters for listing sources.
type ListSourcesInput struct {
	Filter        string       `query_param:"filter_state,allowed=active|paused|quarantined|disabled,default=*"`
	SortDirection string       `query_param:"sort,allowed=ASC|DESC,default=DESC"`
	Search        string       `query_param:"search,max_length=512"`
	PageSize      int          `query_param:"limit,min=1,max=100,default=20"`
	Cursor        types.Cursor `query_param:"cursor"`
}

// listSources handles GET requests to list all sources.
func listSources(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[ListSourcesInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"page_size":      input.PageSize,
			"cursor":         input.Cursor,
			"search":         input.Search,
			"filter_status":  input.Filter,
			"sort_direction": input.SortDirection,
			"query":          r.URL.RawQuery,
		}).Debug("Received listSources request")

		sources, err := svc.ListSources(
			r.Context(),
			input.Cursor,
			input.PageSize,
			input.Search,
			input.Filter,
			input.SortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}

		sourceDTOs := mapperv1.ToSourceDTOs(sources, svc.Config)

		logrus.WithField("returned_count", len(sourceDTOs)).Debug("Returning sources")

		var nextCursor types.Cursor
		if len(sourceDTOs) > input.PageSize {
			lastSource := sourceDTOs[len(sourceDTOs)-1]
			nextCursor = types.NewCursor(
				&lastSource.UpdatedAt,
				&lastSource.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			sourceDTOs,
			input.PageSize,
			nextCursor,
		)

		if err := api.JSON(w, http.StatusOK, paginatedResponse); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			}
		}
	}
}

// GetSourceByIDInput defines the expected input parameters for retrieving a source by its ID.
type GetSourceByIDInput struct {
	SourceID int64 `url_param:"event_id,required,min=1"`
}

// getSourceByID handles GET requests to retrieve a source by its ID.
func getSourceByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input, err := api.ParseRequestInput[GetSourceByIDInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"source_id": input.SourceID,
			"query":     r.URL.RawQuery,
		}).Debug("Received getSourceByID request")

		source, err := svc.GetSourceByID(r.Context(), input.SourceID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("source_id", input.SourceID).Info("Source not found")
				http.Error(w, "Source not found", http.StatusNotFound)
				return
			}
			logrus.WithField("source_id", input.SourceID).WithError(err).Error("Failed to get source")
			http.Error(w, "Failed to get source", http.StatusInternalServerError)
			return
		}

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusOK, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to get source", http.StatusInternalServerError)
			}
		}
	}
}

// CreateSourceData defines the expected input parameters for creating a new source.
type CreateSourceInput struct {
	EgressUrl     string            `json:"EgressUrl"`
	StaticHeaders map[string]string `json:"StaticHeaders,omitempty"`
	Description   string            `json:"Description,omitempty"`
}

// createSource handles POST requests to create a new source.
func createSource(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"query": r.URL.RawQuery,
		}).Debug("Received createSource request")

		input, err := api.ParseRequestInput[CreateSourceInput](r)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse input parameters")
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"egress_url":  input.EgressUrl,
			"description": input.Description,
		}).Debug("Create source request data")

		if len(input.Description) > maxDescriptionLen {
			logrus.Error("Description too long")
			http.Error(
				w,
				fmt.Sprintf("Description must be %d characters or less", maxDescriptionLen),
				http.StatusBadRequest,
			)
			return
		}

		if len(input.StaticHeaders) > maxHeaders {
			http.Error(w, "Too many headers", http.StatusBadRequest)
			return
		}
		for k, v := range input.StaticHeaders {
			if len(k) > maxHeaderKeyLen || len(v) > maxHeaderValueLen {
				http.Error(w, "Header key or value too long", http.StatusBadRequest)
				return
			}
		}

		staticHeaders, staticHeadersErr := json.Marshal(input.StaticHeaders)
		if staticHeadersErr != nil {
			logrus.WithError(staticHeadersErr).Error("Failed to marshal static headers")
			http.Error(w, "Invalid static headers", http.StatusBadRequest)
			return
		}

		if !validateEgressUrl(input.EgressUrl, svc.Config.Env) {
			logrus.WithField("egressUrl", input.EgressUrl).Error("Invalid egress URL")
			http.Error(w, "Invalid egress URL", http.StatusBadRequest)
			return
		}

		source, err := svc.CreateSource(r.Context(), db.CreateSourceParams{
			EgressUrl:     input.EgressUrl,
			StaticHeaders: staticHeaders,
			Description: pgtype.Text{
				String: input.Description,
				Valid:  input.Description != "",
			},
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			http.Error(w, "Failed to create source", http.StatusInternalServerError)
			return
		}
		logrus.WithField("source_id", source.ID).Info("Created new source")

		sourceDTO := mapperv1.ToSourceDTO(source, svc.Config)

		if err := api.JSON(w, http.StatusCreated, sourceDTO); err != nil {
			var writeErr *api.JSONWriteError
			if errors.As(err, &writeErr) {
				logrus.WithError(err).Error("Failed to write response")
			} else {
				logrus.WithError(err).Error("Failed to marshal response")
				http.Error(w, "Failed to create source", http.StatusInternalServerError)
			}
		}
	}
}

// NOTE: This SSRF protection is primarily for defense-in-depth. Since this is an internal tool
// with no user input, the risk is low. If the tool becomes public or user-facing, keep/enhance this check.
// validateEgressUrl checks if the provided egress URL is valid and does not point to local or private network addresses.
func validateEgressUrl(egressUrl, env string) bool {
	if egressUrl == "" {
		return false
	}
	parsedUrl, err := url.Parse(egressUrl)
	if err != nil {
		return false
	}
	if len(parsedUrl.String()) > 2048 ||
		(parsedUrl.Scheme != "http" &&
			parsedUrl.Scheme != "https") {
		return false
	}
	if env == "dev" {
		return true
	}
	host := parsedUrl.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}
	for _, ip := range ips {
		if ip.IsLoopback() ||
			ip.IsPrivate() ||
			ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast() {
			return false
		}
		// Block IPv4-mapped IPv6 loopback
		if ip.To4() == nil && ip.String() == "::1" {
			return false
		}
		// Block IPv4-mapped IPv6 for 127.0.0.0/8
		if ip.To4() == nil &&
			len(ip) == net.IPv6len &&
			ip[0] == 0 &&
			ip[1] == 0 &&
			ip[2] == 0 &&
			ip[3] == 0 &&
			ip[4] == 0 &&
			ip[5] == 0 &&
			ip[6] == 0 &&
			ip[7] == 0 &&
			ip[8] == 0 &&
			ip[9] == 0 &&
			ip[10] == 0xff &&
			ip[11] == 0xff &&
			ip[12] == 127 {
			return false
		}
	}
	return true
}

// sourcesRouter sets up the router for sources-related endpoints.
func sourcesRouter(svc *service.Service) chi.Router {
	r := chi.NewRouter()
	r.Get("/", listSources(svc))
	r.Get("/{source_id}", getSourceByID(svc))
	r.Post("/", createSource(svc))
	r.Mount("/{source_id}/events", eventsRouter(svc))
	return r
}
