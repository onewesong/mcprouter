package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// APIServer as the proxy server for API request
type APIServer struct {
	server *echo.Echo // http server built with echo
}

// NewAPIServer will create API server
func NewAPIServer() *APIServer {
	return &APIServer{
		server: echo.New(),
	}
}

// Route will create the routes for http server
func (s *APIServer) Route(route func(e *echo.Echo)) {
	s.server.Validator = NewValidator()
	s.server.Use(middleware.Logger())
	s.server.Use(createAPIMiddleware())

	route(s.server)
}

// Start will start the http server
func (s *APIServer) Start(port int) {
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%d", port)))
}

// StartWithContext will start the http server with context support for graceful shutdown
func (s *APIServer) StartWithContext(ctx context.Context, port int) error {
	// Start server
	go func() {
		if err := s.server.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			s.server.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctxShutDown); err != nil {
		s.server.Logger.Fatal(err)
	}

	return nil
}

// Shutdown will gracefully shutdown the server
func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
