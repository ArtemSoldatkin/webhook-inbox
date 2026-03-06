package service

import (
	"context"
	"fmt"
	"time"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/sirupsen/logrus"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(
	ctx context.Context,
	sourceID int64,
	cursor api.Cursor,
	pageSize int,
) ([]db.Event, error) {
	cursorTS, cursorID := cursor.ToDBParams()
	return svc.queries.ListEventsBySource(ctx, db.ListEventsBySourceParams{
		SourceID: sourceID,
		CursorTs: cursorTS,
		CursorID: cursorID,
		PageSize: int32(pageSize),
	})
}

// GetEventByID retrieves a specific event by its ID from the database.
func (svc *Service) GetEventByID(ctx context.Context, eventID int64) (db.Event, error) {
	cacheKey := fmt.Sprintf("GetEventByID|%d", eventID)

	if cachedEvent, ok := svc.Cache.Get(cacheKey); ok {
		return cachedEvent.(db.Event), nil
	}

	event, err := svc.queries.GetEventByID(ctx, eventID)
	if err != nil {
		return db.Event{}, err
	}

	cacheCost, err := utils.EstimateStructSize(event)
	if err != nil {
		logrus.
			WithError(err).
			WithField("event_id", event.ID).
			Warning("failed to estimate cache cost for source, using default cost")
		cacheCost = svc.Config.APICacheDefaultCost
	}

	svc.Cache.SetWithTTL(
		cacheKey,
		event,
		cacheCost,
		time.Duration(svc.Config.APICacheEventTTLSec)*time.Second,
	)

	return event, nil
}

// CreateEvent inserts a new event into the database and returns its ID.
func (svc *Service) CreateEvent(
	ctx context.Context,
	event db.CreateEventParams,
) (int64, error) {
	return svc.queries.CreateEvent(ctx, event)
}
