package routev1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	mapperv1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/mapper/v1"
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

var filterStatusOptions = map[string]bool{
	"active":      true,
	"paused":      true,
	"quarantined": true,
	"disabled":    true,
}

// listSources handles GET requests to list all sources.
func listSources(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageSize, cursor, err := api.ParsePaginationParams(
			r.URL.Query(),
			svc.Config.APIDefaultPageSize,
			svc.Config.APIMinPageSize,
			svc.Config.APIMaxPageSize,
		)
		if err != nil {
			logrus.
				WithError(err).
				Error("Invalid pagination parameters")
			http.Error(w, "Invalid pagination parameters", http.StatusBadRequest)
			return
		}

		searchQuery := r.URL.Query().Get("search")
		if len(searchQuery) > svc.Config.APIMaxSearchQueryLength {
			logrus.WithField("search_query_length", len(searchQuery)).Error("Search query is too long")
			http.Error(w, "Search query is too long", http.StatusBadRequest)
			return
		}

		filterStatus := api.ParseFilter(r.URL.Query(), "status", filterStatusOptions)
		sortDirection := api.ParseSortDirection(r.URL.Query(), api.SortDirectionDesc)

		logrus.WithFields(logrus.Fields{
			"pageSize":      pageSize,
			"cursor":        cursor,
			"search":        searchQuery,
			"filterStatus":  filterStatus,
			"sortDirection": sortDirection,
			"query":         r.URL.RawQuery,
		}).Debug("Received listSources request")

		sources, err := svc.ListSources(
			r.Context(),
			cursor,
			pageSize,
			searchQuery,
			filterStatus,
			sortDirection,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}

		sourceDTOs := mapperv1.ToSourceDTOs(sources, svc.Config)

		logrus.WithField("returned_count", len(sourceDTOs)).Debug("Returning sources")

		var nextCursor api.Cursor
		if len(sourceDTOs) > pageSize {
			lastSource := sourceDTOs[len(sourceDTOs)-1]
			nextCursor = api.NewCursor(
				&lastSource.UpdatedAt,
				&lastSource.ID,
			)
		}

		paginatedResponse := api.ToPaginatedResponse(
			sourceDTOs,
			pageSize,
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

// getSourceByID handles GET requests to retrieve a source by its ID.
func getSourceByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceIDRaw := chi.URLParam(r, "sourceID")

		logrus.WithFields(logrus.Fields{
			"source_id": sourceIDRaw,
			"query":     r.URL.RawQuery,
		}).Debug("Received getSourceByID request")

		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid source ID")
			http.Error(w, "Invalid source ID", http.StatusBadRequest)
			return
		}

		if sourceID <= 0 {
			logrus.WithField("source_id", sourceID).Error("Source ID must be a positive integer")
			http.Error(w, "Source ID must be a positive integer", http.StatusBadRequest)
			return
		}

		source, err := svc.GetSourceByID(r.Context(), sourceID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logrus.WithField("source_id", sourceID).Info("Source not found")
				http.Error(w, "Source not found", http.StatusNotFound)
				return
			}
			logrus.WithField("source_id", sourceID).WithError(err).Error("Failed to get source")
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

// CreateSourceData defines the parameters required to create a new source.
type CreateSourceData struct {
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

		var data CreateSourceData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logrus.WithError(err).Error("Failed to decode create source request")
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"egress_url":  data.EgressUrl,
			"description": data.Description,
		}).Debug("Create source request data")

		if len(data.Description) > maxDescriptionLen {
			logrus.Error("Description too long")
			http.Error(
				w,
				fmt.Sprintf("Description must be %d characters or less", maxDescriptionLen),
				http.StatusBadRequest,
			)
			return
		}

		if len(data.StaticHeaders) > maxHeaders {
			http.Error(w, "Too many headers", http.StatusBadRequest)
			return
		}
		for k, v := range data.StaticHeaders {
			if len(k) > maxHeaderKeyLen || len(v) > maxHeaderValueLen {
				http.Error(w, "Header key or value too long", http.StatusBadRequest)
				return
			}
		}

		staticHeaders, staticHeadersErr := json.Marshal(data.StaticHeaders)
		if staticHeadersErr != nil {
			logrus.WithError(staticHeadersErr).Error("Failed to marshal static headers")
			http.Error(w, "Invalid static headers", http.StatusBadRequest)
			return
		}

		if !validateEgressUrl(data.EgressUrl, svc.Config.Env) {
			logrus.WithField("egressUrl", data.EgressUrl).Error("Invalid egress URL")
			http.Error(w, "Invalid egress URL", http.StatusBadRequest)
			return
		}

		source, err := svc.CreateSource(r.Context(), db.CreateSourceParams{
			EgressUrl:     data.EgressUrl,
			StaticHeaders: staticHeaders,
			Description: pgtype.Text{
				String: data.Description,
				Valid:  data.Description != "",
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
	r.Get("/{sourceID}", getSourceByID(svc))
	r.Post("/", createSource(svc))
	r.Mount("/{sourceID}/events", eventsRouter(svc))
	return r
}
