package mapperv1

import (
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestToSourceDTO(t *testing.T) {
	t.Parallel()

	publicID := mustUUID(t, "c5c2cb2b-a66d-4f1d-a661-8e5e1ed82171")
	createdAt := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	updatedAt := createdAt.Add(10 * time.Minute)
	disableAt := updatedAt.Add(time.Hour)

	cfg := &config.Config{
		APIProtocol: "https",
		APIHost:     "api.example.com",
		APIPort:     8443,
	}

	source := db.Source{
		ID:            303,
		PublicID:      publicID,
		EgressUrl:     "https://consumer.example.com/webhook",
		StaticHeaders: []byte(`{"Authorization":"Bearer token"}`),
		Status:        "active",
		StatusReason:  pgtype.Text{String: "healthy", Valid: true},
		Description:   pgtype.Text{String: "billing events", Valid: true},
		CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
		UpdatedAt:     pgtype.Timestamptz{Time: updatedAt, Valid: true},
		DisableAt:     pgtype.Timestamptz{Time: disableAt, Valid: true},
	}

	dto := ToSourceDTO(source, cfg)

	require.NotNil(t, dto.StatusReason)
	require.NotNil(t, dto.Description)
	require.NotNil(t, dto.DisableAt)

	assert.Equal(t, int64(303), dto.ID)
	assert.Equal(t, publicID.String(), dto.PublicID)
	assert.Equal(t, "https://api.example.com:8443/ingest/"+publicID.String(), dto.IngressUrl)
	assert.Equal(t, "https://consumer.example.com/webhook", dto.EgressUrl)
	assert.Equal(t, map[string]string{"Authorization": "Bearer token"}, dto.StaticHeaders)
	assert.Equal(t, "active", dto.Status)
	assert.Equal(t, "healthy", *dto.StatusReason)
	assert.Equal(t, "billing events", *dto.Description)
	assert.Equal(t, createdAt, dto.CreatedAt)
	assert.Equal(t, updatedAt, dto.UpdatedAt)
	assert.Equal(t, disableAt, *dto.DisableAt)
}

func TestToSourceDTO_InvalidStaticHeadersFallback(t *testing.T) {
	t.Parallel()

	dto := ToSourceDTO(db.Source{
		StaticHeaders: []byte(`{`),
	}, &config.Config{
		APIProtocol: "http",
		APIHost:     "localhost",
		APIPort:     8080,
	})

	assert.Equal(t, map[string]string{
		"__error": "Webhook Inbox Error - Failed to parse",
	}, dto.StaticHeaders)
	assert.Nil(t, dto.StatusReason)
	assert.Nil(t, dto.Description)
	assert.Nil(t, dto.DisableAt)
}

func TestToSourceDTOs(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		APIProtocol: "http",
		APIHost:     "localhost",
		APIPort:     8080,
	}
	firstID := mustUUID(t, "11111111-1111-1111-1111-111111111111")
	secondID := mustUUID(t, "22222222-2222-2222-2222-222222222222")

	sources := []db.Source{
		{ID: 1, PublicID: firstID, StaticHeaders: []byte(`{"A":"1"}`)},
		{ID: 2, PublicID: secondID, StaticHeaders: []byte(`{"B":"2"}`)},
	}

	dtos := ToSourceDTOs(sources, cfg)

	require.Len(t, dtos, 2)
	assert.Equal(t, int64(1), dtos[0].ID)
	assert.Equal(t, firstID.String(), dtos[0].PublicID)
	assert.Equal(t, int64(2), dtos[1].ID)
	assert.Equal(t, secondID.String(), dtos[1].PublicID)
}

func mustUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()

	var id pgtype.UUID
	require.NoError(t, id.Scan(value))
	return id
}
