package handler

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"go.uber.org/zap"
)

// Pong is the response for the /ping endpoint.
const Pong string = "pong"

// PingHandler is an HTTP handler for the /ping endpoint.
type PingHandler struct {
	logger *zap.Logger
}

// NewPingHandler returns a new instance of PingHandler.
func NewPingHandler(logger *zap.Logger) *PingHandler {
	return &PingHandler{
		logger: logger,
	}
}

// ServeHTTP serves the /ping endpoint.
func (h *PingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(xhttp.ContentType, xhttp.TextPlain)

	_, err := w.Write([]byte(Pong))
	if err != nil {
		h.logger.Error("failed to write response", zap.Error(err))

		serror.JSON(w, h.logger, serror.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to write response. Please try again later.",
		})

		return
	}
}
