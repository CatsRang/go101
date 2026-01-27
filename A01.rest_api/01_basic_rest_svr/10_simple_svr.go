package main

// Simple REST API Server - Server Initialization and Graceful Shutdown
//
// This is a beginner-friendly example demonstrating the essential patterns
// for building a production-ready HTTP server in Go. It focuses on:
// - Proper server initialization
// - Signal-based graceful shutdown
// - Context propagation for cancellation
//
// For more advanced patterns (worker pools, backpressure, sync.Pool),
// see A02.rest_api_pattern and A03.high_perf_rest_svr.

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// =============================================================================
// SERVER STRUCTURE
// =============================================================================
// Encapsulating dependencies in a struct is idiomatic Go.
// It makes the code testable and dependencies explicit.
// =============================================================================

type Server struct {
	srv    *http.Server
	router *http.ServeMux
}

// NewServer creates and configures a new HTTP server.
// This constructor pattern ensures proper initialization.
func NewServer(addr string) *Server {
	s := &Server{
		router: http.NewServeMux(),
	}

	// Register routes
	s.registerRoutes()

	// Configure HTTP server with production timeouts
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
		// =================================================================
		// HTTP Server Timeouts (Important for Production)
		// =================================================================
		// Without these, malicious clients can hold connections forever.
		// - ReadTimeout: max time to read entire request including body
		// - WriteTimeout: max time to write response
		// - IdleTimeout: max time for keep-alive connections
		// - ReadHeaderTimeout: prevents Slowloris attacks
		// =================================================================
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

// =============================================================================
// ROUTE REGISTRATION
// =============================================================================

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/", s.handleRoot)
	s.router.HandleFunc("/health", s.handleHealth)
	s.router.HandleFunc("/api/echo", s.handleEcho)
}

// =============================================================================
// HTTP HANDLERS
// =============================================================================

// handleRoot returns a welcome message
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]string{
		"message": "Welcome to Simple REST API Server",
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleHealth provides a health check endpoint.
// Used by load balancers and container orchestrators (e.g., Kubernetes).
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleEcho echoes back the received JSON payload.
// Demonstrates basic request body parsing and response.
func (s *Server) handleEcho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"received": payload,
		"echo":     true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		log.Printf("Server starting on %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	// srv.Shutdown() gracefully shuts down the server:
	// - Stops accepting new connections
	// - Waits for existing connections to finish
	// - Returns when all connections are closed or timeout reached
	//
	// The timeout prevents hanging forever if clients don't disconnect.
	// =======================================================================
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
		return err
	}

	// =======================================================================
	// RESOURCE CLEANUP
	// =======================================================================
	// Place resource cleanup logic HERE, after srv.Shutdown() completes.
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

	log.Println("=== Simple REST API Server ===")
	log.Printf("Endpoints:")
	log.Printf("  GET  /        - Welcome message")
	log.Printf("  GET  /health  - Health check")
	log.Printf("  POST /api/echo - Echo JSON payload")

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// =============================================================================
// USAGE EXAMPLES
// =============================================================================
//
// Run the server:
//   go run 10_simple_svr.go
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
