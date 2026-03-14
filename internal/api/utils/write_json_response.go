package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON is a helper function to write a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, statusCode int, data any) error {
	response, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to marshal response data: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(response); err != nil {
		return fmt.Errorf("Failed to write response: %w", err)
	}
	return nil
}
