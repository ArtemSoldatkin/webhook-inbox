package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListSources retrieves all sources from the database.
func (svc *Service) ListSources(ctx context.Context) ([]db.Source, error) {
	return svc.queries.ListSources(ctx)
}

// GetSourceByID retrieves a source by its ID from the database.
func (svc *Service) GetSourceByID(ctx context.Context, id int64) (db.Source, error) {
	return svc.queries.GetSourceByID(ctx, id)
}

// GetSourceByPublicID retrieves a source by its public ID from the database.
func (svc *Service) GetSourceByPublicID(ctx context.Context, publicID string) (db.Source, error) {
	var publicUUID pgtype.UUID
	err := publicUUID.Scan(publicID)
	if err != nil {
		return db.Source{}, err
	}
	return svc.queries.GetSourceByPublicID(ctx, publicUUID)
}

// CreateSource creates a new source in the database with the provided parameters.
func (svc *Service) CreateSource(ctx context.Context, source db.CreateSourceParams) (db.Source, error) {
	return svc.queries.CreateSource(ctx, source)
}