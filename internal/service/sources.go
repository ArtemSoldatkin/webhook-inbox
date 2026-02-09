package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// ListSources retrieves all sources from the database.
func (svc *Service) ListSources(ctx context.Context) ([]db.Source, error) {
	return svc.queries.ListSources(ctx)
}
