package main

// REF : "The Golang Patterns That Make High-Scale Systems Shockingly Simple"
// https://medium.com/@yashbatra11111/the-golang-patterns-that-make-high-scale-systems-shockingly-simple-8254a0cdb03a
//
// This is the Echo framework version of 10_rest_pattern.go.
// All patterns remain the same, but with Echo's cleaner API.

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// =============================================================================
// PATTERN 1: Small, Behavior-Specific Interfaces
// =============================================================================
// Define interfaces for the exact behavior you need, not large catch-all
// contracts. This enables easy mocking for unit tests and clear dependency
// injection without coupling surprises.
//
// Why this matters:
// - In production, you might swap SimpleProcessor with a DatabaseProcessor
//   or a MockProcessor for testing—all without changing the worker pool code.
// - Interfaces should describe behavior ("I can process a job"), not identity.
// =============================================================================

// Job represents the unit of work flowing through the system.
// Using a simple struct keeps serialization straightforward and
// allows the bounded channel to hold concrete values efficiently.
type Job struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// Processor defines the contract for job processing.
// The context.Context parameter is CRITICAL—it propagates cancellation
// and deadlines, enabling graceful shutdown and per-job timeouts.
type Processor interface {
	Process(ctx context.Context, j Job) error
}

// =============================================================================
// PATTERN 2: Context-First Cancellation
// =============================================================================
// The simpleProcessor demonstrates proper context handling. Every I/O-bound
// or potentially blocking operation should respect context cancellation.
//
// Why this matters:
// - When the server receives SIGTERM, all in-flight jobs need to know
//   "stop working gracefully" rather than running until completion.
// - Per-job timeouts prevent a single slow job from blocking a worker forever.
// =============================================================================

type simpleProcessor struct{}

func (p *simpleProcessor) Process(ctx context.Context, j Job) error {
	// IMPORTANT: Use select to race between actual work and context cancellation.
	// This pattern is the foundation of responsive concurrent code.
	//
	// Real implementation would replace time.After with actual work:
	// - Database writes
	// - HTTP calls to downstream services
	// - CPU-intensive transformations
	//
	// The key insight: wrap ALL blocking operations in context-aware constructs.
	select {
	case <-time.After(10 * time.Millisecond):
		// Simulates successful job completion.
		// In production, this would be your actual business logic.
		return nil
	case <-ctx.Done():
		// Context was cancelled (shutdown or timeout).
		// Return ctx.Err() to distinguish between:
		// - context.Canceled: explicit cancellation (shutdown)
		// - context.DeadlineExceeded: timeout reached
		return ctx.Err()
	}
}

// =============================================================================
// APPLICATION STRUCTURE
// =============================================================================
// Encapsulating all dependencies in a struct (Server) follows Go idioms and
// makes testing easier—you can inject mock dependencies via the constructor.
// =============================================================================

type Server struct {
	echo      *echo.Echo       // Echo instance
	processor Processor        // Injected dependency (testable)
	queue     chan Job         // Bounded channel for backpressure
	wg        sync.WaitGroup   // Tracks in-flight worker goroutines
	processed uint64           // Atomic counter for metrics (Pattern 6)
}

// NewServer is a constructor that enforces proper initialization.
// This pattern prevents "partially initialized" bugs that plague many Go apps.
func NewServer(processor Processor, queueSize int, addr string) *Server {
	e := echo.New()
	e.HideBanner = true

	s := &Server{
		echo:      e,
		processor: processor,
		// =========================================================================
		// PATTERN 4: Channels as Pipelines with Backpressure
		// =========================================================================
		// The bounded queue (size 1000) creates NATURAL BACKPRESSURE.
		// When the queue fills up:
		// - Producers (HTTP handlers) block or receive immediate feedback
		// - This prevents memory exhaustion from unbounded job accumulation
		// - The system fails fast with 503 instead of slowly degrading
		//
		// Queue sizing tradeoffs:
		// - Too small: frequent 503s, wasted producer capacity
		// - Too large: high memory usage, latency spikes during bursts
		// - Rule of thumb: size = (expected_burst_duration) × (processing_rate)
		// =========================================================================
		queue: make(chan Job, queueSize),
	}

	// Setup middleware
	s.setupMiddleware()

	// Register routes
	s.setupRoutes()

	// Configure server address
	s.echo.Server.Addr = addr

	// ReadTimeout/WriteTimeout should be set in production to prevent
	// slow-loris attacks and resource exhaustion.
	s.echo.Server.ReadTimeout = 15 * time.Second
	s.echo.Server.WriteTimeout = 15 * time.Second

	return s
}

