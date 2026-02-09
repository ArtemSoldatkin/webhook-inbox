package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListDeliveryAttempts retrieves all delivery attempts for a given event ID from the database.
func (svc *Service) ListDeliveryAttempts(ctx context.Context, eventID int64) ([]db.DeliveryAttempt, error) {
	return svc.queries.ListDeliveryAttemptsByEvent(ctx, eventID)
}
