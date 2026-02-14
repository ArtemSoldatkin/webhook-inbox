// Package deliveryengine implements the background process responsible for processing and delivering webhooks.
package deliveryengine

import (
	"context"
	"net/http"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/sirupsen/logrus"
)

// Start initializes and starts the delivery engine, which periodically checks for pending deliveries and processes them.
func Start(svc *service.Service, pollIntervalMs time.Duration) {
	logrus.Info("Starting delivery engine...")
	ctx := context.Background()
	ticker := time.NewTicker(pollIntervalMs * time.Millisecond)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	defer ticker.Stop()
	for {
		pendingDeliveries, err := svc.ListPendingDeliveryAttempts(ctx)
		if err != nil {
			logrus.Error("Error listing pending deliveries:", err)
			continue
		}
		for _, delivery := range pendingDeliveries {
			go AttemptDelivery(ctx, svc, httpClient, delivery)
		}
		<-ticker.C
	}
}
// TODO use semaphore to limit number of concurrent deliveries
// AttemptDelivery processes a single delivery attempt by retrieving the associated event and source, merging headers, and sending the HTTP request to the configured egress URL.
func AttemptDelivery(ctx context.Context, svc *service.Service, httpClient *http.Client, delivery db.ListPendingDeliveryAttemptsRow) {
	event, err := svc.GetEventByID(ctx, delivery.EventID)
	if err != nil {
		logrus.Error("Error retrieving event for delivery attempt:", err)
		return
	}
	source, err := svc.GetSourceByID(ctx, event.SourceID)
	if err != nil {
		logrus.Error("Error retrieving source for event:", err)
		return
	}
	staticHeaders, err := utils.JSONBtoType[map[string]string](source.StaticHeaders)
	if err != nil {
		logrus.Error("Error retrieving static headers for source:", err)
		return
	}
	rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
	if err != nil {
		logrus.Error("Error retrieving headers for event:", err)
		return
	}
	queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams)
	if err != nil {
		logrus.Error("Error retrieving query params for event:", err)
		return
	}
	body, err := utils.JSONBtoType[map[string]string](event.Body)
	if err != nil {
		logrus.Error("Error retrieving body for event:", err)
		return
	}
	headers := utils.MergeHeaders(staticHeaders, rawHeaders)
	req, err := http.NewRequest(event.Method, source.EgressUrl, nil)
	if err != nil {
		logrus.Error("Error creating HTTP request for delivery attempt:", err)
		return
	}
	req.Header = http.Header(headers)
	res, err := httpClient.Do(req)
	if err != nil {
		logrus.Error("Error sending HTTP request for delivery attempt:", err)
		return
	}
	defer res.Body.Close()
	logrus.Infof("Delivery attempt for event ID %d returned status code %d", event.ID, res.StatusCode)

	logrus.Infof("%s, %s, %s", headers, queryParams, body)
}