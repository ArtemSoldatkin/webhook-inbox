package service

import (
	"context"
	"encoding/json"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListEndpoints retrieves all endpoints associated with a specific user.
func (svc *Service) ListEndpoints(ctx context.Context, userID int64) ([]db.Endpoint, error) {
    return svc.queries.ListEndpoints(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

// RegisterEndpoint creates a new endpoint for a user with the provided details.
func (svc *Service) RegisterEndpoint(ctx context.Context, userID int64, url, name, description string, headers map[string]string) (db.Endpoint, error) {
    jsonHeaders, err := json.Marshal(headers)
    if err != nil {
        return db.Endpoint{}, err
    }
    return svc.queries.RegisterEndpoint(ctx, db.RegisterEndpointParams{
        UserID:      pgtype.Int8{Int64: userID, Valid: true},
        Url:         url,
        Name:        name,
        Description: pgtype.Text{String: description, Valid: true},
        Headers:     jsonHeaders,
    })
}

// ToggleEndpoint enables or disables an endpoint based on its current state.
func (svc *Service) ToggleEndpoint(ctx context.Context, endpointID int64) (db.Endpoint, error) {
    return svc.queries.ToggleEndpoint(ctx, endpointID)
}
