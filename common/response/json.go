package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// StandardResponse represents a common API response format
type StandardResponse[T any] struct {
	Message string `json:"message,omitempty" example:"Success"`
	Data    T      `json:"data,omitempty"`
}

// SendSuccess sends a success response
func SendSuccess[T any](w http.ResponseWriter, status StatusCode, message string, data *T) error {
	if status < 200 || status > 299 {
		return fmt.Errorf("SendSuccess status code must be between 200-299, got %d", status)
	}
	response := StandardResponse[T]{}

	if message != "" {
		response.Message = message
	}
	if data != nil {
		response.Data = *data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))

	logrus.WithFields(logrus.Fields{
		"status_code": status,
		"message":     message,
	}).Debug("sending success response")

	return json.NewEncoder(w).Encode(response)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, status StatusCode, message string) error {
	if status < 400 || status > 599 {
		return fmt.Errorf("SendError status code must be between 400-599, got %d", status)
	}
	response := StandardResponse[interface{}]{}

	if message != "" {
		response.Message = message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))

	logrus.WithFields(logrus.Fields{
		"status_code": status,
		"message":     message,
	}).Error("sending error response")

	return json.NewEncoder(w).Encode(response)
}
