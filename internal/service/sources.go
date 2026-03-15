package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
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
	sortDirection api.SortDirection,
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

// CreateSource creates a new source in the database with the provided parameters.
func (svc *Service) CreateSource(
	ctx context.Context,
	source db.CreateSourceParams,
) (db.Source, error) {
	return svc.queries.CreateSource(ctx, source)
}
