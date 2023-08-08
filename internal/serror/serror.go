// Package serror provides custom error types for the service.
package serror

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// ErrorResponse is the response returned by the API when an error occurs.
type ErrorResponse struct {
	// Message is a human-readable message describing the error.
	Message string `json:"message"`

	// Documentation is a URL to the documentation with more information about
	// the error.
	Documentation string `json:"documentation,omitempty"`

	// Code is a machine-readable code describing the error.
	Code uint `json:"code"`
}

// JSON sends an ErrorResponse to the HTTP response writer as JSON.
func JSON(w http.ResponseWriter, logger *zap.Logger, response ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(response.Code))

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode error response", zap.String("error", err.Error()))
	}
}
