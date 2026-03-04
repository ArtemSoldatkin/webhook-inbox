// Package deliveryengine implements the background process responsible for processing and delivering webhooks.
package deliveryengine

import (
	"context"
	"net/http"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// Start initializes and starts the delivery engine, which periodically checks for pending deliveries and processes them.
func Start(svc *service.Service, pollInterval time.Duration) {
	ctx := context.Background()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	httpClient := &http.Client{
		Timeout: time.Duration(svc.Config.APIDeliveryRequestTimeoutSec) * time.Second,
	}

	semaphore := make(chan struct{}, svc.Config.APIDeliveryMaxConcurrency)
	for {
		pendingDeliveries, err := svc.ListPendingDeliveryAttempts(
			ctx,
			int32(svc.Config.APIDeliveryMaxConcurrency),
		)
		if err != nil {
			logrus.WithError(err).Error("Error listing pending deliveries")
			continue
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
				attemptDelivery(svc, httpClient, delivery)
			}(delivery)
		}
		<-ticker.C
	}
}
