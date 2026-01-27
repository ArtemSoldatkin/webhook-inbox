// Package service provides business logic for the application.
package service

import "github.com/ArtemSoldatkin/webhook-inbox/internal/db"

// Service encapsulates the business logic and database queries.
type Service struct {
	queries *db.Queries
}

// NewService creates a new instance of Service.
func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}