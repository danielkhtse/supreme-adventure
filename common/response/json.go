package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// StandardResponse represents the standard response format for successful API calls
type StandardResponse[T any] struct {
	// Message is an optional success message
	Message string `json:"message,omitempty" example:"Operation completed successfully"`
	// Data contains the optional response payload
	Data T `json:"data,omitempty"`
}

// ErrorResponse represents the standard response format for API errors
type ErrorResponse struct {
	// Message contains details about what went wrong
	Message string `json:"message" example:"Invalid request parameters"`
}

// SendSuccess sends a success response with optional data
//
// @Summary Send success response
// @Description Sends a standardized success response with optional data payload
// @Tags response
// @Accept json
// @Produce json
// @Param w body http.ResponseWriter true "HTTP response writer"
// @Param status path StatusCode true "HTTP status code (200-299)"
// @Param data body T false "Optional response data"
// @Success 200 {object} StandardResponse[T] "Success response with data"
// @Failure 500 {object} ErrorResponse "Invalid status code error"
// @Return error
func SendSuccess[T any](w http.ResponseWriter, status StatusCode, data *T) error {
	if status < 200 || status > 299 {
		return fmt.Errorf("SendSuccess status code must be between 200-299, got %d", status)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))

	logrus.WithFields(logrus.Fields{
		"status_code": status,
	}).Debug("sending success response")

	if data == nil {
		return json.NewEncoder(w).Encode(struct{}{})
	}

	response := StandardResponse[T]{
		Data: *data,
	}
	return json.NewEncoder(w).Encode(response.Data)
}

// SendError sends an error response with a message
//
// @Summary Send error response
// @Description Sends a standardized error response with message
// @Tags response
// @Accept json
// @Produce json
// @Param w body http.ResponseWriter true "HTTP response writer"
// @Param status path StatusCode true "HTTP status code (400-599)"
// @Param message body string false "Error message"
// @Success 400-599 {object} ErrorResponse "Error response"
// @Failure 500 {object} ErrorResponse "Invalid status code error"
// @Return error
func SendError(w http.ResponseWriter, status StatusCode, message string) error {
	if status < 400 || status > 599 {
		return fmt.Errorf("SendError status code must be between 400-599, got %d", status)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))

	logrus.WithFields(logrus.Fields{
		"status_code": status,
		"message":     message,
	}).Error("sending error response")

	if message == "" {
		return json.NewEncoder(w).Encode(struct{}{})
	}

	response := ErrorResponse{
		Message: message,
	}
	return json.NewEncoder(w).Encode(response)
}
