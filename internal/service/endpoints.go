package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateEndpoint creates a new endpoint in the database.
func(service *Service) CreateEndpoint(name string, description *string) (db.Endpoint, error) {
	result, err := service.queries.RegisterEndpoint(context.Background(), db.RegisterEndpointParams{
		UserID: pgtype.Int8{Int64: 1, Valid: true}, // TODO : replace with actual user ID
		PublicKey: "test_public_key", // TODO : replace with actual public key generation
		Name: name,
		Description: pgtype.Text{
			String: *description,
			Valid: description != nil,
		},
	})
	return result, err
}
