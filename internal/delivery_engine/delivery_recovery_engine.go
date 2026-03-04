package deliveryengine

import (
	"context"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/sirupsen/logrus"
)

// RecoverStuckDeliveryAttempts identifies and resets delivery attempts that have been in-flight for too long, allowing them to be retried by the delivery engine.
func recoverStuckDeliveryAttempts(svc *service.Service) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(svc.Config.APIRecoveryTimeoutSec)*time.Second,
	)
	defer cancel()

	return svc.RecoverStuckDeliveryAttempts(ctx)
}

// StartRecoveryEngine initializes and starts the delivery recovery engine, which periodically checks for stuck delivery attempts and resets them for retrying.
func StartRecoveryEngine(svc *service.Service, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		logrus.Debug("Checking for stuck delivery attempts to recover...")
		if err := recoverStuckDeliveryAttempts(svc); err != nil {
			logrus.WithError(err).Error("Error recovering stuck delivery attempts")
		}
		<-ticker.C
	}
}
