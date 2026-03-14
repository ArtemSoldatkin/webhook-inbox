package api

import (
	"encoding/json"
	"net/http"
)

// JSONMarshalError represents an error that occurred during JSON marshaling.
type JSONMarshalError struct {
	Message string `json:"error"`
	Err     error  `json:"-"`
}

// Error returns the error message for JSONMarshalError.
func (e *JSONMarshalError) Error() string {
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work with JSONMarshalError.
func (e *JSONMarshalError) Unwrap() error { return e.Err }

// JSONWriteError represents an error that occurred while writing the JSON response.
type JSONWriteError struct {
	Message string `json:"error"`
	Err     error  `json:"-"`
}

// Error returns the error message for JSONWriteError.
func (e *JSONWriteError) Error() string {
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work with JSONWriteError.
func (e *JSONWriteError) Unwrap() error { return e.Err }

// JSON is a helper function to write a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, statusCode int, data any) error {
	response, err := json.Marshal(data)
	if err != nil {
		return &JSONMarshalError{
			Message: "failed to marshal response data",
			Err:     err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(response); err != nil {
		return &JSONWriteError{
			Message: "failed to write response",
			Err:     err,
		}
	}
	return nil
}
