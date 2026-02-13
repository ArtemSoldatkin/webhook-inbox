// Package deliveryengine implements the background process responsible for processing and delivering webhooks.
package deliveryengine

import (
	"context"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// Start initializes and starts the delivery engine, which periodically checks for pending deliveries and processes them.
func Start(svc *service.Service, pollIntervalMs time.Duration) {
	logrus.Info("Starting delivery engine...")
	ctx := context.Background()
	ticker := time.NewTicker(pollIntervalMs * time.Millisecond)
	defer ticker.Stop()
	for {
		pendingDeliveries, err := svc.ListPendingDeliveryAttempts(ctx)
		if err != nil {
			logrus.Error("Error listing pending deliveries:", err)
			continue
		}
		for _, delivery := range pendingDeliveries {
			event, err := svc.GetEvent(ctx, delivery.EventID)
			if err != nil {
				logrus.Error("Error retrieving event for delivery attempt:", err)
				continue
			}
			logrus.Infof("Processing delivery attempt %d for event %d", delivery.ID, event.ID)
		}
		logrus.Infof("Found %s pending deliveries", pendingDeliveries)
		<-ticker.C
	}
}