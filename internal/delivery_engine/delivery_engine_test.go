package deliveryengine

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestProcessPendingDeliveries_ReturnsListError(t *testing.T) {
	t.Parallel()

	svc := newEngineTestService(t, newEngineTestDB())
	httpClient := &http.Client{Timeout: time.Second}
	semaphore := make(chan struct{}, 1)
	expectedErr := errors.New("list failure")

	err := processPendingDeliveries(
		context.Background(),
		svc,
		httpClient,
		semaphore,
		func(context.Context, int32) ([]service.PendingDeliveryAttempt, error) {
			return nil, expectedErr
		},
		func(*service.Service, *http.Client, service.PendingDeliveryAttempt) {
			t.Fatalf("attempt should not be called when listing fails")
		},
	)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Len(t, semaphore, 0)
}

func TestProcessPendingDeliveries_DispatchesAllAttempts(t *testing.T) {
	t.Parallel()

	svc := newEngineTestService(t, newEngineTestDB())
	svc.Config.APIDeliveryMaxConcurrency = 2
	httpClient := &http.Client{Timeout: time.Second}
	semaphore := make(chan struct{}, svc.Config.APIDeliveryMaxConcurrency)

	deliveries := []service.PendingDeliveryAttempt{
		{ID: 1, EventID: 101, AttemptNumber: 1},
		{ID: 2, EventID: 102, AttemptNumber: 2},
	}
	done := make(chan service.PendingDeliveryAttempt, len(deliveries))

	err := processPendingDeliveries(
		context.Background(),
		svc,
		httpClient,
		semaphore,
		func(_ context.Context, n int32) ([]service.PendingDeliveryAttempt, error) {
			assert.Equal(t, int32(2), n)
			return deliveries, nil
		},
		func(gotSvc *service.Service, gotClient *http.Client, delivery service.PendingDeliveryAttempt) {
			assert.Same(t, svc, gotSvc)
			assert.Same(t, httpClient, gotClient)
			done <- delivery
		},
	)

	require.NoError(t, err)

	var received []service.PendingDeliveryAttempt
	for range deliveries {
		select {
		case delivery := <-done:
			received = append(received, delivery)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for dispatched delivery")
		}
	}

	assert.ElementsMatch(t, deliveries, received)
	require.Eventually(t, func() bool { return len(semaphore) == 0 }, time.Second, 10*time.Millisecond)
}
