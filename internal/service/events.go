package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(ctx context.Context, sourceID int64) ([]db.Event, error) {
	return svc.queries.ListEventsBySource(ctx, sourceID)
}
