package routev1

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

var (
	httpRegexp        = regexp.MustCompile(`^https?://`)
	localhostRegexp   = regexp.MustCompile(`^https?://(localhost|127\.0\.0\.1|0\.0\.0\.0|\[?::1\]?)(/|:|$)`)
	private10Regexp   = regexp.MustCompile(`^https?://10\.`)
	private192Regexp  = regexp.MustCompile(`^https?://192\.168\.`)
	private172Regexp  = regexp.MustCompile(`^https?://172\.(1[6-9]|2[0-9]|3[0-1])\.`)
	metadata169Regexp = regexp.MustCompile(`^https?://169\.254\.169\.254(/|:|$)`)
)

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
		sources, err := svc.ListSources(r.Context(), cursor, pageSize)
		if err != nil {
			logrus.WithError(err).Error("Failed to list sources")
			http.Error(w, "Failed to list sources", http.StatusInternalServerError)
			return
		}
		sourceDTOs := make([]dtov1.SourceDTO, 0, len(sources))
		for _, source := range sources {
			staticHeaders, err := utils.JSONBtoType[map[string]string](source.StaticHeaders)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal static headers")
				staticHeaders = map[string]string{
					"__error": "Webhook Inbox Error - Failed to parse",
				}
			}
			var disableAt *time.Time
			if source.DisableAt.Valid {
				disableAt = &source.DisableAt.Time
			} else {
				disableAt = nil
			}
			sourceDTOs = append(sourceDTOs, dtov1.SourceDTO{
				ID:            source.ID,
				PublicID:      source.PublicID.String(),
				IngressUrl:    utils.GenerateIngressURL(svc.Config.APIProtocol, svc.Config.APIHost, svc.Config.APIPort, source.PublicID.String()),
				EgressUrl:     source.EgressUrl,
				StaticHeaders: staticHeaders,
				Status:        source.Status,
				StatusReason:  utils.PtrIfValid(source.StatusReason.String, source.StatusReason.Valid),
				Description:   utils.PtrIfValid(source.Description.String, source.Description.Valid),
				CreatedAt:     source.CreatedAt.Time,
				UpdatedAt:     source.UpdatedAt.Time,
				DisableAt:     disableAt,
			})
		}
		var nextCursor api.Cursor
		if len(sourceDTOs) > 0 {
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
		response, err := json.Marshal(paginatedResponse)
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

// getSourceByID handles GET requests to retrieve a source by its ID.
func getSourceByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceIDRaw := chi.URLParam(r, "sourceID")
		sourceID, err := strconv.ParseInt(sourceIDRaw, 10, 64)
		if err != nil {
			logrus.WithError(err).Error("Invalid source ID")
			http.Error(w, "Invalid source ID", http.StatusBadRequest)
			return
		}
		source, err := svc.GetSourceByID(r.Context(), sourceID)
		if err != nil {
			logrus.WithError(err).Error("Failed to get source by ID")
			http.Error(w, "Failed to get source", http.StatusInternalServerError)
			return
		}
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
		sourceDTO := dtov1.SourceDTO{
			ID:            source.ID,
			PublicID:      source.PublicID.String(),
			IngressUrl:    utils.GenerateIngressURL(svc.Config.APIProtocol, svc.Config.APIHost, svc.Config.APIPort, source.PublicID.String()),
			EgressUrl:     source.EgressUrl,
			StaticHeaders: staticHeaders,
			Status:        source.Status,
			StatusReason:  utils.PtrIfValid(source.StatusReason.String, source.StatusReason.Valid),
			Description:   utils.PtrIfValid(source.Description.String, source.Description.Valid),
			CreatedAt:     source.CreatedAt.Time,
			UpdatedAt:     source.UpdatedAt.Time,
			DisableAt:     disbaleAt,
		}
		response, err := json.Marshal(sourceDTO)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal source")
			http.Error(w, "Failed to get source", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// CreateSourceData defines the parameters required to create a new source.
type CreateSourceData struct {
	IngressUrl    string            `json:"IngressUrl"`
	EgressUrl     string            `json:"EgressUrl"`
	StaticHeaders map[string]string `json:"StaticHeaders,omitempty"`
	Description   string            `json:"Description,omitempty"`
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
		if !validateEgressUrl(data.EgressUrl, svc.Config.Env) {
			logrus.WithField("egressUrl", data.EgressUrl).Error("Invalid egress URL")
			http.Error(w, "Invalid egress URL", http.StatusBadRequest)
			return
		}
		source, err := svc.CreateSource(r.Context(), db.CreateSourceParams{
			EgressUrl:     data.EgressUrl,
			StaticHeaders: staticHeaders,
			Description:   pgtype.Text{String: data.Description, Valid: data.Description != ""},
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to create source")
			http.Error(w, "Failed to create source", http.StatusInternalServerError)
			return
		}
		sourceDTO := dtov1.SourceDTO{
			ID:            source.ID,
			PublicID:      source.PublicID.String(),
			IngressUrl:    utils.GenerateIngressURL(svc.Config.APIProtocol, svc.Config.APIHost, svc.Config.APIPort, source.PublicID.String()),
			EgressUrl:     source.EgressUrl,
			StaticHeaders: data.StaticHeaders,
			Status:        source.Status,
			StatusReason:  utils.PtrIfValid(source.StatusReason.String, source.StatusReason.Valid),
			Description:   utils.PtrIfValid(source.Description.String, source.Description.Valid),
			CreatedAt:     source.CreatedAt.Time,
			UpdatedAt:     source.UpdatedAt.Time,
		}
		response, err := json.Marshal(sourceDTO)
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

// NOTE: This SSRF protection is primarily for defense-in-depth. Since this is an internal tool
// with no user input, the risk is low. If the tool becomes public or user-facing, keep/enhance this check.
// validateEgressUrl checks if the provided egress URL is valid and does not point to local or private network addresses.
func validateEgressUrl(egressUrl, env string) bool {
	parsedUrl, err := url.Parse(egressUrl)
	if err != nil {
		return false
	}
	if len(parsedUrl.String()) > 2048 || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") {
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
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return false
		}
		// Block IPv4-mapped IPv6 loopback
		if ip.To4() == nil && ip.String() == "::1" {
			return false
		}
		// Block IPv4-mapped IPv6 for 127.0.0.0/8
		if ip.To4() == nil && len(ip) == net.IPv6len && ip[0] == 0 && ip[1] == 0 && ip[2] == 0 && ip[3] == 0 &&
			ip[4] == 0 && ip[5] == 0 && ip[6] == 0 && ip[7] == 0 && ip[8] == 0 && ip[9] == 0 && ip[10] == 0xff && ip[11] == 0xff &&
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
