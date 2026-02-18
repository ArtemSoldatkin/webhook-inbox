package deliveryengine

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

// markInFlight updates the state of a delivery attempt to "in_flight" and sets the started_at timestamp to the current time.
func markInFlight(ctx context.Context, svc *service.Service, deliveryID int64) error {
	return svc.UpdateDeliveryAttempt(ctx, db.UpdateDeliveryAttemptParams{
		ID: deliveryID,
		State: "in_flight",
		StartedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}

// markInPending updates the state of a delivery attempt to "pending" and clears the started_at timestamp.
func markInPending(ctx context.Context, svc *service.Service, deliveryID int64) error {
	return svc.UpdateDeliveryAttempt(ctx, db.UpdateDeliveryAttemptParams{
		ID: deliveryID,
		State: "pending",
	})
}

// finalizeDeliveryAttempt updates the delivery attempt with the final result of the delivery, including the status code, error type, and error message if applicable, and sets the finished_at timestamp to the current time.
func finalizeDeliveryAttempt(ctx context.Context, svc *service.Service, deliveryID int64, result *DeliveryResult) error {
	return svc.UpdateDeliveryAttempt(ctx, db.UpdateDeliveryAttemptParams{
		ID: deliveryID,
		State: result.DeliveryState,
		StatusCode: pgtype.Int4{Int32: int32(result.StatusCode), Valid: true},
		ErrorType: pgtype.Text{String: result.ErrorType, Valid: result.ErrorType != ""},
		ErrorMessage: pgtype.Text{String: result.ErrorMessage, Valid: result.ErrorMessage != ""},
		FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}


// scheduleRetry creates a new delivery attempt with an incremented attempt number and a state of "pending" to schedule a retry for a failed delivery attempt.
func scheduleRetry(ctx context.Context, svc *service.Service, delivery db.ListPendingDeliveryAttemptsRow) (int64, error) {
	// TODO implement retry scheduling logic, e.g. using exponential backoff
	return svc.CreateDeliveryAttempt(ctx, db.CreateDeliveryAttemptParams{
				EventID: delivery.EventID,
				AttemptNumber: delivery.AttemptNumber + 1,
				State: "pending",
			})
}

// DeliveryPayload represents the data needed to send a delivery request, including the target URL, HTTP method, headers, query parameters, and body.
type DeliveryPayload struct {
	URL string
	Method string
	Headers map[string][]string
	QueryParams map[string][]string
	Body []byte
}

// loadDeliveryPayload retrieves the event and source associated with a delivery attempt, merges static and dynamic headers, and constructs the payload for the delivery request.
func loadDeliveryPayload(ctx context.Context, svc *service.Service, delivery db.ListPendingDeliveryAttemptsRow) (*DeliveryPayload, error) {
	event, err := svc.GetEventByID(ctx, delivery.EventID)
	if err != nil {
		return nil, err
	}
	source, err := svc.GetSourceByID(ctx, event.SourceID)
	if err != nil {
		return nil, err
	}
	staticHeaders, err := utils.JSONBtoType[map[string]string](source.StaticHeaders)
	if err != nil {
		return 	nil, err
	}
	rawHeaders, err := utils.JSONBtoType[map[string][]string](event.RawHeaders)
	if err != nil {
		return nil, err
	}
	headers := utils.MergeHeaders(staticHeaders, rawHeaders)
	queryParams, err := utils.JSONBtoType[map[string][]string](event.QueryParams)
	if err != nil {
		return nil, err
	}
	return &DeliveryPayload{
		URL: source.EgressUrl,
		Method: event.Method,
		Headers: headers,
		QueryParams: queryParams,
		Body: event.Body,
	}, nil
}

// sendDeliveryRequest constructs and sends an HTTP request based on the provided delivery payload and returns the response or an error if the request fails.
func sendDeliveryRequest(ctx context.Context, httpClient *http.Client, payload *DeliveryPayload) (*http.Response, error) {
	URL, err := url.Parse(payload.URL)
    if err != nil {
        return nil, err
    }
	query := URL.Query()
	for key, values := range payload.QueryParams {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	URL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, payload.Method, URL.String(), bytes.NewReader(payload.Body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header(payload.Headers)
	res, err := httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Error sending HTTP request for delivery attempt")
		return nil, err
	}
	return res, nil
}

// DeliveryResult represents the outcome of a delivery attempt, including the HTTP status code, final delivery state, and any error information if the delivery failed.
type DeliveryResult struct {
	StatusCode int
	DeliveryState string
	ErrorType string
	ErrorMessage string
}

// interpretDeliveryResponse analyzes the HTTP response from a delivery attempt and determines the final delivery state, error type, and error message if applicable, based on the status code of the response.
func interpretDeliveryResponse(res *http.Response) *DeliveryResult {
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
	return &DeliveryResult{
		StatusCode: res.StatusCode,
		DeliveryState: deliveryState,
		ErrorType: errorType,
		ErrorMessage: errorMessage,
	}
}

// handleDeliveryFinalizationAndRetry finalizes the delivery attempt by updating its state and result in the database, and if the delivery failed, it schedules a retry if the maximum number of attempts has not been reached.
func handleDeliveryFinalizationAndRetry(ctx context.Context, svc *service.Service, delivery db.ListPendingDeliveryAttemptsRow, result *DeliveryResult) {
	if err := finalizeDeliveryAttempt(ctx, svc, delivery.ID, result); err != nil {
		logrus.WithError(err).Error("Error finalizing delivery attempt")
		return
	}

	logrus.Infof("Finalized delivery attempt for event ID %d with status code %d and delivery state %s", delivery.EventID, result.StatusCode, result.DeliveryState)

	if result.DeliveryState == "failed" {
		logrus.Warnf("Delivery attempt for event ID %d failed with status code %d: %s", delivery.EventID, result.StatusCode, result.ErrorMessage)
		if delivery.AttemptNumber < 3 {
			deliveryAttemptID, err := scheduleRetry(ctx, svc, delivery)
			if err != nil {
				logrus.WithError(err).Error("Failed to schedule retry for delivery attempt")
				return
			}
			logrus.Infof("Scheduled retry delivery attempt with ID: %d for event ID: %d", deliveryAttemptID, delivery.EventID)
		}
	}
}

// TODO add delay before retrying failed deliveries
// attemptDelivery processes a single delivery attempt by retrieving the associated event and source, merging headers, and sending the HTTP request to the configured egress URL.
func attemptDelivery(svc *service.Service, httpClient *http.Client, delivery db.ListPendingDeliveryAttemptsRow) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // TODO make timeout configurable
	defer cancel()

	if err := markInFlight(ctx, svc, delivery.ID); err != nil {
		logrus.WithError(err).Error("Error updating delivery attempt state to in_flight")
		return
	}

	payload, err := loadDeliveryPayload(ctx, svc, delivery)
	if err != nil {
		logrus.WithError(err).Error("Error loading delivery payload")
		if markErr := markInPending(ctx, svc, delivery.ID); markErr != nil {
			logrus.WithError(err).Error("Error reverting delivery attempt state to pending after load failure:", markErr)
		}
		return
	}

	res, err := sendDeliveryRequest(ctx, httpClient, payload)
	if err != nil {
		logrus.WithError(err).Error("Error sending delivery request")
		if errors.Is(err, context.DeadlineExceeded) {
			// If the context deadline is exceeded, suspend delivery and let the recovery worker resume it.
			return
		}
		var errorType string
		var errorMessage string
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			errorType = "timeout"
			errorMessage = nerr.Error()
		} else if dnsErr, ok := err.(*net.DNSError); ok {
			errorType = "dns_error"
			errorMessage = dnsErr.Error()
		} else if opErr, ok := err.(*net.OpError); ok {
			errorType = "connection_error"
			errorMessage = opErr.Error()
		} else {
			errorType = "network_error"
			errorMessage = err.Error()
		}
		logrus.Warnf("Delivery request failed for event ID %d with error type %s: %s", delivery.EventID, errorType, errorMessage)
		result := &DeliveryResult{
			StatusCode: 0,
			DeliveryState: "failed",
			ErrorType: errorType,
			ErrorMessage: errorMessage,
		}
		handleDeliveryFinalizationAndRetry(ctx, svc, delivery, result)
		return
	}
	defer res.Body.Close()

	result := interpretDeliveryResponse(res)
	handleDeliveryFinalizationAndRetry(ctx, svc, delivery, result)
}