package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(ctx context.Context, sourceID int64) ([]db.Event, error) {
	return svc.queries.ListEventsBySource(ctx, sourceID)
}

// GetEvent retrieves a specific event by its ID from the database.
func (svc *Service) GetEvent(ctx context.Context, eventID int64) (db.Event, error) {
	return svc.queries.GetEventByID(ctx, eventID)
}

// CreateEvent inserts a new event into the database and returns its ID.
func (svc *Service) CreateEvent(ctx context.Context, event db.CreateEventParams) (int64, error) {
	return svc.queries.CreateEvent(ctx, event)
}
