// Package server wires the net/http server and owns graceful shutdown logic.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"my-app/internal/config"
)

// ShutdownTimeout is the maximum time allowed for in-flight requests to finish
// before the process exits. Exported so main.go can size its context correctly.
const ShutdownTimeout = 30 * time.Second

// Server wraps the standard library HTTP server.
type Server struct {
	httpServer *http.Server
}

// New constructs a Server from config and the application's root handler.
func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.AppPort,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start begins listening and serving. It returns a non-nil error only when the
// server exits unexpectedly (not on graceful shutdown via http.ErrServerClosed).
func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server within the deadline carried by ctx.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
