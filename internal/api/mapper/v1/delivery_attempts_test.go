package mapperv1

import (
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestToDeliveryAttemptDTO(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	finishedAt := startedAt.Add(2 * time.Minute)
	createdAt := startedAt.Add(-1 * time.Minute)
	nextAttemptAt := finishedAt.Add(5 * time.Minute)

	deliveryAttempt := db.DeliveryAttempt{
		ID:            11,
		EventID:       22,
		AttemptNumber: 3,
		State:         "failed",
		StatusCode:    pgtype.Int4{Int32: 500, Valid: true},
		ErrorType:     pgtype.Text{String: "network_error", Valid: true},
		ErrorMessage:  pgtype.Text{String: "upstream timeout", Valid: true},
		StartedAt:     pgtype.Timestamptz{Time: startedAt, Valid: true},
		FinishedAt:    pgtype.Timestamptz{Time: finishedAt, Valid: true},
		CreatedAt:     pgtype.Timestamptz{Time: createdAt, Valid: true},
		NextAttemptAt: pgtype.Timestamptz{Time: nextAttemptAt, Valid: true},
	}

	dto := ToDeliveryAttemptDTO(deliveryAttempt)

	require.NotNil(t, dto.StatusCode)
	require.NotNil(t, dto.ErrorType)
	require.NotNil(t, dto.ErrorMessage)
	require.NotNil(t, dto.StartedAt)
	require.NotNil(t, dto.FinishedAt)
	require.NotNil(t, dto.NextAttemptAt)

	assert.Equal(t, int64(11), dto.ID)
	assert.Equal(t, int64(22), dto.EventID)
	assert.Equal(t, int32(3), dto.AttemptNumber)
	assert.Equal(t, "failed", dto.State)
	assert.Equal(t, int32(500), *dto.StatusCode)
	assert.Equal(t, "network_error", *dto.ErrorType)
	assert.Equal(t, "upstream timeout", *dto.ErrorMessage)
	assert.Equal(t, startedAt, *dto.StartedAt)
	assert.Equal(t, finishedAt, *dto.FinishedAt)
	assert.Equal(t, createdAt, dto.CreatedAt)
	assert.Equal(t, nextAttemptAt, *dto.NextAttemptAt)
}

func TestToDeliveryAttemptDTO_InvalidOptionalFields(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	dto := ToDeliveryAttemptDTO(db.DeliveryAttempt{
		ID:        1,
		EventID:   2,
		CreatedAt: pgtype.Timestamptz{Time: createdAt, Valid: true},
	})

	assert.Nil(t, dto.StatusCode)
	assert.Nil(t, dto.ErrorType)
	assert.Nil(t, dto.ErrorMessage)
	assert.Nil(t, dto.StartedAt)
	assert.Nil(t, dto.FinishedAt)
	assert.Nil(t, dto.NextAttemptAt)
	assert.Equal(t, createdAt, dto.CreatedAt)
}

func TestToDeliveryAttemptDTOs(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	deliveryAttempts := []db.DeliveryAttempt{
		{ID: 1, EventID: 10, CreatedAt: pgtype.Timestamptz{Time: createdAt, Valid: true}},
		{ID: 2, EventID: 20, CreatedAt: pgtype.Timestamptz{Time: createdAt.Add(time.Minute), Valid: true}},
	}

	dtos := ToDeliveryAttemptDTOs(deliveryAttempts)

	require.Len(t, dtos, 2)
	assert.Equal(t, int64(1), dtos[0].ID)
	assert.Equal(t, int64(10), dtos[0].EventID)
	assert.Equal(t, int64(2), dtos[1].ID)
	assert.Equal(t, int64(20), dtos[1].EventID)
}
