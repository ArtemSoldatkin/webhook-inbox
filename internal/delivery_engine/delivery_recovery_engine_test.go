package deliveryengine

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestRecoverStuckDeliveryAttempts(t *testing.T) {
	t.Parallel()

	dbtx := newEngineTestDB()
	svc := newEngineTestService(t, dbtx)
	svc.Config.APIRecoveryTimeoutSec = 5
	dbtx.execHandlers["-- name: RecoverStuckDeliveryAttempts :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		t.Fatalf("recover query should be called without arguments")
		return pgconn.NewCommandTag(""), nil
	}

	called := false
	dbtx.execHandlers["-- name: RecoverStuckDeliveryAttempts :exec"] = func(args ...any) (pgconn.CommandTag, error) {
		called = true
		assert.Len(t, args, 0)
		return pgconn.NewCommandTag("UPDATE 1"), nil
	}

	err := recoverStuckDeliveryAttempts(svc)

	require.NoError(t, err)
	assert.True(t, called)
}

func TestRecoverStuckDeliveryAttempts_UsesTimeoutContext(t *testing.T) {
	t.Parallel()

	dbtx := &recoveryContextDB{}
	svc := newEngineTestService(t, dbtx)
	svc.Config.APIRecoveryTimeoutSec = 3

	err := recoverStuckDeliveryAttempts(svc)

	require.NoError(t, err)
	require.NotNil(t, dbtx.deadline)
	assert.WithinDuration(t, time.Now().Add(3*time.Second), *dbtx.deadline, 2*time.Second)
}

type recoveryContextDB struct {
	deadline *time.Time
}

func (dbtx *recoveryContextDB) Query(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (dbtx *recoveryContextDB) QueryRow(_ context.Context, _ string, _ ...interface{}) pgx.Row {
	return nil
}

func (dbtx *recoveryContextDB) Exec(ctx context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	if deadline, ok := ctx.Deadline(); ok {
		dbtx.deadline = &deadline
	}
	return pgconn.NewCommandTag("UPDATE 1"), nil
}
