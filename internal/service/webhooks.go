package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListWebhooks retrieves all webhooks associated with a specific endpoint.
func (svc *Service) ListWebhooks(ctx context.Context, endpointID int64) ([]db.Webhook, error) {
	return svc.queries.ListWebhooks(ctx, pgtype.Int8{Int64: endpointID, Valid: true})
}

// CreateWebhook creates a new webhook for a given endpoint with the provided details.
func (svc *Service) CreateWebhook(ctx context.Context, endpointID int64, name, description string) (db.Webhook, error) {
	return svc.queries.CreateWebhook(ctx, db.CreateWebhookParams{
		EndpointID: pgtype.Int8{Int64: endpointID, Valid: true},
		PublicKey:  generatePublicKey(), // Assume this function generates a unique public key
		Name:        name,
		Description: pgtype.Text{String: description, Valid: true},
	})
}

// ToggleWebhook enables or disables a webhook based on its current state.
func (svc *Service) ToggleWebhook(ctx context.Context, webhookID int64) (db.Webhook, error) {
	return svc.queries.ToggleWebhook(ctx, webhookID)
}

// generatePublicKey is a placeholder function for generating a unique public key for a webhook.
func generatePublicKey() string {
	// TODO introduce actual key generation logic
	return "unique-public-key"
}