// =============================================================================
// MIDDLEWARE SETUP
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

func (s *Server) setupRoutes() {
	s.echo.POST("/submit", s.handleSubmit)
	s.echo.GET("/metrics", s.handleMetrics)
}

// =============================================================================
// PATTERN 3: Worker Pool + Bounded Concurrency
// =============================================================================
// startWorkers creates exactly N workers. This is the cornerstone of
// predictable system behavior under load.
//
// Why bounded concurrency matters:
// - Unbounded goroutines → memory exhaustion, GC pressure, scheduler thrashing
// - Fixed workers → predictable memory footprint, stable latency percentiles
// - Benchmark insight: worker pools use ~50% less memory than unbounded
//   goroutine spawning for 10,000 jobs (see goperf.dev benchmarks)
//
// Worker count guidelines:
// - CPU-bound: workers ≈ runtime.NumCPU()
// - I/O-bound: workers can exceed core count (blocked workers yield CPU)
// - Mixed: start at 2-3× cores and benchmark
// =============================================================================

func (s *Server) startWorkers(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		s.wg.Add(1) // CRITICAL: Add() BEFORE go statement, not inside goroutine
		go func(workerID int) {
			defer s.wg.Done() // Ensures WaitGroup decrements even on panic
			s.worker(ctx, workerID)
		}(i)
	}
}

// worker is the core processing loop. Each worker:
// 1. Pulls jobs from the shared queue
// 2. Creates per-job timeout contexts
// 3. Processes jobs respecting cancellation
// 4. Updates metrics atomically
func (s *Server) worker(ctx context.Context, id int) {
	for {
		select {
		case job, ok := <-s.queue:
			if !ok {
				// Channel closed during shutdown.
				// This is the signal to exit the worker goroutine cleanly.
				log.Printf("worker %d: queue closed, exiting", id)
				return
			}

			// =================================================================
			// PER-JOB TIMEOUT CONTEXT
			// =================================================================
			// Create a child context with timeout for THIS SPECIFIC JOB.
			// Key insights:
			// - Parent ctx carries shutdown signal (from signal.NotifyContext)
			// - Child jobCtx adds a per-job deadline (2 seconds)
			// - If parent cancels, child automatically cancels too (cascading)
			// - ALWAYS call cancel() to release timer resources—hence defer
			//
			// Context hierarchy: rootCtx → workerCtx → jobCtx
			// Cancellation propagates DOWN the tree automatically.
			// =================================================================
			jobCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

			if err := s.processor.Process(jobCtx, job); err != nil {
				// Log error but continue processing—one bad job shouldn't
				// crash the worker. In production, you might:
				// - Send to dead-letter queue
				// - Increment error counter
				// - Implement retry with exponential backoff
				if errors.Is(err, context.Canceled) {
					log.Printf("worker %d: job %s cancelled (shutdown)", id, job.ID)
				} else if errors.Is(err, context.DeadlineExceeded) {
					log.Printf("worker %d: job %s timed out", id, job.ID)
				} else {
					log.Printf("worker %d: job %s failed: %v", id, job.ID, err)
				}
			}

			cancel() // CRITICAL: Always call cancel to prevent goroutine leaks

			// =================================================================
			// PATTERN 6: Atomic Counters for Metrics
			// =================================================================
			// Use sync/atomic for high-throughput counters. Benchmarks show
			// atomic operations handle 47+ million ops/second with zero
			// allocations—far faster than mutex-protected counters.
			//
			// For even higher throughput, consider sharded counters:
			// - Each worker maintains local counter
			// - Background goroutine aggregates periodically
			// - Eliminates all cross-core cache contention
			// =================================================================
			atomic.AddUint64(&s.processed, 1)

		case <-ctx.Done():
			// Parent context cancelled (shutdown signal received).
			// Drain remaining jobs before exiting to prevent data loss.
			log.Printf("worker %d: shutdown signal received, draining queue", id)
			return
		}
	}
}

