package service

import (
	"context"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(ctx context.Context, sourceID int64, cursor *time.Time, pageSize int) ([]db.Event, error) {
	var cursorValue pgtype.Timestamp
	if cursor != nil {
		cursorValue = pgtype.Timestamp{Time: *cursor, Valid: true}
	} else {
		cursorValue = pgtype.Timestamp{Valid: false}
	}
	return svc.queries.ListEventsBySource(ctx, db.ListEventsBySourceParams{
		SourceID: sourceID,
		Cursor:   cursorValue,
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
