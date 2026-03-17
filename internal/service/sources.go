package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
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

// ListSources retrieves all sources from the database.
func (svc *Service) ListSources(
	ctx context.Context,
	cursor types.Cursor,
	pageSize int,
	searchQuery string,
	filterStatus string,
	sortDirection string,
) ([]db.Source, error) {
	cursorTS, cursorID := cursor.ToDBParams()
	return svc.queries.ListSources(ctx, db.ListSourcesParams{
		CursorTs:      cursorTS,
		CursorID:      cursorID,
		SearchQuery:   searchQuery,
		PageSize:      int32(pageSize),
		FilterStatus:  filterStatus,
		SortDirection: string(sortDirection),
	})
}

// GetSourceByID retrieves a source by its ID from the database.
func (svc *Service) GetSourceByID(ctx context.Context, id int64) (db.Source, error) {
	return svc.queries.GetSourceByID(ctx, id)
}

// GetSourceByPublicID retrieves a source by its public ID from the database.
func (svc *Service) GetSourceByPublicID(
	ctx context.Context,
	publicID string,
) (db.Source, error) {
	cacheKey := fmt.Sprintf("GetSourceByPublicID|%s", publicID)

	if cachedSource, ok := svc.Cache.Get(cacheKey); ok {
		source, ok := cachedSource.(db.Source)
		if ok {
			return source, nil
		}

		logrus.
			WithField("public_id", publicID).
			Warning("cache hit for GetSourceByPublicID but value has unexpected type, ignoring cache")
	}

	var publicUUID pgtype.UUID
	if err := publicUUID.Scan(publicID); err != nil {
		return db.Source{}, err
	}

	source, err := svc.queries.GetSourceByPublicID(ctx, publicUUID)
	if err != nil {
		return db.Source{}, err
	}

	cacheCost, err := utils.EstimateStructSize(source)
	if err != nil {
		logrus.
			WithError(err).
			WithField("source_id", source.ID).
			Warning("failed to estimate cache cost for source, using default cost")
		cacheCost = svc.Config.APICacheDefaultCost
	}

	svc.Cache.SetWithTTL(
		cacheKey,
		source,
		cacheCost,
		time.Duration(svc.Config.APICacheSourceTTLSec)*time.Second,
	)

	return source, nil
}

// ListSources retrieves a paginated list of sources based on the provided parameters.
type CreateSourceInput struct {
	EgressUrl     string
	StaticHeaders map[string]string
	Description   string
}

// CreateSource creates a new source in the database with the provided parameters.
func (svc *Service) CreateSource(
	ctx context.Context,
	source CreateSourceInput,
) (db.Source, error) {
	staticHeaders, staticHeadersErr := json.Marshal(source.StaticHeaders)
	if staticHeadersErr != nil {
		return db.Source{}, fmt.Errorf("failed to marshal static headers: %w", staticHeadersErr)
	}

	if !validateEgressUrl(source.EgressUrl, svc.Config.Env) {
		return db.Source{}, fmt.Errorf("invalid egress URL: %s", source.EgressUrl)
	}

	return svc.queries.CreateSource(ctx, db.CreateSourceParams{
		EgressUrl:     source.EgressUrl,
		StaticHeaders: staticHeaders,
		Description: pgtype.Text{
			String: source.Description,
			Valid:  source.Description != "",
		},
	})
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
	// Use a context with timeout to avoid hanging on DNS resolution.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ipAddrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return false
	}
	for _, addr := range ipAddrs {
		ip := addr.IP
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
