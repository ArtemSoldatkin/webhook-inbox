// Package deliveryengine implements the background process responsible for processing
// and delivering webhooks.
package deliveryengine

import (
	"context"
	"net/http"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

type listPendingDeliveriesFunc func(context.Context, int32) ([]service.PendingDeliveryAttempt, error)
type attemptDeliveryFunc func(*service.Service, *http.Client, service.PendingDeliveryAttempt)

// processPendingDeliveries retrieves pending deliveries and dispatches them for processing,
// respecting the concurrency limit defined by the semaphore.
func processPendingDeliveries(
	ctx context.Context,
	svc *service.Service,
	httpClient *http.Client,
	semaphore chan struct{},
	listPending listPendingDeliveriesFunc,
	attempt attemptDeliveryFunc,
) error {
	pendingDeliveries, err := listPending(
		ctx,
		int32(svc.Config.APIDeliveryMaxConcurrency),
	)
	if err != nil {
		logrus.WithError(err).Error("Error listing pending deliveries")
		return err
	}

	logrus.WithField("count", len(pendingDeliveries)).Debug("Pending deliveries found")
	for _, delivery := range pendingDeliveries {
		logrus.WithFields(logrus.Fields{
			"delivery_id": delivery.ID,
			"event_id":    delivery.EventID,
			"attempt":     delivery.AttemptNumber,
		}).Debug("Starting delivery attempt")

		semaphore <- struct{}{}

		go func(delivery service.PendingDeliveryAttempt) {
			defer func() { <-semaphore }()

			attempt(svc, httpClient, delivery)
		}(delivery)
	}

	return nil
}

// Start initializes and starts the delivery engine,
// which periodically checks for pending deliveries and processes them.
func Start(svc *service.Service, pollInterval time.Duration) {
	ctx := context.Background()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	httpClient := &http.Client{
		Timeout: time.Duration(svc.Config.APIDeliveryRequestTimeoutSec) * time.Second,
	}

	semaphore := make(chan struct{}, svc.Config.APIDeliveryMaxConcurrency)
	for {
		_ = processPendingDeliveries(
			ctx,
			svc,
			httpClient,
			semaphore,
			svc.ListPendingDeliveryAttempts,
			attemptDelivery,
		)
		<-ticker.C
	}
}
