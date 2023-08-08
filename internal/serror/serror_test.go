package serror_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"go.uber.org/zap"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		response serror.ErrorResponse
		want     string
		wantCode int
	}{
		{
			name: "Sample Error Response",
			response: serror.ErrorResponse{
				Message:       "Sample error",
				Documentation: "https://example.com/docs",
				Code:          http.StatusBadRequest,
			},
			want:     `{"message":"Sample error","documentation":"https://example.com/docs","code":400}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()

			serror.JSON(rec, logger, tt.response)

			if got := rec.Body.String(); strings.TrimSpace(got) != tt.want {
				t.Errorf("JSON() = %v, want %v", got, tt.want)
			}

			if got := rec.Code; got != tt.wantCode {
				t.Errorf("Response Code = %v, want %v", got, tt.wantCode)
			}

			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content Type = %v, want 'application/json'", got)
			}
		})
	}
}
