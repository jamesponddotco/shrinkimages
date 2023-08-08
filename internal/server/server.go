// Package server provides a simple HTTP server for the service.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/config"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/endpoint"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/fetch"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/serror"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/server/handler"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/server/middleware"
	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
	"go.uber.org/zap"
)

// Server represents the Shrink Images server.
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// New creates a new Shrink Images server instance.
func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(cfg.Server.TLS.Certificate, cfg.Server.TLS.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	var tlsConfig *tls.Config

	if cfg.Server.TLS.Version == "1.3" {
		tlsConfig = xtls.ModernServerConfig()
	}

	if cfg.Server.TLS.Version == "1.2" {
		tlsConfig = xtls.IntermediateServerConfig()
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	middlewares := []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return middleware.PanicRecovery(logger, h) },
		func(h http.Handler) http.Handler { return middleware.UserAgent(logger, h) },
		func(h http.Handler) http.Handler { return middleware.Authorization(cfg.Service.APIKey, logger, h) },
		func(h http.Handler) http.Handler {
			return middleware.AcceptRequests(
				[]string{
					http.MethodPost,
					http.MethodGet,
					http.MethodHead,
				},
				logger,
				h,
			)
		},
		func(h http.Handler) http.Handler { return middleware.PrivacyPolicy(cfg.Service.PrivacyPolicy, h) },
		func(h http.Handler) http.Handler { return middleware.TermsOfService(cfg.Service.TermsOfService, h) },
	}

	var (
		fetchInstance = fetch.New(cfg.Service.Name, cfg.Service.Contact)
		pingHandler   = handler.NewPingHandler(logger)
		shrinkHandler = handler.NewShrinkHandler(cfg, fetchInstance, logger)
	)

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint.Root, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case endpoint.Root:
			http.Redirect(w, r, cfg.Service.Homepage, http.StatusMovedPermanently)
		default:
			serror.JSON(w, logger, serror.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Page not found. Check the URL and try again.",
			})
		}
	})

	mux.Handle(endpoint.Ping, middleware.Chain(pingHandler, middlewares...))
	mux.Handle(endpoint.Shrink, middleware.Chain(shrinkHandler, middlewares...))

	httpServer := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

// Start starts the Shrink Images server.
func (s *Server) Start() error {
	var (
		sigint            = make(chan os.Signal, 1)
		shutdownCompleted = make(chan struct{})
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server Shutdown:", zap.Error(err))
		}

		close(shutdownCompleted)
	}()

	if err := s.httpServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	<-shutdownCompleted

	return nil
}

// Stop gracefully shuts down the Shrink Images server.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
