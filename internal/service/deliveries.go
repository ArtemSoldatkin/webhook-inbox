package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListDeliveries retrieves all deliveries associated with a specific endpoint.
func (svc *Service) ListDeliveries(ctx context.Context, endpointID int64) ([]db.Delivery, error) {
	return svc.queries.ListDeliveries(ctx, pgtype.Int8{Int64: endpointID, Valid: true})
}
