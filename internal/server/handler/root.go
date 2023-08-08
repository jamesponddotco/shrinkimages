package handler

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/config"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/endpoint"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"go.uber.org/zap"
)

// RootHandler is an HTTP handler for the / endpoint.
type RootHandler struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewRootHandler creates a new instance of RootHandler.
func NewRootHandler(cfg *config.Config, logger *zap.Logger) *RootHandler {
	return &RootHandler{
		cfg:    cfg,
		logger: logger,
	}
}

// ServeHTTP handles HTTP requests for the / endpoint.
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case endpoint.Root:
		http.Redirect(w, r, h.cfg.Service.Homepage, http.StatusMovedPermanently)
	default:
		serror.JSON(w, h.logger, serror.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Page not found. Please check the URL and try again.",
		})
	}
}
