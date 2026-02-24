package dtov1

import "time"

// DeliveryAttemptDTO represents the data transfer object for a delivery attempt in the API.
type DeliveryAttemptDTO struct {
	ID            int64
	EventID       int64
	AttemptNumber int32
	State         string
	StatusCode    *int32
	ErrorType     *string
	ErrorMessage  *string
	StartedAt     *time.Time
	FinishedAt    *time.Time
	CreatedAt     time.Time
	NextAttemptAt *time.Time
}
