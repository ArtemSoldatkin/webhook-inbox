// Package deliveryengine implements the background process responsible for processing and delivering webhooks.
package deliveryengine

import (
	"context"
	"net/http"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
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
	inFlightErr := svc.UpdateDeliveryAttempt(ctx, db.UpdateDeliveryAttemptParams{
		ID: delivery.ID,
		State: "in_flight",
		StartedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if inFlightErr != nil {
		logrus.Error("Error updating delivery attempt state to in_flight:", inFlightErr)
		return
	}
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
	var deliveryState string
	var errorType string
	if res.StatusCode >= 200 && res.StatusCode < 400 {
		deliveryState = "succeeded"
	} else {
		deliveryState = "failed"
		if res.StatusCode >= 500 {
			errorType = "http_5xx"
		} else {
			errorType = "http_4xx"
		}
	}
	var errorMessage string
	if deliveryState == "failed" {
		errorMessage = http.StatusText(res.StatusCode)
	}
	finishErr := svc.UpdateDeliveryAttempt(ctx, db.UpdateDeliveryAttemptParams{
		ID: delivery.ID,
		State: deliveryState,
		StatusCode: pgtype.Int4{Int32: int32(res.StatusCode), Valid: true},
		ErrorType: pgtype.Text{String: errorType, Valid: errorType != ""},
		ErrorMessage: pgtype.Text{String: errorMessage, Valid: errorMessage != ""},
		FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if finishErr != nil {
		logrus.Error("Error updating delivery attempt with result:", finishErr)
		return
	}
	logrus.Infof("Delivery attempt for event ID %d returned status code %d", event.ID, res.StatusCode)
	logrus.Infof("%s, %s, %s", headers, queryParams, body)
}