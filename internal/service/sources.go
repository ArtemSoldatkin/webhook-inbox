package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListSources retrieves all sources from the database.
func (svc *Service) ListSources(ctx context.Context) ([]db.Source, error) {
	return svc.queries.ListSources(ctx)
}

// CreateSource creates a new source in the database with the provided parameters.
func (svc *Service) CreateSource(ctx context.Context, source db.CreateSourceParams) (db.Source, error) {
	return svc.queries.CreateSource(ctx, source)
}