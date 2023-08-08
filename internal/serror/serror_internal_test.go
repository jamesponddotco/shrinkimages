package serror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type errorResponseWriter struct{}

func (*errorResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (*errorResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("forced error")
}

func (*errorResponseWriter) WriteHeader(_ int) {}

func TestJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response ErrorResponse
		writer   http.ResponseWriter
		wantLog  string
	}{
		{
			name: "Normal Error Response",
			response: ErrorResponse{
				Message: "Sample error",
				Code:    http.StatusBadRequest,
			},
			writer:  httptest.NewRecorder(),
			wantLog: "",
		},
		{
			name: "JSON Encoding Error",
			response: ErrorResponse{
				Message: "Sample error",
				Code:    http.StatusBadRequest,
			},
			writer:  &errorResponseWriter{},
			wantLog: "failed to encode error response",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			observedZapCore, logs := observer.New(zap.InfoLevel)
			logger := zap.New(observedZapCore)

			JSON(tt.writer, logger, tt.response)

			if tt.wantLog != "" {
				found := false
				for _, log := range logs.All() {
					if log.Message == tt.wantLog {
						found = true

						break
					}
				}

				if !found {
					t.Errorf("Expected log message not found in logs")
				}
			}
		})
	}
}
