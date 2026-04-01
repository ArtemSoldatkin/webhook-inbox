// Package service provides business logic for the application.
package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/dgraph-io/ristretto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var beginTxFunc = func(dbPool *pgxpool.Pool, ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return dbPool.BeginTx(ctx, opts)
}

// Service encapsulates the business logic and database queries.
type Service struct {
	dbPool  *pgxpool.Pool
	queries *db.Queries
	Config  *config.Config
	Cache   *ristretto.Cache
}

// NewService creates a new instance of Service.
func NewService(
	dbPool *pgxpool.Pool,
	config *config.Config,
	cache *ristretto.Cache,
) *Service {
	queries := db.New(dbPool)

	return &Service{
		dbPool:  dbPool,
		queries: queries,
		Config:  config,
		Cache:   cache,
	}
}

// BeginTx starts a new database transaction and returns it.
// The caller is responsible for committing or rolling back the transaction and releasing the connection.
func (s *Service) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return beginTxFunc(s.dbPool, ctx, opts)
}
