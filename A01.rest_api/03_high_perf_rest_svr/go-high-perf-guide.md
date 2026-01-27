# High Performance Go Server Guide

This guide details the architectural patterns and configuration tuning required to build a production-ready, high-performance Go HTTP server. It focuses on resilience (graceful shutdown, cancellation) and throughput (concurrency, resource reuse).

## 1. Executive Summary: The Production Checklist

| Component | Key Optimization |
| :--- | :--- |
| **Graceful Shutdown** | Use `signal.NotifyContext` + `server.Shutdown(ctx)` + Load Balancer drain sleep. |
| **Concurrency** | Bounded concurrency via **Worker Pools** or **Semaphores** (never unbounded `go func()`). |
| **Network Tuning** | Set explicit `ReadTimeout`, `WriteTimeout`, and `ReadHeaderTimeout`. |
| **Database** | Limit `MaxOpenConns` to prevent DB saturation; tune `MaxIdleConns` for reuse. |
| **Memory** | Use `sync.Pool` for hot-path object reuse (e.g., byte buffers, JSON encoders). |
| **Serialization** | Replace `encoding/json` with `goccy/go-json` or `sonic` for 2-5x perf gain. |

---

## 2. Architecture: Graceful Shutdown & Cancellation

A high-performance server must handle termination without dropping active requests.

### The Shutdown Lifecycle
1.  **Listen for Signals:** Catch `SIGINT` and `SIGTERM`.
2.  **Drain Traffic (K8s/LB):** Sleep briefly (e.g., 5-10s) to allow Load Balancers to remove the pod from rotation.
3.  **Stop Accepting:** Call `server.Shutdown()`.
4.  **Wait/Cancel:** Wait for active requests to finish or hit a hard timeout.

### Implementation Pattern

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Create a context that listens for OS signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: myRouter(),
		// 2. Production Timeouts (Critical)
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		// Prevent Slowloris attacks by setting Header timeout
		ReadHeaderTimeout: 2 * time.Second,
	}

	// 3. Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 4. Block until signal is received
	<-ctx.Done()
	log.Println("Shutdown signal received...")

	// 5. (Optional) Readiness Drain: Sleep to let LB remove pod
	// time.Sleep(5 * time.Second)

	// 6. Graceful Shutdown with hard timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func myRouter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
```

### Understanding `signal.NotifyContext()` Internals

`signal.NotifyContext()` (Go 1.16+) is syntactic sugar that combines signal handling with context cancellation. Understanding its internals helps debug shutdown issues.

**What it does under the hood:**

```go
// Simplified internal implementation of signal.NotifyContext
func NotifyContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(parent)
    c := make(chan os.Signal, 1)
    signal.Notify(c, signals...)  // Register channel to receive signals

    go func() {  // <-- Hidden goroutine spawned here
        select {
        case <-c:
            cancel()  // Cancel context when signal arrives
        case <-ctx.Done():
            // Parent cancelled or stop() called
        }
        signal.Stop(c)  // Unregister signal notification
    }()

    return ctx, cancel
}
```

**Key insights:**

| Aspect | Detail |
| :--- | :--- |
| **Hidden goroutine** | A background goroutine waits for signals - this IS the signal handler |
| **`stop()` function** | Calls both `signal.Stop()` (unregister) AND `cancel()` (cancel context) |
| **Channel buffer** | Uses buffered channel (size 1) to prevent signal loss |
| **Cleanup** | `defer stop()` is critical - prevents goroutine leak and cleans up signal registration |

**Comparison with pre-Go 1.16 pattern:**

```go
// Old way (explicit, verbose)
ctx, cancel := context.WithCancel(context.Background())
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
go func() {
    <-sigCh
    cancel()
}()
defer signal.Stop(sigCh)
defer cancel()

// New way (Go 1.16+, recommended)
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()
```

### Context Propagation
Never use `context.Background()` inside a request handler. Always propagate the request context to ensure that if a user disconnects, the database query or internal work is cancelled immediately.

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // This context is cancelled if the client disconnects or server shuts down
    ctx := r.Context() 

    // Pass ctx to DB, Redis, or API calls
    // user, err := db.GetUser(ctx, userID) 
}
```