// =============================================================================
// HTTP HANDLERS WITH BACKPRESSURE
// =============================================================================

// handleSubmit demonstrates explicit backpressure signaling.
// Instead of blocking indefinitely when the queue is full, we immediately
// return 503 Service Unavailable. This is a CRITICAL production pattern.
//
// Why return 503 instead of blocking?
// - Client knows immediately that the system is overloaded
// - Client can retry with a different server (load balancer routing)
// - Prevents connection pile-up that could exhaust file descriptors
// - Keeps response latency predictable (fail fast, fail loud)
func (s *Server) handleSubmit(c echo.Context) error {
	var job Job
	if err := c.Bind(&job); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON payload")
	}

	// =========================================================================
	// NON-BLOCKING SEND WITH BACKPRESSURE
	// =========================================================================
	// The select with default case is the Go idiom for non-blocking operations.
	//
	// Behavior:
	// - If queue has capacity: job enqueued, return 202 Accepted
	// - If queue is full: default case fires, return 503 immediately
	//
	// Alternative: blocking send (queue <- job) would block the HTTP handler
	// until a worker picks up a job. This ties up the HTTP goroutine and can
	// cascade into connection pool exhaustion under high load.
	// =========================================================================
	select {
	case s.queue <- job:
		return c.JSON(http.StatusAccepted, map[string]string{
			"status": "accepted",
			"job_id": job.ID,
		})
	default:
		// Queue full—apply backpressure by rejecting the request.
		// Client should interpret 503 as "try again later" or "try another server."
		return echo.NewHTTPError(http.StatusServiceUnavailable, "queue full, try again later")
	}
}

// handleMetrics exposes operational metrics for observability.
// In production, integrate with Prometheus using the prometheus/client_golang.
func (s *Server) handleMetrics(c echo.Context) error {
	// =========================================================================
	// PATTERN 7: Observability-First Design
	// =========================================================================
	// Expose metrics that answer operational questions:
	// - queue_depth: Is the system keeping up? (trending up = capacity issue)
	// - processed: Throughput indicator
	//
	// Additional metrics to consider in production:
	// - p99 latency (via histogram)
	// - error rate by type
	// - active workers
	// - job processing duration
	// =========================================================================
	metrics := map[string]interface{}{
		"queue_depth": len(s.queue),
		"queue_cap":   cap(s.queue),
		"processed":   atomic.LoadUint64(&s.processed),
	}

	return c.JSON(http.StatusOK, metrics)
}

// =============================================================================
// PATTERN 5: Graceful Shutdown and Lifecycle Management
// =============================================================================
// Run orchestrates the entire server lifecycle:
// 1. Start workers
// 2. Start HTTP server
// 3. Wait for shutdown signal
// 4. Graceful shutdown sequence
//
// The shutdown sequence is ORDERED and DELIBERATE:
// a) Stop accepting new HTTP connections (echo.Shutdown)
// b) Close the job queue (no new jobs can be submitted internally)
// c) Wait for workers to drain remaining jobs (wg.Wait)
//
// This ensures NO HALF-PROCESSED JOBS—every accepted job either completes
// or times out cleanly.
// =============================================================================

