package service

import (
	"context"

	api "github.com/ArtemSoldatkin/webhook-inbox/internal/api/utils"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5"
)

// This file contains service methods related to managing delivery attempts,
// including listing attempts for an event, creating new attempts,
// updating attempt status, and recovering stuck attempts.
type PendingDeliveryAttempt struct {
	ID            int64
	EventID       int64
	AttemptNumber int32
}

// ListDeliveryAttempts retrieves all delivery attempts for a given event ID from the database.
func (svc *Service) ListDeliveryAttempts(
	ctx context.Context,
	eventID int64,
	cursor api.Cursor,
	pageSize int,
	searchQuery string,
	filterState string,
) ([]db.DeliveryAttempt, error) {
	cursorTS, cursorID := cursor.ToDBParams()
	return svc.queries.ListDeliveryAttemptsByEvent(
		ctx,
		db.ListDeliveryAttemptsByEventParams{
			EventID:     eventID,
			CursorTs:    cursorTS,
			CursorID:    cursorID,
			SearchQuery: searchQuery,
			PageSize:    int32(pageSize),
			FilterState: filterState,
		},
	)
}

// CreateDeliveryAttempt inserts a new delivery attempt into the database and returns its ID.
func (svc *Service) CreateDeliveryAttempt(
	ctx context.Context,
	delivery db.CreateDeliveryAttemptParams,
) (int64, error) {
	return svc.queries.CreateDeliveryAttempt(ctx, delivery)
}

// ListPendingDeliveryAttempts retrieves a list of pending delivery attempts that are ready to be processed by the delivery engine,
// marking them as in-flight to prevent multiple workers from processing the same attempt concurrently.
func (svc *Service) ListPendingDeliveryAttempts(
	ctx context.Context,
	nDeliveries int32,
) ([]PendingDeliveryAttempt, error) {
	tx, err := svc.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := db.New(tx)

	deliveryAttemptIds, err := queries.SelectPendingDeliveryAttemptIDs(
		ctx,
		nDeliveries,
	)
	if err != nil {
		return nil, err
	}

	rows, err := queries.UpdateDeliveryAttemptsToInFlight(
		ctx,
		deliveryAttemptIds,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	pendingDeliveryAttempts := make([]PendingDeliveryAttempt, len(rows))
	for i, row := range rows {
		pendingDeliveryAttempts[i] = PendingDeliveryAttempt{
			ID:            row.ID,
			EventID:       row.EventID,
			AttemptNumber: row.AttemptNumber,
		}
	}
	return pendingDeliveryAttempts, nil
}

// UpdateDeliveryAttempt updates the status
// and other relevant fields of a delivery attempt in the database
// after it has been processed by the delivery engine.
func (svc *Service) UpdateDeliveryAttempt(
	ctx context.Context,
	attempt db.UpdateDeliveryAttemptParams,
) error {
	return svc.queries.UpdateDeliveryAttempt(ctx, attempt)
}

// RecoverStuckDeliveryAttempts identifies
// and resets delivery attempts that have been in-flight for too long,
// allowing them to be retried by the delivery engine.
func (svc *Service) RecoverStuckDeliveryAttempts(ctx context.Context) error {
	return svc.queries.RecoverStuckDeliveryAttempts(ctx)
}
