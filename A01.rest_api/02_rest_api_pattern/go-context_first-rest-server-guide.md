# Go High-Scale System Patterns: A Comprehensive Reference Guide

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Core Language Primitives](#core-language-primitives)
   - 2.1 [context.Context](#21-contextcontext)
   - 2.2 [Channels](#22-channels)
   - 2.3 [sync.WaitGroup](#23-syncwaitgroup)
   - 2.4 [sync/atomic](#24-syncatomic)
3. [Concurrency Patterns](#concurrency-patterns)
   - 3.1 [Worker Pool Pattern](#31-worker-pool-pattern)
   - 3.2 [Pipeline Pattern](#32-pipeline-pattern)
   - 3.3 [Semaphore Pattern](#33-semaphore-pattern)
   - 3.4 [Fan-Out/Fan-In Pattern](#34-fan-outfan-in-pattern)
4. [Resilience Patterns](#resilience-patterns)
   - 4.1 [Backpressure](#41-backpressure)
   - 4.2 [Circuit Breaker](#42-circuit-breaker)
   - 4.3 [Graceful Shutdown](#43-graceful-shutdown)
5. [Observability Tools](#observability-tools)
   - 5.1 [Structured Logging (slog)](#51-structured-logging-slog)
   - 5.2 [Prometheus Metrics](#52-prometheus-metrics)
   - 5.3 [pprof Profiling](#53-pprof-profiling)
6. [Advanced Patterns](#advanced-patterns)
   - 6.1 [errgroup Package](#61-errgroup-package)
   - 6.2 [Sharded Counters](#62-sharded-counters)
7. [External Libraries Reference](#external-libraries-reference)
8. [Quick Reference Cheat Sheet](#quick-reference-cheat-sheet)

---

## Executive Summary

This guide provides a comprehensive reference for the tools, frameworks, and architectural patterns used to build high-scale systems in Go. The patterns covered enable:

- **Bounded Concurrency**: Prevent resource exhaustion through worker pools and semaphores
- **Predictable Behavior**: Context-driven cancellation and graceful shutdown
- **Resilience**: Backpressure mechanisms and circuit breakers for fault tolerance
- **Observability**: Structured logging, metrics, and profiling for production debugging

| Pattern Category | Key Tools | Primary Benefit |
|-----------------|-----------|-----------------|
| Concurrency Control | Worker Pool, Semaphore | Predictable resource usage |
| Error Propagation | context, errgroup | Cascading cancellation |
| Flow Control | Bounded Channels | Natural backpressure |
| Fault Tolerance | Circuit Breaker | Prevent cascade failures |
| Observability | pprof, Prometheus, slog | Production debugging |

---

## Core Language Primitives

### 2.1 context.Context

The `context` package is the foundation for cancellation, deadlines, and request-scoped values in Go. It should be the **first parameter** of any function that does I/O or spawns goroutines.

#### Core Functions

| Function | Purpose | When to Use |
|----------|---------|-------------|
| `context.Background()` | Root context, never cancelled | Entry points (main, init) |
| `context.TODO()` | Placeholder when unsure | Refactoring, temporary |
| `context.WithCancel(parent)` | Manual cancellation | Explicit shutdown control |
| `context.WithTimeout(parent, duration)` | Auto-cancel after duration | Per-request timeouts |
| `context.WithDeadline(parent, time)` | Auto-cancel at specific time | Absolute deadlines |
| `context.WithValue(parent, key, val)` | Pass request-scoped data | Trace IDs, auth tokens |

#### Example: Context with Timeout

```go
func FetchData(ctx context.Context, url string) ([]byte, error) {
    // Create child context with 5-second timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel() // ALWAYS call cancel to release resources

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        // Check if error was due to cancellation
        if errors.Is(err, context.Canceled) {
            return nil, fmt.Errorf("request cancelled: %w", err)
        }
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("request timed out: %w", err)
        }
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

#### Best Practices

1. **Always pass context as first parameter**: `func DoWork(ctx context.Context, ...)`
2. **Never store context in structs**: Pass explicitly to methods
3. **Always call cancel()**: Use `defer cancel()` immediately after creation
4. **Check ctx.Done() in loops**: Enable responsive cancellation
5. **Use WithValue sparingly**: Only for request-scoped data crossing API boundaries

#### Context Cancellation Flow

```
                    ┌─────────────────┐
                    │ context.Background() │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │ WithCancel(ctx) │ ← Server-level context
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
     ┌────────▼────────┐ ┌───▼───┐ ┌───────▼───────┐
     │ WithTimeout(ctx)│ │Worker2│ │WithTimeout(ctx)│
     │   Request A     │ │       │ │   Request B    │
     └─────────────────┘ └───────┘ └────────────────┘

     Cancellation cascades DOWN the tree automatically
```

---

### 2.2 Channels

Channels are Go's primary mechanism for communication between goroutines. For high-scale systems, **bounded (buffered) channels** are essential for backpressure.

#### Channel Types

| Type | Declaration | Behavior |
|------|-------------|----------|
| Unbuffered | `make(chan T)` | Synchronous; blocks until receiver ready |
| Buffered | `make(chan T, n)` | Asynchronous up to capacity n |
| Send-only | `chan<- T` | Can only send to channel |
| Receive-only | `<-chan T` | Can only receive from channel |

#### Example: Bounded Channel for Backpressure

```go
func ProcessJobs(ctx context.Context) {
    // Bounded queue creates natural backpressure
    queue := make(chan Job, 1000)

    // Producer
    go func() {
        for job := range incomingJobs {
            select {
            case queue <- job:
                // Job accepted
            default:
                // Queue full - apply backpressure
                log.Println("queue full, rejecting job")
            case <-ctx.Done():
                return
            }
        }
    }()

    // Consumer
    for job := range queue {
        process(job)
    }
}
```

#### Channel Patterns

```go
// 1. Generator Pattern - produces values on demand
func Generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// 2. Multiplexer (Fan-In) - merge multiple channels
func Merge(cs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, c := range cs {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

---

### 2.3 sync.WaitGroup

`sync.WaitGroup` is a counting semaphore for waiting on a collection of goroutines to complete.

#### Core Methods

| Method | Purpose |
|--------|---------|
| `Add(delta int)` | Increment counter (call BEFORE go statement) |
| `Done()` | Decrement counter (call with defer inside goroutine) |
| `Wait()` | Block until counter reaches zero |

#### Example: Worker Pool with WaitGroup

```go
func ProcessItems(items []Item) {
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)  // CRITICAL: Add BEFORE go statement
        go func(i Item) {
            defer wg.Done()  // CRITICAL: Use defer for safety
            process(i)
        }(item)
    }
    
    wg.Wait()  // Block until all goroutines complete
    fmt.Println("All items processed")
}
```

#### Common Mistakes to Avoid

```go
// ❌ WRONG: Add() inside goroutine (race condition)
go func() {
    wg.Add(1)
    defer wg.Done()
    // ...
}()

// ✅ CORRECT: Add() before goroutine
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()

// ❌ WRONG: Forgetting Done() on early return
go func() {
    if err != nil {
        return  // wg.Done() never called!
    }
    wg.Done()
}()

// ✅ CORRECT: Use defer for guaranteed execution
go func() {
    defer wg.Done()
    if err != nil {
        return  // defer still executes
    }
}()
```

---

### 2.4 sync/atomic

The `sync/atomic` package provides low-level atomic memory primitives for lock-free operations. Use for simple counters and flags in high-contention scenarios.

#### Common Operations

| Operation | Function | Description |
|-----------|----------|-------------|
| Add | `atomic.AddInt64(&val, delta)` | Atomically add and return new value |
| Load | `atomic.LoadInt64(&val)` | Atomically read value |
| Store | `atomic.StoreInt64(&val, new)` | Atomically write value |
| Swap | `atomic.SwapInt64(&val, new)` | Atomically swap and return old value |
| CompareAndSwap | `atomic.CompareAndSwapInt64(&val, old, new)` | CAS operation |

#### Example: High-Throughput Counter

```go
type Metrics struct {
    requestCount  uint64
    errorCount    uint64
    bytesReceived uint64
}

func (m *Metrics) IncrementRequests() {
    atomic.AddUint64(&m.requestCount, 1)
}

func (m *Metrics) IncrementErrors() {
    atomic.AddUint64(&m.errorCount, 1)
}

func (m *Metrics) AddBytes(n uint64) {
    atomic.AddUint64(&m.bytesReceived, n)
}

func (m *Metrics) GetStats() (requests, errors, bytes uint64) {
    return atomic.LoadUint64(&m.requestCount),
           atomic.LoadUint64(&m.errorCount),
           atomic.LoadUint64(&m.bytesReceived)
}
```

#### When to Use Atomics vs Mutex

| Use Case | Recommendation |
|----------|---------------|
| Simple counter increment | `atomic.AddInt64` |
| Single value read/write | `atomic.Load/Store` |
| Multiple related values | `sync.Mutex` |
| Complex data structures | `sync.Mutex` or `sync.RWMutex` |
| Read-heavy workloads | `sync.RWMutex` |

---

## Concurrency Patterns

### 3.1 Worker Pool Pattern

A fixed number of workers process jobs from a shared queue. This is the cornerstone of predictable concurrency.

#### Why Worker Pools?

- **Bounded memory**: Fixed goroutine count prevents memory exhaustion
- **Predictable latency**: No goroutine scheduling overhead spikes
- **Resource control**: Limit connections to databases, APIs, etc.

#### Implementation

```go
type WorkerPool struct {
    workers   int
    jobQueue  chan Job
    wg        sync.WaitGroup
    processor Processor
}

func NewWorkerPool(workers, queueSize int, p Processor) *WorkerPool {
    return &WorkerPool{
        workers:   workers,
        jobQueue:  make(chan Job, queueSize),
        processor: p,
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx, i)
    }
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case job, ok := <-wp.jobQueue:
            if !ok {
                return  // Channel closed, exit worker
            }
            // Create per-job timeout
            jobCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
            if err := wp.processor.Process(jobCtx, job); err != nil {
                log.Printf("worker %d: job failed: %v", id, err)
            }
            cancel()
            
        case <-ctx.Done():
            return  // Shutdown signal received
        }
    }
}

func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobQueue <- job:
        return nil
    default:
        return errors.New("queue full")
    }
}

func (wp *WorkerPool) Shutdown() {
    close(wp.jobQueue)
    wp.wg.Wait()
}
```

#### Worker Count Guidelines

| Workload Type | Recommended Workers | Rationale |
|--------------|---------------------|-----------|
| CPU-bound | `runtime.NumCPU()` | No benefit beyond core count |
| I/O-bound | `2-10× NumCPU()` | Blocked workers yield CPU |
| Mixed | `2-3× NumCPU()` | Start here, then benchmark |
| Database writes | Match connection pool | Prevent connection exhaustion |

---

### 3.2 Pipeline Pattern

Pipelines connect processing stages via channels, where each stage's output becomes the next stage's input.

#### Structure

```
┌─────────┐     ┌───────────┐     ┌─────────┐     ┌────────┐
│ Generate │────▶│ Transform │────▶│ Filter  │────▶│ Output │
└─────────┘     └───────────┘     └─────────┘     └────────┘
    Stage 1         Stage 2          Stage 3        Stage 4
```

#### Implementation

```go
// Stage 1: Generate values
func Generate(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Stage 2: Transform (square)
func Square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Stage 3: Filter (only even)
func FilterEven(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n%2 == 0 {
                select {
                case out <- n:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

// Usage: Construct pipeline
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Pipeline: Generate → Square → Filter
    pipeline := FilterEven(ctx, Square(ctx, Generate(ctx, 1, 2, 3, 4, 5)))

    for result := range pipeline {
        fmt.Println(result)  // 4, 16
    }
}
```

---

### 3.3 Semaphore Pattern

A semaphore limits the number of concurrent operations. In Go, use a buffered channel as a counting semaphore.

#### Implementation

```go
type Semaphore struct {
    sem chan struct{}
}

func NewSemaphore(maxConcurrency int) *Semaphore {
    return &Semaphore{
        sem: make(chan struct{}, maxConcurrency),
    }
}

func (s *Semaphore) Acquire() {
    s.sem <- struct{}{}
}

func (s *Semaphore) Release() {
    <-s.sem
}

func (s *Semaphore) TryAcquire() bool {
    select {
    case s.sem <- struct{}{}:
        return true
    default:
        return false
    }
}
```

#### Usage with HTTP Requests

```go
func FetchAllURLs(urls []string) []Response {
    sem := NewSemaphore(10)  // Max 10 concurrent requests
    var wg sync.WaitGroup
    results := make([]Response, len(urls))

    for i, url := range urls {
        wg.Add(1)
        go func(i int, url string) {
            defer wg.Done()
            
            sem.Acquire()
            defer sem.Release()
            
            results[i] = fetch(url)
        }(i, url)
    }

    wg.Wait()
    return results
}
```

#### Alternative: golang.org/x/sync/semaphore

```go
import "golang.org/x/sync/semaphore"

func main() {
    ctx := context.Background()
    sem := semaphore.NewWeighted(10)  // Max 10 concurrent

    for _, url := range urls {
        // Acquire with weight 1
        if err := sem.Acquire(ctx, 1); err != nil {
            log.Printf("failed to acquire semaphore: %v", err)
            continue
        }
        
        go func(url string) {
            defer sem.Release(1)
            fetch(url)
        }(url)
    }
}
```

---

### 3.4 Fan-Out/Fan-In Pattern

**Fan-Out**: Distribute work to multiple workers  
**Fan-In**: Collect results from multiple workers into a single channel

```go
// Fan-Out: Start multiple workers on same input channel
func FanOut(ctx context.Context, input <-chan Job, workers int) []<-chan Result {
    outputs := make([]<-chan Result, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(ctx, input)
    }
    return outputs
}

// Worker processes jobs and produces results
func worker(ctx context.Context, jobs <-chan Job) <-chan Result {
    out := make(chan Result)
    go func() {
        defer close(out)
        for job := range jobs {
            select {
            case out <- process(job):
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Fan-In: Merge multiple channels into one
func FanIn(ctx context.Context, channels ...<-chan Result) <-chan Result {
    out := make(chan Result)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan Result) {
            defer wg.Done()
            for result := range c {
                select {
                case out <- result:
                case <-ctx.Done():
                    return
                }
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

---

## Resilience Patterns

### 4.1 Backpressure

Backpressure prevents system overload by signaling upstream producers to slow down when downstream consumers can't keep up.

#### Strategies

| Strategy | Implementation | When to Use |
|----------|---------------|-------------|
| **Blocking** | Unbuffered channel | Synchronous processing |
| **Bounded Buffer** | Buffered channel | Smooth out bursts |
| **Drop** | `select` with `default` | Lossy acceptable |
| **Reject (503)** | HTTP 503 response | Client can retry elsewhere |

#### HTTP Handler with Backpressure

```go
func NewBackpressureHandler(queueSize int) http.Handler {
    queue := make(chan Request, queueSize)
    
    // Start consumer
    go func() {
        for req := range queue {
            processRequest(req)
        }
    }()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        req := parseRequest(r)
        
        select {
        case queue <- req:
            w.WriteHeader(http.StatusAccepted)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "accepted",
            })
        default:
            // Queue full - apply backpressure
            w.Header().Set("Retry-After", "5")
            http.Error(w, "Service overloaded, try again later", 
                       http.StatusServiceUnavailable)
        }
    })
}
```

---

### 4.2 Circuit Breaker

The Circuit Breaker pattern prevents cascade failures by stopping requests to a failing service.

#### States

```
          Success
    ┌──────────────┐
    │              │
    ▼              │
┌───────┐    ┌─────┴─────┐    ┌────────┐
│ Closed │───▶│ Half-Open │◀───│  Open  │
└───┬───┘    └───────────┘    └────┬───┘
    │                              │
    │      Failure Threshold       │
    └──────────────────────────────┘
              Timeout
```

#### Library: sony/gobreaker

```bash
go get github.com/sony/gobreaker/v2
```

```go
import "github.com/sony/gobreaker/v2"

func main() {
    settings := gobreaker.Settings{
        Name:        "HTTP Client",
        MaxRequests: 3,                    // Max requests in half-open
        Interval:    10 * time.Second,     // Interval to clear counts
        Timeout:     30 * time.Second,     // Time in open state
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures > 5
        },
        OnStateChange: func(name string, from, to gobreaker.State) {
            log.Printf("Circuit breaker %s: %s → %s", name, from, to)
        },
    }

    cb := gobreaker.NewCircuitBreaker[[]byte](settings)

    // Use circuit breaker
    body, err := cb.Execute(func() ([]byte, error) {
        resp, err := http.Get("https://api.example.com/data")
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        return io.ReadAll(resp.Body)
    })
}
```

#### Library: mercari/go-circuitbreaker

```go
import "github.com/mercari/go-circuitbreaker"

cb := circuitbreaker.New(
    circuitbreaker.WithFailOnContextCancel(true),
    circuitbreaker.WithOpenTimeout(10*time.Second),
    circuitbreaker.WithTripFunc(
        circuitbreaker.NewTripFuncFailureRate(10, 0.4),  // 40% failure rate
    ),
)

result, err := cb.Do(ctx, func() (interface{}, error) {
    return callExternalService()
})
```

---

### 4.3 Graceful Shutdown

Graceful shutdown ensures in-flight requests complete before termination.

#### Using signal.NotifyContext (Go 1.16+)

```go
func main() {
    // Create context that cancels on SIGINT/SIGTERM
    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt,
        syscall.SIGTERM,
    )
    defer stop()

    // Create server
    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
        BaseContext: func(l net.Listener) context.Context {
            return ctx  // Propagate cancellation to handlers
        },
    }

    // Start server in background
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for shutdown signal
    <-ctx.Done()
    log.Println("Shutdown signal received")

    // Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("Forced shutdown: %v", err)
    }

    log.Println("Server stopped gracefully")
}
```

#### Complete Lifecycle Management

```go
type Application struct {
    server   *http.Server
    workers  *WorkerPool
    db       *sql.DB
}

func (app *Application) Run(ctx context.Context) error {
    // Start all components
    app.workers.Start(ctx)
    
    errChan := make(chan error, 1)
    go func() {
        if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
            errChan <- err
        }
    }()

    // Wait for shutdown or error
    select {
    case err := <-errChan:
        return fmt.Errorf("server error: %w", err)
    case <-ctx.Done():
        return app.shutdown()
    }
}

func (app *Application) shutdown() error {
    var errs []error

    // 1. Stop accepting new HTTP requests
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    
    if err := app.server.Shutdown(shutdownCtx); err != nil {
        errs = append(errs, fmt.Errorf("http shutdown: %w", err))
    }

    // 2. Stop worker pool (drain queue)
    app.workers.Shutdown()

    // 3. Close database connections
    if err := app.db.Close(); err != nil {
        errs = append(errs, fmt.Errorf("db close: %w", err))
    }

    return errors.Join(errs...)
}
```

---

## Observability Tools

### 5.1 Structured Logging (slog)

Go 1.21 introduced `log/slog` for structured logging in the standard library.

#### Basic Usage

```go
import "log/slog"

func main() {
    // JSON handler for production
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Log with structured fields
    slog.Info("request processed",
        slog.String("method", "GET"),
        slog.String("path", "/api/users"),
        slog.Int("status", 200),
        slog.Duration("latency", 45*time.Millisecond),
    )
}
```

Output:
```json
{"time":"2025-11-28T09:00:00Z","level":"INFO","msg":"request processed","method":"GET","path":"/api/users","status":200,"latency":45000000}
```

#### Logger with Context

```go
func RequestHandler(w http.ResponseWriter, r *http.Request) {
    // Create logger with request context
    logger := slog.With(
        slog.String("request_id", r.Header.Get("X-Request-ID")),
        slog.String("user_agent", r.UserAgent()),
    )

    logger.Info("handling request",
        slog.String("method", r.Method),
        slog.String("path", r.URL.Path),
    )

    // Pass logger via context
    ctx := context.WithValue(r.Context(), "logger", logger)
    processRequest(ctx)
}
```

---

### 5.2 Prometheus Metrics

#### Installation

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

#### Metric Types

| Type | Use Case | Example |
|------|----------|---------|
| Counter | Cumulative values | Total requests, errors |
| Gauge | Current value | Active connections, queue depth |
| Histogram | Value distribution | Request latency |
| Summary | Percentiles | Response time percentiles |

#### Implementation

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )

    queueDepth = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "job_queue_depth",
            Help: "Current number of jobs in queue",
        },
    )
)

// Middleware to record metrics
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap ResponseWriter to capture status
        wrapped := &statusRecorder{ResponseWriter: w, status: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start).Seconds()
        
        requestsTotal.WithLabelValues(
            r.Method, 
            r.URL.Path, 
            strconv.Itoa(wrapped.status),
        ).Inc()
        
        requestDuration.WithLabelValues(
            r.Method, 
            r.URL.Path,
        ).Observe(duration)
    })
}

// Expose metrics endpoint
func main() {
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
```

---

### 5.3 pprof Profiling

Go's built-in profiler for CPU, memory, and concurrency analysis.

#### Enable in HTTP Server

```go
import _ "net/http/pprof"

func main() {
    // pprof endpoints automatically registered at /debug/pprof/
    go http.ListenAndServe(":6060", nil)
    
    // Your main application...
}
```

#### Profile Types

| Profile | Endpoint | Purpose |
|---------|----------|---------|
| CPU | `/debug/pprof/profile?seconds=30` | CPU usage over time |
| Heap | `/debug/pprof/heap` | Memory allocations |
| Goroutine | `/debug/pprof/goroutine` | All goroutine stack traces |
| Block | `/debug/pprof/block` | Blocking synchronization |
| Mutex | `/debug/pprof/mutex` | Mutex contention |

#### Command-Line Analysis

```bash
# CPU profile (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Interactive commands
(pprof) top10           # Top 10 functions
(pprof) web             # Open in browser
(pprof) list funcName   # Show source code
```

#### Benchmark Profiling

```bash
# Run benchmark with CPU profile
go test -bench=. -cpuprofile=cpu.prof

# Run benchmark with memory profile
go test -bench=. -memprofile=mem.prof

# Analyze
go tool pprof cpu.prof
```

---

## Advanced Patterns

### 6.1 errgroup Package

`errgroup` combines `sync.WaitGroup` with error propagation and context cancellation.

#### Installation

```bash
go get golang.org/x/sync/errgroup
```

#### Basic Usage

```go
import "golang.org/x/sync/errgroup"

func FetchAllData(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url  // Capture loop variables
        g.Go(func() error {
            data, err := fetch(ctx, url)
            if err != nil {
                return err  // First error cancels all other goroutines
            }
            results[i] = data
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

#### With Concurrency Limit

```go
func ProcessItems(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(10)  // Max 10 concurrent goroutines

    for _, item := range items {
        item := item
        g.Go(func() error {
            return processItem(ctx, item)
        })
    }

    return g.Wait()
}
```

---

### 6.2 Sharded Counters

For extremely high-throughput counters, sharding reduces contention.

#### Implementation

```go
type ShardedCounter struct {
    shards []uint64
    mask   uint64
}

func NewShardedCounter(numShards int) *ShardedCounter {
    // Round to power of 2 for fast modulo via bitmask
    n := 1
    for n < numShards {
        n <<= 1
    }
    return &ShardedCounter{
        shards: make([]uint64, n),
        mask:   uint64(n - 1),
    }
}

func (c *ShardedCounter) Inc() {
    // Use goroutine ID or random shard for distribution
    shard := fastrand() & c.mask
    atomic.AddUint64(&c.shards[shard], 1)
}

func (c *ShardedCounter) Value() uint64 {
    var total uint64
    for i := range c.shards {
        total += atomic.LoadUint64(&c.shards[i])
    }
    return total
}

// Fast random using runtime trick
//go:linkname fastrand runtime.fastrand
func fastrand() uint64
```

#### When to Use

| Counter Type | Use Case | Throughput |
|--------------|----------|------------|
| Single atomic | Low-moderate contention | ~10M ops/sec |
| Sharded (8 shards) | High contention | ~80M ops/sec |
| Sharded (64 shards) | Extreme contention | ~500M ops/sec |

---

## External Libraries Reference

| Library | Purpose | Installation |
|---------|---------|--------------|
| `golang.org/x/sync/errgroup` | Error-aware goroutine groups | `go get golang.org/x/sync` |
| `golang.org/x/sync/semaphore` | Weighted semaphore | `go get golang.org/x/sync` |
| `github.com/sony/gobreaker` | Circuit breaker | `go get github.com/sony/gobreaker/v2` |
| `github.com/mercari/go-circuitbreaker` | Context-aware circuit breaker | `go get github.com/mercari/go-circuitbreaker` |
| `github.com/prometheus/client_golang` | Prometheus metrics | `go get github.com/prometheus/client_golang` |
| `github.com/alitto/pond` | Worker pool library | `go get github.com/alitto/pond` |

---

## Quick Reference Cheat Sheet

### Context Patterns

```go
// Create with cancellation
ctx, cancel := context.WithCancel(parent)
defer cancel()

// Create with timeout
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// Check if cancelled
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // Continue processing
}
```

### Worker Pool Template

```go
workers := 50
queue := make(chan Job, 1000)
var wg sync.WaitGroup

for i := 0; i < workers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for job := range queue {
            process(job)
        }
    }()
}

// Submit jobs...
close(queue)
wg.Wait()
```

### Graceful Shutdown Template

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

// Start server...
<-ctx.Done()

shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
server.Shutdown(shutdownCtx)
```

### Metrics Template

```go
var counter = promauto.NewCounter(prometheus.CounterOpts{
    Name: "myapp_processed_total",
    Help: "Total processed items",
})

counter.Inc()
```

### Benchmark Command

```bash
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

---

## Further Reading

- [Go Concurrency Patterns](https://go.dev/blog/pipelines) - Official Go Blog
- [Context Package Documentation](https://pkg.go.dev/context) - pkg.go.dev
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency) - Official Guide
- [Go Memory Model](https://go.dev/ref/mem) - Understanding synchronization
