package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
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
func (svc *Service) GetSourceByID(ctx context.Context, sourceID int64) (db.Source, error) {
	cacheKey := fmt.Sprintf("GetSourceByID|%s", strconv.FormatInt(sourceID, 10))

	if cachedSource, ok := svc.Cache.Get(cacheKey); ok {
		source, ok := cachedSource.(db.Source)
		if ok {
			return source, nil
		}

		logrus.
			WithField("source_id", sourceID).
			Warning("cache hit for GetSourceByID but value has unexpected type, ignoring cache")
	}

	source, err := svc.queries.GetSourceByID(ctx, sourceID)
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

// CreateSourceInput contains the parameters required to create a new source.
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

// allowedStatusTransitions defines the valid status transitions for sources.
var allowedStatusTransitions = map[string][]string{
	"active":      {"paused", "quarantined"},
	"paused":      {"active", "quarantined", "disabled"},
	"quarantined": {"active", "disabled"},
	"disabled":    {"active"},
}

// isValidStatusTransition checks if transitioning from currentStatus to newStatus is allowed based on predefined rules.
func isValidStatusTransition(currentStatus, newStatus string) bool {
	allowedTransitions, ok := allowedStatusTransitions[currentStatus]

	if !ok {
		return false
	}

	return slices.Contains(allowedTransitions, newStatus)
}

// UpdateSourceStatusInput contains the parameters required to update the status of a source.
type UpdateSourceStatusInput struct {
	SourceID     int64
	Status       string
	StatusReason string
}

// UpdateSourceStatus updates the status of a source in the database.
func (svc *Service) UpdateSourceStatus(ctx context.Context, sourceStatusInput UpdateSourceStatusInput) (db.Source, error) {
	source, err := svc.GetSourceByID(ctx, sourceStatusInput.SourceID)
	if err != nil {
		return db.Source{}, fmt.Errorf("failed to retrieve source with ID %d: %w", sourceStatusInput.SourceID, err)
	}

	if source.Status == sourceStatusInput.Status {
		return source, nil
	}

	if !isValidStatusTransition(source.Status, sourceStatusInput.Status) {
		return source, fmt.Errorf("invalid status transition from %s to %s", source.Status, sourceStatusInput.Status)
	}

	if err := svc.queries.UpdateSourceStatus(ctx, db.UpdateSourceStatusParams{
		SourceID: sourceStatusInput.SourceID,
		Status:   sourceStatusInput.Status,
		StatusReason: pgtype.Text{
			String: sourceStatusInput.StatusReason,
			Valid:  sourceStatusInput.StatusReason != "",
		},
	}); err != nil {
		return source, fmt.Errorf("failed to update status for source with ID %d: %w", sourceStatusInput.SourceID, err)
	}

	// Invalidate cache for this source since its status has changed
	svc.Cache.Del(fmt.Sprintf("GetSourceByID|%s", strconv.FormatInt(sourceStatusInput.SourceID, 10)))
	svc.Cache.Del(fmt.Sprintf("GetSourceByPublicID|%s", source.PublicID.String()))

	return svc.GetSourceByID(ctx, sourceStatusInput.SourceID)
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
