package service

import (
    "context"
    "crypto/rand"
    "encoding/base64"

    "github.com/ArtemSoldatkin/webhook-inbox/internal/db"
    "github.com/jackc/pgx/v5/pgtype"
)

// CreateEndpoint creates a new endpoint in the database.
func (service *Service) CreateEndpoint(name string, description *string) (db.Endpoint, error) {
    var desc pgtype.Text
    if description != nil {
        desc = pgtype.Text{String: *description, Valid: true}
    }

    result, err := service.queries.RegisterEndpoint(context.Background(), db.RegisterEndpointParams{
        UserID:      pgtype.Int8{Int64: 1, Valid: true}, // TODO : replace with actual user ID
        Url:         generatePublicKey(),
        Name:        name,
        Description: desc,
        Headers:     []byte(`{}`),
    })
    return result, err
}

// generatePublicKey returns a URL-safe base64-encoded 32-byte random string.
func generatePublicKey() string {
    const n = 32
    var b [n]byte
    if _, err := rand.Read(b[:]); err != nil {
        // best-effort fallback; still return deterministic-length key
    }
    return base64.RawURLEncoding.EncodeToString(b[:])
}
