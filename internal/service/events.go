package service

import (
	"context"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
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
	return svc.queries.GetEventByID(ctx, eventID)
}

// CreateEvent inserts a new event into the database and returns its ID.
func (svc *Service) CreateEvent(ctx context.Context, event db.CreateEventParams) (int64, error) {
	return svc.queries.CreateEvent(ctx, event)
}