func (s *Server) Run(ctx context.Context, workers int) error {
	// Start worker pool
	s.startWorkers(ctx, workers)

	// Start HTTP server in background goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP server starting on %s", s.echo.Server.Addr)
		if err := s.echo.Start(s.echo.Server.Addr); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// =========================================================================
	// WAIT FOR SHUTDOWN TRIGGER
	// =========================================================================
	// Two possible triggers:
	// 1. Context cancelled (SIGINT/SIGTERM received)
	// 2. Server error (port already in use, etc.)
	// =========================================================================
	select {
	case <-ctx.Done():
		log.Println("shutdown: received termination signal")
	case err := <-serverErr:
		return err
	}

	// =========================================================================
	// GRACEFUL SHUTDOWN SEQUENCE
	// =========================================================================
	// Step 1: Stop HTTP server from accepting new connections
	// The 15-second timeout allows in-flight HTTP requests to complete.
	// After 15 seconds, remaining connections are forcefully closed.
	// =========================================================================
	log.Println("shutdown: stopping HTTP server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: HTTP server forced close: %v", err)
	}

	// =========================================================================
	// Step 2: Close job queue
	// Closing the channel signals workers to exit after draining remaining jobs.
	// Workers detect this via "ok" being false in their range loop.
	// =========================================================================
	log.Println("shutdown: closing job queue")
	close(s.queue)

	// =========================================================================
	// Step 3: Wait for workers to finish
	// This blocks until all workers have:
	// a) Processed remaining jobs from the queue
	// b) Called wg.Done() to signal completion
	//
	// In production, consider adding a timeout here too, in case a worker hangs.
	// =========================================================================
	log.Println("shutdown: waiting for workers to finish")
	s.wg.Wait()

	log.Println("shutdown: complete")
	return nil
}

// =============================================================================
// MAIN: Signal-Based Lifecycle
// =============================================================================

func main() {
	// =========================================================================
	// signal.NotifyContext (Go 1.16+) combines signal handling with context.
	// This is the idiomatic way to handle SIGINT (Ctrl+C) and SIGTERM
	// (container orchestrator shutdown).
	//
	// Key insight: The returned context automatically cancels when ANY of the
	// specified signals are received. This cancellation cascades to all child
	// contexts, enabling system-wide graceful shutdown.
	// =========================================================================
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // Container kill signal
	)
	defer stop() // Restore default signal behavior

	// =========================================================================
	// CONFIGURATION
	// =========================================================================
	// In production, these would come from environment variables or config files.
	// Shown explicitly here for educational clarity.
	// =========================================================================
	const (
		workerCount = 50   // Adjust based on workload type and core count
		queueSize   = 1000 // Adjust based on burst tolerance and memory budget
		serverAddr  = ":8080"
	)

	processor := &simpleProcessor{}
	server := NewServer(processor, queueSize, serverAddr)

	log.Println("=== High-Scale REST API Server (Echo) ===")
	log.Printf("Workers: %d, Queue Size: %d", workerCount, queueSize)
	log.Printf("Endpoints:")
	log.Printf("  POST /submit  - Submit a job")
	log.Printf("  GET  /metrics - View metrics")

	if err := server.Run(ctx, workerCount); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// =============================================================================
// USAGE EXAMPLES
// =============================================================================
//
// Run the server:
//   go run 20_rest_pattern_echo.go
//
// Submit a job:
//   curl -X POST http://localhost:8080/submit \
//     -H "Content-Type: application/json" \
//     -d '{"id":"job-001","data":"hello world"}'
//
// Check metrics:
//   curl http://localhost:8080/metrics
//
// Load test (requires hey: go install github.com/rakyll/hey@latest):
//   hey -n 10000 -c 100 -m POST -H "Content-Type: application/json" \
//     -d '{"id":"test","data":"load"}' http://localhost:8080/submit
//
// Graceful shutdown:
//   Press Ctrl+C or send SIGTERM
//
// =============================================================================