---

## 3. High Concurrency Patterns

Unbounded concurrency (spawning a goroutine per task) is the #1 cause of Go server crashes under load.

### Pattern A: The Worker Pool (Processing Pipeline)
Use this for CPU-intensive tasks or to limit load on downstream systems (e.g., rate-limited 3rd party APIs).

```go
type Job struct {
    Payload string
}

func StartDispatcher(workerCount int) chan Job {
    jobQueue := make(chan Job, 100) // Buffered channel

    for i := 0; i < workerCount; i++ {
        go func(id int) {
            for job := range jobQueue {
                process(job) // Blocking processing
            }
        }(i)
    }
    return jobQueue
}

// Usage in Handler
// jobs <- Job{Payload: "data"} // Non-blocking if buffer isn't full
```

### Pattern B: Semaphore (Bounded Concurrency)
Use `golang.org/x/sync/semaphore` when you don't need a dedicated pool but want to limit concurrent heavy operations (e.g., "only allow 10 concurrent image resizes").

```go
var sem = semaphore.NewWeighted(10) // Max 10 concurrent ops

func heavyHandler(w http.ResponseWriter, r *http.Request) {
    // Acquire 1 slot
    if err := sem.Acquire(r.Context(), 1); err != nil {
        http.Error(w, "Server busy", http.StatusTooManyRequests)
        return
    }
    defer sem.Release(1)

    // Do heavy work...
}
```

---

## 4. Performance Tuning & Configuration

### 4.1 HTTP Server Timeouts
Default Go `http.Server` has **no timeouts**, meaning a malicious client can hold a connection open forever (Slowloris attack).

- **ReadHeaderTimeout**: **Must be set**. Controls how long to wait for headers.
- **ReadTimeout**: Max time to read the entire request (including body).
- **WriteTimeout**: Max time to write the response.

### 4.2 Database Connection Pooling
Misconfigured DB pools cause high latency or connection resets.

- **SetMaxOpenConns(N)**: **Required**. Set to `(CPU Cores * 2)` or limit based on DB capacity. Prevents opening 10k connections during a spike.
- **SetMaxIdleConns(N)**: Set equal to `MaxOpenConns` if memory allows. Tearing down and rebuilding connections is expensive; keep them warm.
- **SetConnMaxLifetime(duration)**: Set to 5-10 minutes (shorter than load balancer/firewall idle timeouts) to prevent using stale connections that get reset by the network.

### 4.3 Memory Optimization (Hot Paths)

**Sync.Pool**
Use `sync.Pool` to reuse heavy structs or byte buffers to reduce GC pressure.

```go
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func getBuffer() *bytes.Buffer {
    b := bufPool.Get().(*bytes.Buffer)
    b.Reset()
    return b
}
```

**JSON Serialization**
Standard `encoding/json` uses reflection and is slow. For high-throughput endpoints, switch to drop-in replacements:
- **goccy/go-json**: [GitHub](https://github.com/goccy/go-json)
- **sonic** (Bytedance): [GitHub](https://github.com/bytedance/sonic)

---

## 5. Common Pitfalls to Avoid

1.  **Ignoring Request Body Close:** Always `defer r.Body.Close()` in HTTP clients.
2.  **Goroutine Leaks:** Ensure every `go func` has a way to exit (e.g., listens to `ctx.Done()`).
3.  **Log Contention:** Use a high-performance structured logger like **[zerolog](https://github.com/rs/zerolog)** or **[zap](https://github.com/uber-go/zap)** instead of `log.Printf`, which locks on every write.

---

## 6. References & Further Reading

- **Graceful Shutdown Patterns**: [VictoriaMetrics Blog](https://victoriametrics.com/blog/go-graceful-shutdown/)
- **Go Context & Cancellation**: [Go Blog: Context](https://go.dev/blog/context)
- **Database Connection Management**: [Go Docs: Managing Connections](https://go.dev/doc/database/manage-connections)
- **HTTP Timeouts Guide**: [Cloudflare Blog](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- **Worker Pool Implementation**: [Go Patterns: Worker Pool](https://go-patterns.dev/parallel-computing/worker-pool)
