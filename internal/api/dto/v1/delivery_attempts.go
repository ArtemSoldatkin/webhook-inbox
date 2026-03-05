package dtov1

import "time"

// DeliveryAttemptDTO represents the data transfer object for a delivery attempt in the API.
type DeliveryAttemptDTO struct {
	ID            int64      `json:"id"`
	EventID       int64      `json:"event_id"`
	AttemptNumber int32      `json:"attempt_number"`
	State         string     `json:"state"`
	StatusCode    *int32     `json:"status_code,omitempty"`
	ErrorType     *string    `json:"error_type,omitempty"`
	ErrorMessage  *string    `json:"error_message,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	NextAttemptAt *time.Time `json:"next_attempt_at,omitempty"`
}
