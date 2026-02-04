package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListEvents retrieves all events associated with a specific webhook.
func (svc *Service) ListEvents(ctx context.Context, webhookID int64) ([]db.Event, error) {
	return svc.queries.ListEvents(ctx, pgtype.Int8{Int64: webhookID, Valid: true})
}
