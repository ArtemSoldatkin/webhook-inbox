package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestListDeliveryAttempts(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	ts := testNow()
	cursorID := int64(55)
	cursor := types.NewCursor(&ts, &cursorID)
	dbtx.queryHandlers["-- name: ListDeliveryAttemptsByEvent :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 7)
		assert.Equal(t, int64(8), args[0])
		assert.Equal(t, ts, args[1].(pgtype.Timestamptz).Time)
		assert.Equal(t, "ASC", args[2])
		assert.Equal(t, int64(55), args[3])
		assert.Equal(t, "timeout", args[4])
		assert.Equal(t, "failed", args[5])
		assert.Equal(t, int32(10), args[6])
		return &serviceTestRows{rows: [][]any{
			{int64(1), int64(8), int32(1), "failed", pgtype.Int4{}, pgtype.Text{}, pgtype.Text{}, pgtype.Timestamptz{}, pgtype.Timestamptz{}, pgtype.Timestamptz{}, pgtype.Timestamptz{}},
		}}, nil
	}

	attempts, err := svc.ListDeliveryAttempts(context.Background(), 8, cursor, 10, "timeout", "failed", "ASC")

	require.NoError(t, err)
	require.Len(t, attempts, 1)
	assert.Equal(t, int64(1), attempts[0].ID)
}

func TestCreateDeliveryAttempt(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	params := db.CreateDeliveryAttemptParams{
		EventID:       9,
		AttemptNumber: 2,
		State:         "pending",
	}
	dbtx.queryRowHandlers["-- name: CreateDeliveryAttempt :one"] = func(args ...any) pgx.Row {
		require.Len(t, args, 9)
		assert.Equal(t, int64(9), args[0])
		assert.Equal(t, int32(2), args[1])
		assert.Equal(t, "pending", args[2])
		return &serviceTestRow{values: []any{int64(101)}}
	}

	id, err := svc.CreateDeliveryAttempt(context.Background(), params)

	require.NoError(t, err)
	assert.Equal(t, int64(101), id)
}

func TestListPendingDeliveryAttempts(t *testing.T) {
	t.Parallel()

	svc := newServiceUnderTest(t, newServiceTestDB())
	tx := newServiceTestTx()
	tx.queryHandlers["-- name: SelectPendingDeliveryAttemptIDs :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 1)
		assert.Equal(t, int32(3), args[0])
		return &serviceTestRows{rows: [][]any{{int64(11)}, {int64(22)}}}, nil
	}
	tx.queryHandlers["-- name: UpdateDeliveryAttemptsToInFlight :many"] = func(args ...any) (pgx.Rows, error) {
		require.Len(t, args, 1)
		assert.Equal(t, []int64{11, 22}, args[0])
		return &serviceTestRows{rows: [][]any{
			{int64(11), int64(101), int32(1)},
			{int64(22), int64(202), int32(2)},
		}}, nil
	}

	originalBeginTx := beginTxFunc
	beginTxFunc = func(_ *pgxpool.Pool, ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
		return tx, nil
	}
	defer func() { beginTxFunc = originalBeginTx }()

	pending, err := svc.ListPendingDeliveryAttempts(context.Background(), 3)

	require.NoError(t, err)
	assert.Equal(t, 1, tx.commitCount)
	assert.Equal(t, 1, tx.rollbackCount)
	require.Len(t, pending, 2)
	assert.Equal(t, PendingDeliveryAttempt{ID: 11, EventID: 101, AttemptNumber: 1}, pending[0])
	assert.Equal(t, PendingDeliveryAttempt{ID: 22, EventID: 202, AttemptNumber: 2}, pending[1])
}

func TestListPendingDeliveryAttempts_BeginTxError(t *testing.T) {
	t.Parallel()

	svc := newServiceUnderTest(t, newServiceTestDB())
	expectedErr := errors.New("begin tx failed")
	originalBeginTx := beginTxFunc
	beginTxFunc = func(_ *pgxpool.Pool, ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
		return nil, expectedErr
	}
	defer func() { beginTxFunc = originalBeginTx }()

	pending, err := svc.ListPendingDeliveryAttempts(context.Background(), 3)

	assert.Nil(t, pending)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestUpdateDeliveryAttempt(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	params := db.UpdateDeliveryAttemptParams{
		State:             "failed",
		DeliveryAttemptID: 15,
	}
	dbtx.execHandlers["-- name: UpdateDeliveryAttempt :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		require.Len(t, args, 7)
		assert.Equal(t, "failed", args[0])
		assert.Equal(t, int64(15), args[6])
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	err := svc.UpdateDeliveryAttempt(context.Background(), params)

	require.NoError(t, err)
}

func TestRecoverStuckDeliveryAttempts(t *testing.T) {
	t.Parallel()

	dbtx := newServiceTestDB()
	svc := newServiceUnderTest(t, dbtx)
	dbtx.execHandlers["-- name: RecoverStuckDeliveryAttempts :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		assert.Len(t, args, 0)
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	err := svc.RecoverStuckDeliveryAttempts(context.Background())

	require.NoError(t, err)
}
