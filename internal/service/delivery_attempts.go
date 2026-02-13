package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListDeliveryAttempts retrieves all delivery attempts for a given event ID from the database.
func (svc *Service) ListDeliveryAttempts(ctx context.Context, eventID int64) ([]db.DeliveryAttempt, error) {
	return svc.queries.ListDeliveryAttemptsByEvent(ctx, eventID)
}

// CreateDeliveryAttempt inserts a new delivery attempt into the database and returns its ID.
func (svc *Service) CreateDeliveryAttempt(ctx context.Context, delivery db.CreateDeliveryAttemptParams) (int64, error) {
	return svc.queries.CreateDeliveryAttempt(ctx, delivery)
}

// ListPendingDeliveryAttempts retrieves all delivery attempts that are currently pending and need to be processed by the delivery engine.
func (svc *Service) ListPendingDeliveryAttempts(ctx context.Context) ([]db.ListPendingDeliveryAttemptsRow, error) {
	return svc.queries.ListPendingDeliveryAttempts(ctx)
}
