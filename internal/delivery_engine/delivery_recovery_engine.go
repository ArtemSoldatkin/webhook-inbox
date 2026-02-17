package deliveryengine

import (
	"context"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// RecoverStuckDeliveryAttempts identifies and resets delivery attempts that have been in-flight for too long, allowing them to be retried by the delivery engine.
func recoverStuckDeliveryAttempts(svc *service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // TODO make timeout configurable
	defer cancel()

	return svc.RecoverStuckDeliveryAttempts(ctx)
}

// StartRecoveryEngine initializes and starts the delivery recovery engine, which periodically checks for stuck delivery attempts and resets them for retrying.
func StartRecoveryEngine(svc *service.Service, interval time.Duration) {
	logrus.Info("Starting delivery recovery engine...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := recoverStuckDeliveryAttempts(svc); err != nil {
			logrus.Error("Error recovering stuck delivery attempts:", err)
		}
		<-ticker.C
	}
}
