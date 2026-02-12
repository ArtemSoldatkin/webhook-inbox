package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(ctx context.Context, sourceID int64) ([]db.Event, error) {
	return svc.queries.ListEventsBySource(ctx, sourceID)
}

// CreateEvent inserts a new event into the database and returns its ID.
func (svc *Service) CreateEvent(ctx context.Context, event db.CreateEventParams) (int64, error) {
	return svc.queries.CreateEvent(ctx, event)
}
