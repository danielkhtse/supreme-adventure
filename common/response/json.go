package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StandardResponse represents a common API response format
type StandardResponse[T any] struct {
	Message string     `json:"message" example:"Success"`
	Data    T          `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    StatusCode `json:"code" example:"400"`
	Message string     `json:"message" example:"Invalid input parameters"`
}

// SendSuccess sends a success response
func SendSuccess[T any](w http.ResponseWriter, status StatusCode, message string, data *T) error {
	if status < 200 || status > 299 {
		return fmt.Errorf("SendSuccess status code must be between 200-299, got %d", status)
	}
	response := StandardResponse[T]{
		Message: message,
	}
	if data != nil {
		response.Data = *data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	return json.NewEncoder(w).Encode(response)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, status StatusCode, message string) error {
	if status < 400 || status > 599 {
		return fmt.Errorf("SendError status code must be between 400-599, got %d", status)
	}
	response := StandardResponse[interface{}]{
		Message: "Error",
		Error: &ErrorInfo{
			Code:    status,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	return json.NewEncoder(w).Encode(response)
}
