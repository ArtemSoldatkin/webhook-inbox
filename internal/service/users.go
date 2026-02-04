package service

import (
	"context"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
)

// CreateUser creates a new user with the given email.
func (svc *Service) CreateUser(ctx context.Context, email string) (db.User, error) {
	return svc.queries.CreateUser(ctx, email)
}
