// Package deliveryengine implements the background process responsible for processing and delivering webhooks.
package deliveryengine

import (
	"context"
	"net/http"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// Start initializes and starts the delivery engine, which periodically checks for pending deliveries and processes them.
func Start(svc *service.Service, pollInterval time.Duration) {
	logrus.Info("Starting delivery engine...")
	ctx := context.Background()
	ticker := time.NewTicker(pollInterval)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	defer ticker.Stop()
	semaphore := make(chan struct{}, 10) // TODO make max concurrency configurable
	for {
		pendingDeliveries, err := svc.ListPendingDeliveryAttempts(ctx)
		if err != nil {
			logrus.Error("Error listing pending deliveries:", err)
			continue
		}
		for _, delivery := range pendingDeliveries {
			semaphore <- struct{}{}
			go func(delivery db.ListPendingDeliveryAttemptsRow) {
				defer func() { <-semaphore }()
				attemptDelivery(svc, httpClient, delivery)
			}(delivery)
		}
		<-ticker.C
	}
}
