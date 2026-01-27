package main

// Simple Echo Framework REST API Server - Server Initialization and Graceful Shutdown
//
// This is a beginner-friendly example demonstrating the essential patterns
// for building a production-ready HTTP server using Echo framework.
// It focuses on:
// - Proper server initialization with Echo
// - Signal-based graceful shutdown
// - Context propagation for cancellation
//
// Compared to 10_simple_svr.go (standard library), Echo provides:
// - Cleaner routing with method-specific handlers
// - Built-in middleware (recovery, logging)
// - Automatic content-type negotiation
// - Request binding for JSON payloads
// - Better error handling with HTTPError
//
// For more advanced patterns (worker pools, backpressure, sync.Pool),
// see A02.rest_api_pattern and A03.high_perf_rest_svr.

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// =============================================================================
// SERVER STRUCTURE
// =============================================================================
// Encapsulating dependencies in a struct is idiomatic Go.
// It makes the code testable and dependencies explicit.
// =============================================================================

type Server struct {
	echo *echo.Echo
	addr string
}

// NewServer creates and configures a new Echo server.
// This constructor pattern ensures proper initialization.
func NewServer(addr string) *Server {
	s := &Server{
		echo: echo.New(),
		addr: addr,
	}

	// Configure Echo settings
	s.echo.HideBanner = true // Cleaner log output

	// =================================================================
	// HTTP Server Timeouts (Important for Production)
	// =================================================================
	// Without these, malicious clients can hold connections forever.
	// Echo wraps http.Server, so we configure timeouts on the underlying server.
	// - ReadTimeout: max time to read entire request including body
	// - WriteTimeout: max time to write response
	// - IdleTimeout: max time for keep-alive connections
	// =================================================================
	s.echo.Server.ReadTimeout = 15 * time.Second
	s.echo.Server.WriteTimeout = 15 * time.Second
	s.echo.Server.IdleTimeout = 60 * time.Second

	// Setup middleware
	s.setupMiddleware()

	// Register routes
	s.setupRoutes()

	return s
}

// =============================================================================
// MIDDLEWARE SETUP
// =============================================================================
// Middleware order matters! The execution order is:
// 1. Recover (catches panics from all subsequent middleware)
// 2. Logger (logs all requests)
// 3. Handler
// =============================================================================

func (s *Server) setupMiddleware() {
	// Recovery middleware - converts panics to 500 errors
	// MUST be first to catch panics from subsequent middleware
	s.echo.Use(middleware.Recover())

	// Logger middleware - logs HTTP requests
	s.echo.Use(middleware.Logger())
}

// =============================================================================
// ROUTE REGISTRATION
// =============================================================================
// Echo uses method-specific registration: e.GET(), e.POST(), etc.
// This is cleaner than http.HandleFunc which requires method checks inside handlers.
// =============================================================================

func (s *Server) setupRoutes() {
	s.echo.GET("/", s.handleRoot)
	s.echo.GET("/health", s.handleHealth)
	s.echo.POST("/api/echo", s.handleEcho)
}

// =============================================================================
// HTTP HANDLERS
// =============================================================================
// Echo handlers return error instead of writing directly to response.
// This enables cleaner error handling and automatic content-type negotiation.
// =============================================================================

// handleRoot returns a welcome message
func (s *Server) handleRoot(c echo.Context) error {
	// Echo automatically handles 404 for unmatched routes,
	// unlike net/http where "/" matches all paths
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Welcome to Simple Echo REST API Server",
		"version": "1.0.0",
	})
}

// handleHealth provides a health check endpoint.
// Used by load balancers and container orchestrators (e.g., Kubernetes).
func (s *Server) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleEcho echoes back the received JSON payload.
// Demonstrates basic request body parsing and response.
func (s *Server) handleEcho(c echo.Context) error {
	// Echo's Bind() automatically parses JSON based on Content-Type header
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		// Echo's HTTPError provides status code and message
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON payload")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"received": payload,
		"echo":     true,
	})
}

// =============================================================================
// SERVER LIFECYCLE
// =============================================================================

// Run starts the server and handles graceful shutdown.
// This is the heart of proper server lifecycle management.
func (s *Server) Run(ctx context.Context) error {
	// Channel to capture server errors
	serverErr := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on %s", s.addr)
		if err := s.echo.Start(s.addr); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// =======================================================================
	// WAIT FOR SHUTDOWN SIGNAL OR ERROR
	// =======================================================================
	// Two possible triggers for shutdown:
	// 1. Context cancelled (SIGINT/SIGTERM received)
	// 2. Server error (e.g., port already in use)
	// =======================================================================
	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErr:
		return err
	}

	// =======================================================================
	// GRACEFUL SHUTDOWN
	// =======================================================================
	// echo.Shutdown() gracefully shuts down the server:
	// - Stops accepting new connections
	// - Waits for existing connections to finish
	// - Returns when all connections are closed or timeout reached
	//
	// The timeout prevents hanging forever if clients don't disconnect.
	// =======================================================================
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
		return err
	}

	// =======================================================================
	// RESOURCE CLEANUP
	// =======================================================================
	// Place resource cleanup logic HERE, after Shutdown() completes.
	// At this point:
	// - No new requests are being accepted
	// - All in-flight requests have completed
	// - Safe to close shared resources
	//
	// Example cleanup (if Server had these dependencies):
	//   s.db.Close()           // Close database connections
	//   s.cache.Close()        // Close Redis/cache connections
	//   s.messageQueue.Close() // Close message queue connections
	//   s.logger.Sync()        // Flush log buffers
	//
	// Why this order matters:
	// 1. Shutdown HTTP first → Stops accepting new requests
	// 2. Wait for in-flight requests → They may still need DB/cache
	// 3. Then cleanup resources → No active requests using them
	// =======================================================================

	log.Println("Server shutdown complete")
	return nil
}

// =============================================================================
// MAIN FUNCTION
// =============================================================================

func main() {
	// =========================================================================
	// SIGNAL HANDLING WITH CONTEXT
	// =========================================================================
	// signal.NotifyContext (Go 1.16+) is the idiomatic way to handle OS signals.
	// It creates a context that automatically cancels when signals are received.
	//
	// SIGINT:  Triggered by Ctrl+C
	// SIGTERM: Triggered by container orchestrators (Docker, K8s) during shutdown
	// =========================================================================
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // Container shutdown signal
	)
	defer stop() // Restore default signal handling on exit

	// Configuration
	addr := ":8080"

	// Create and run server
	server := NewServer(addr)

	log.Println("=== Simple Echo REST API Server ===")
	log.Printf("Endpoints:")
	log.Printf("  GET  /         - Welcome message")
	log.Printf("  GET  /health   - Health check")
	log.Printf("  POST /api/echo - Echo JSON payload")

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// =============================================================================
// USAGE EXAMPLES
// =============================================================================
//
// Install Echo (if not already in go.mod):
//   go get github.com/labstack/echo/v4
//
// Run the server:
//   go run 20_simple_echo_svr.go
//
// Test endpoints:
//   curl http://localhost:8080/
//   curl http://localhost:8080/health
//   curl -X POST http://localhost:8080/api/echo -H "Content-Type: application/json" -d '{"hello":"world"}'
//
// Graceful shutdown:
//   Press Ctrl+C or send SIGTERM
//
// =============================================================================
