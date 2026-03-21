package server

import (
	"context"
	"net/http"
)

// Server wraps http.Server with Start and Shutdown helpers.
type Server struct {
	httpServer *http.Server
}

// New creates a Server bound to the given address with the provided handler.
func New(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Start begins listening and serving HTTP requests. It blocks until an error occurs.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server, waiting up to the deadline in ctx.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
