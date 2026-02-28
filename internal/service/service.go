// Package service provides business logic for the application.
package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service encapsulates the business logic and database queries.
type Service struct {
	dbPool  *pgxpool.Pool
	queries *db.Queries
	Config  *config.Config
}

// NewService creates a new instance of Service.
func NewService(dbPool *pgxpool.Pool, config *config.Config) *Service {
	queries := db.New(dbPool)

	return &Service{
		queries: queries,
		Config:  config,
	}
}

// BeginTx starts a new database transaction and returns it.
// The caller is responsible for committing or rolling back the transaction and releasing the connection.
func (s *Service) BeginTx(ctx context.Context) (pgx.Tx, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}
	return tx, nil
}
