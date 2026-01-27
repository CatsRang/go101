# Go + Echo REST Development Guide with Context Management

A comprehensive guide for building production-ready REST APIs in Go using Echo framework with proper context propagation, graceful shutdown, and high-scale concurrency patterns.

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Server Lifecycle & Context Flow](#2-server-lifecycle--context-flow)
3. [Signal Handling & Graceful Shutdown](#3-signal-handling--graceful-shutdown)
4. [Context Propagation Patterns](#4-context-propagation-patterns)
5. [Middleware Stack Design](#5-middleware-stack-design)
6. [Worker Pool Integration](#6-worker-pool-integration)
7. [HTTP Handler Patterns](#7-http-handler-patterns)
8. [Configuration & Initialization](#8-configuration--initialization)
9. [Logging Architecture](#9-logging-architecture)
10. [Error Handling Patterns](#10-error-handling-patterns)
11. [Quick Reference](#11-quick-reference)

---

## 1. Architecture Overview

### Stack Components

| Component | Library | Purpose |
|-----------|---------|---------|
| Web Framework | `github.com/labstack/echo/v4` | HTTP routing, middleware |
| CLI | `github.com/spf13/cobra` | Command-line interface |
| Configuration | `github.com/spf13/viper` | Config file + env vars |
| Logging | `log/slog` (stdlib) | Structured JSON logging |
| Transaction Log | `gopkg.in/natefinch/lumberjack.v2` | Rotating file logs |

### Clean Architecture Flow

```
┌──────────────────────────────────────────────────────────────────────┐
│                           cmd/main.go                                │
│                    server.Execute() entry point                      │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│                        pkg/server/server.go                          │
│  • Cobra CLI initialization                                          │
│  • signal.NotifyContext() for root context                           │
│  • Echo server + middleware configuration                            │
│  • Graceful shutdown orchestration                                   │
└──────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│                        pkg/domain/<domain>/                          │
│  handler.go → service.go → repository.go → database                 │
└──────────────────────────────────────────────────────────────────────┘
```

### Context Propagation Tree

```
                    ┌─────────────────────────────┐
                    │   signal.NotifyContext()    │
                    │     (Root Context)          │
                    │  SIGINT/SIGTERM/SIGHUP      │
                    └─────────────┬───────────────┘
                                  │
                    ┌─────────────▼───────────────┐
                    │      Server.Setup(ctx)      │
                    │   Stores ctx in Server      │
                    └─────────────┬───────────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
     ┌────────▼────────┐ ┌────────▼────────┐ ┌───────▼───────┐
     │  Worker Pools   │ │  HTTP Handlers  │ │ Background    │
     │ InitPool(ctx)   │ │ (via middleware)│ │ Goroutines    │
     └────────┬────────┘ └────────┬────────┘ └───────┬───────┘
              │                   │                   │
              └───────────────────┴───────────────────┘
                                  │
                         All respond to ctx.Done()
```

---

## 2. Server Lifecycle & Context Flow

### Server Structure

```go
type Server struct {
    echo         *echo.Echo
    config       *config.AppConf
    db           *sql.DB
    ctx          context.Context    // Root context for propagation
    crawlService *crawl.Service     // Services needing shutdown
}
```

### Complete Lifecycle

```go
func (s *Server) Run() {
    // 1. Create root context with signal handling
    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt,
        syscall.SIGINT,
        syscall.SIGTERM,
        syscall.SIGHUP,
    )
    defer stop()

    // 2. Setup with context propagation
    s.Setup(ctx)

    // 3. Start server in goroutine
    serverErr := make(chan error, 1)
    go func() {
        if err := s.Start(); err != nil && err != http.ErrServerClosed {
            serverErr <- err
        }
    }()

    // 4. Wait for shutdown signal OR server error
    select {
    case <-ctx.Done():
        slog.Info("Shutdown signal received", "reason", ctx.Err())
    case err := <-serverErr:
        slog.Error("Server failed to start", "error", err)
        return
    }

    // 5. Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := s.Shutdown(shutdownCtx); err != nil {
        slog.Error("Shutdown error", "error", err)
    }
}
```

### Setup Phase

```go
func (s *Server) Setup(ctx context.Context) {
    // Store context for middleware propagation
    s.ctx = ctx

    // Initialize logging
    log.InitTranLogger(/* config params */)

    // Configure middleware (uses s.ctx)
    s.setupMiddleware()

    // Configure routes (passes ctx to worker pools)
    s.setupRoutes(ctx)
}
```

---

## 3. Signal Handling & Graceful Shutdown

### Signal Registration (Go 1.16+)

```go
// Modern approach using signal.NotifyContext
ctx, stop := signal.NotifyContext(context.Background(),
    os.Interrupt,      // Ctrl+C
    syscall.SIGINT,    // Interrupt
    syscall.SIGTERM,   // Termination (docker stop, k8s)
    syscall.SIGHUP,    // Hangup (terminal closed)
)
defer stop()

// Wait for signal
<-ctx.Done()
// ctx.Err() returns context.Canceled
```

### Shutdown Sequence (Order Matters)

```go
func (s *Server) Shutdown(ctx context.Context) error {
    slog.Info("Starting graceful shutdown sequence")
    var errs []error

    // 1. Stop accepting NEW HTTP requests first
    slog.Info("shutdown: stopping HTTP server")
    if err := s.echo.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("HTTP server shutdown: %w", err))
    }

    // 2. Drain worker pools (allow in-flight jobs to complete)
    slog.Info("shutdown: stopping worker pools")
    if s.crawlService != nil {
        if err := s.crawlService.Shutdown(ctx); err != nil {
            errs = append(errs, fmt.Errorf("worker pool shutdown: %w", err))
        }
    }

    // 3. Close database connections LAST
    slog.Info("shutdown: closing database connections")
    if err := db.CloseDB(s.db); err != nil {
        errs = append(errs, fmt.Errorf("database close: %w", err))
    }

    if len(errs) > 0 {
        return errors.Join(errs...)  // Go 1.20+ multi-error
    }

    slog.Info("Server shutdown completed successfully")
    return nil
}
```

### Shutdown Order Rationale

| Order | Component | Reason |
|-------|-----------|--------|
| 1 | HTTP Server | Stop accepting new requests |
| 2 | Worker Pools | Drain job queues, complete in-flight work |
| 3 | Background Tasks | Cancel long-running goroutines |
| 4 | Database | Workers may need DB during drain |

---

## 4. Context Propagation Patterns

### Middleware-Based Request Context

The context propagation middleware merges the server's root context with each request's context:

```go
func (s *Server) contextPropagationMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if s.ctx != nil {
                req := c.Request()
                originalCtx := req.Context()

                // Create merged context from server's root context
                mergedCtx, cancel := context.WithCancel(s.ctx)
                defer cancel()

                // Also respect original request cancellation (client disconnect)
                go func() {
                    select {
                    case <-originalCtx.Done():
                        cancel()  // Client disconnected
                    case <-mergedCtx.Done():
                        // Server shutting down
                    }
                }()

                c.SetRequest(req.WithContext(mergedCtx))
            }
            return next(c)
        }
    }
}
```

### Handler Context Usage

```go
func (h *Handler) CreateResource(c echo.Context) error {
    ctx := c.Request().Context()  // Gets merged context

    // Context is cancelled on:
    // 1. Client disconnect
    // 2. Server shutdown signal
    // 3. Request timeout (if set)

    result, err := h.service.Create(ctx, input)
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return echo.NewHTTPError(http.StatusServiceUnavailable, "request cancelled")
        }
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusCreated, result)
}
```

### Service Layer Context

```go
func (s *Service) Create(ctx context.Context, input CreateInput) (*Resource, error) {
    // Check context before expensive operations
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Pass context to repository
    return s.repo.Create(ctx, input)
}
```

### Repository Context (Database)

```go
func (r *Repository) Create(ctx context.Context, input CreateInput) (*Resource, error) {
    query := `INSERT INTO resources (name, value) VALUES (?, ?)`

    result, err := r.db.ExecContext(ctx, query, input.Name, input.Value)
    if err != nil {
        return nil, fmt.Errorf("insert resource: %w", err)
    }

    // ... handle result
}
```

---

## 5. Middleware Stack Design

### Middleware Order (Critical)

```go
func (s *Server) setupMiddleware() {
    // 1. Recovery - MUST be first (catches panics from all subsequent middleware)
    s.echo.Use(middleware.Recover())

    // 2. Context Propagation - Inject root context early
    s.echo.Use(s.contextPropagationMiddleware())

    // 3. CORS - Handle preflight before other processing
    s.applyCORSMiddleware()

    // 4. Request ID - Generate/extract for tracing
    s.echo.Use(s.requestIDMiddleware())

    // 5. Transaction Logging - Log request details
    s.echo.Use(s.transactionLoggingMiddleware())

    // 6. Auth middleware applied per-group in setupRoutes()
}
```

### Request ID Middleware

```go
func (s *Server) requestIDMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            reqID := c.Request().Header.Get("X-Request-ID")
            if reqID == "" {
                reqID = uuid.New().String()
            }

            c.Set("reqID", reqID)
            c.Set("startTime", time.Now())
            c.Response().Header().Set("X-Request-ID", reqID)

            return next(c)
        }
    }
}
```

### CORS Configuration

```go
func (s *Server) applyCORSMiddleware() {
    cfg := s.config.Server.CORS
    if !cfg.Enabled {
        slog.Warn("CORS disabled via configuration")
        return
    }

    corsConfig := middleware.CORSConfig{
        AllowOrigins:     cfg.AllowOrigins,
        AllowMethods:     cfg.AllowMethods,
        AllowHeaders:     cfg.AllowHeaders,
        AllowCredentials: cfg.AllowCredentials,
        ExposeHeaders:    cfg.ExposeHeaders,
        MaxAge:           cfg.MaxAge,
    }

    // Safe defaults
    if len(corsConfig.AllowOrigins) == 0 {
        corsConfig.AllowOrigins = []string{"*"}
    }
    if len(corsConfig.AllowMethods) == 0 {
        corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    }
    if len(corsConfig.AllowHeaders) == 0 {
        corsConfig.AllowHeaders = []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAccept,
            echo.HeaderAuthorization,
        }
    }

    s.echo.Use(middleware.CORSWithConfig(corsConfig))
}
```

---

## 6. Worker Pool Integration

### Pool Initialization with Context

```go
func (s *Server) setupRoutes(ctx context.Context) {
    // ... repository and service initialization ...

    // Initialize worker pool manager WITH context
    if err := crawlService.InitializePoolManager(ctx, s.db, s.config.CrawlWorker); err != nil {
        slog.Error("Failed to initialize worker pools", "error", err)
        panic(fmt.Sprintf("Failed to initialize worker pools: %v", err))
    }

    // Store service reference for shutdown
    s.crawlService = crawlService
}
```

### Worker Pool Pattern

```go
type WorkerPool struct {
    workers   int
    jobQueue  chan Job
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

func NewWorkerPool(ctx context.Context, workers, queueSize int) *WorkerPool {
    poolCtx, cancel := context.WithCancel(ctx)

    return &WorkerPool{
        workers:  workers,
        jobQueue: make(chan Job, queueSize),
        ctx:      poolCtx,
        cancel:   cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()

    for {
        select {
        case job, ok := <-wp.jobQueue:
            if !ok {
                return  // Channel closed
            }
            // Create per-job timeout
            jobCtx, cancel := context.WithTimeout(wp.ctx, 30*time.Second)
            if err := wp.processJob(jobCtx, job); err != nil {
                slog.Error("job failed", "worker", id, "error", err)
            }
            cancel()

        case <-wp.ctx.Done():
            return  // Shutdown signal
        }
    }
}

func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobQueue <- job:
        return nil
    case <-wp.ctx.Done():
        return errors.New("pool shutting down")
    default:
        return errors.New("queue full")
    }
}

func (wp *WorkerPool) Shutdown(ctx context.Context) error {
    // Signal workers to stop
    wp.cancel()

    // Close job queue to unblock workers waiting on channel
    close(wp.jobQueue)

    // Wait for workers with timeout
    done := make(chan struct{})
    go func() {
        wp.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Worker Count Guidelines

| Workload Type | Workers | Rationale |
|---------------|---------|-----------|
| CPU-bound | `runtime.NumCPU()` | No benefit beyond cores |
| I/O-bound (HTTP) | `2-10× NumCPU()` | Blocked workers yield CPU |
| Database writes | Match connection pool | Prevent connection exhaustion |
| Mixed | `2-3× NumCPU()` | Benchmark to tune |

---

## 7. HTTP Handler Patterns

### Route Registration Pattern

```go
// handler.go
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
    g.GET("/resources", h.List)
    g.POST("/resources", h.Create)
    g.GET("/resources/:id", h.Get)
    g.PUT("/resources/:id", h.Update)
    g.DELETE("/resources/:id", h.Delete)
}
```

### Handler Implementation

```go
func (h *Handler) Create(c echo.Context) error {
    ctx := c.Request().Context()

    var input CreateInput
    if err := c.Bind(&input); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
    }

    if err := c.Validate(input); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    result, err := h.service.Create(ctx, input)
    if err != nil {
        return h.handleError(err)
    }

    return c.JSON(http.StatusCreated, result)
}

func (h *Handler) handleError(err error) error {
    switch {
    case errors.Is(err, context.Canceled):
        return echo.NewHTTPError(http.StatusServiceUnavailable, "request cancelled")
    case errors.Is(err, context.DeadlineExceeded):
        return echo.NewHTTPError(http.StatusGatewayTimeout, "request timeout")
    case errors.Is(err, ErrNotFound):
        return echo.NewHTTPError(http.StatusNotFound, err.Error())
    case errors.Is(err, ErrConflict):
        return echo.NewHTTPError(http.StatusConflict, err.Error())
    default:
        slog.Error("handler error", "error", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
    }
}
```

### Server Routes Setup

```go
func (s *Server) setupRoutes(ctx context.Context) {
    // Initialize dependency chain
    repo := NewRepository(s.db)
    service := NewService(repo)
    handler := NewHandler(service)

    // Public routes (no auth)
    authGroup := s.echo.Group("/api/v1/auth")
    authHandler.RegisterRoutes(authGroup)

    // Protected routes
    apiGroup := s.echo.Group("/api/v1")
    if s.config.Auth.Enabled {
        apiGroup.Use(auth.JWTMiddleware(&s.config.JWT))
    } else {
        apiGroup.Use(s.testBypassAuthMiddleware())
    }

    // Register domain handlers
    handler.RegisterRoutes(apiGroup)
}
```

---

## 8. Configuration & Initialization

### Configuration Singleton Pattern

```go
// pkg/config/app_conf.go
var (
    sharedConf *AppConf
    confOnce   sync.Once
)

func SharedAppConf() *AppConf {
    confOnce.Do(func() {
        sharedConf = &AppConf{}
    })
    return sharedConf
}

func (c *AppConf) Init(configPath string) error {
    viper.SetConfigFile(configPath)
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("read config: %w", err)
    }

    if err := viper.Unmarshal(c); err != nil {
        return fmt.Errorf("unmarshal config: %w", err)
    }

    return nil
}
```

### Cobra CLI Integration

```go
var (
    confPath string
    rootCmd  = &cobra.Command{
        Use:     "app",
        Version: config.APP_VERSION,
        Short:   config.APP_DESC_SHORT,
        Run:     rootRun,
    }
)

func Execute() {
    rootCmd.Flags().StringVarP(&confPath, "conf", "c", "config.yml", "Config file path")
    cobra.OnInitialize(onInitializeRoot)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func onInitializeRoot() {
    appConf := config.SharedAppConf()

    if err := appConf.Init(confPath); err != nil {
        fmt.Printf("Failed to initialize config: %v\n", err)
        os.Exit(1)
    }

    log.InitLogger(appConf.Log.LogLevel)

    slog.Info("Application started",
        "app_id", appConf.AppId,
        "version", appConf.AppVersion,
        "config", confPath)
}

func rootRun(cmd *cobra.Command, args []string) {
    appConf := config.SharedAppConf()

    db, err := db.InitDB(&appConf.Database)
    if err != nil {
        slog.Error("Failed to initialize database", "error", err)
        os.Exit(1)
    }

    srv := New(appConf, db)
    srv.Run()  // Blocks until shutdown
}
```

---

## 9. Logging Architecture

### Application Logger (slog)

```go
// pkg/util/log/logger.go
func InitLogger(level string) {
    var lvl slog.Level
    switch strings.ToLower(level) {
    case "debug":
        lvl = slog.LevelDebug
    case "warn", "warning":
        lvl = slog.LevelWarn
    case "error":
        lvl = slog.LevelError
    default:
        lvl = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{Level: lvl}
    handler := slog.NewJSONHandler(os.Stdout, opts)
    slog.SetDefault(slog.New(handler))
}
```

### Transaction Logger (Lumberjack)

```go
// pkg/util/log/tranlog.go
type TranLogger struct {
    logger *slog.Logger
}

func InitTranLogger(folderPath string, maxSizeMB, maxAgeDays, maxBackups int) *TranLogger {
    writer := &lumberjack.Logger{
        Filename:   filepath.Join(folderPath, "transaction.log"),
        MaxSize:    maxSizeMB,
        MaxAge:     maxAgeDays,
        MaxBackups: maxBackups,
        Compress:   true,
    }

    handler := slog.NewJSONHandler(writer, nil)
    logger := slog.New(handler)

    return &TranLogger{logger: logger}
}

func (t *TranLogger) LogRequestStart(reqID, method, path, clientIP, userAgent string, contentLen int64) {
    t.logger.Info("request_start",
        slog.String("req_id", reqID),
        slog.String("method", method),
        slog.String("path", path),
        slog.String("client_ip", clientIP),
        slog.String("user_agent", userAgent),
        slog.Int64("content_length", contentLen),
    )
}
```

### Logging Best Practices

```go
// Use structured fields, not string concatenation
slog.Info("request processed",
    slog.String("method", "GET"),
    slog.String("path", "/api/users"),
    slog.Int("status", 200),
    slog.Duration("latency", 45*time.Millisecond),
)

// Include context for request correlation
func logWithContext(ctx context.Context, msg string, args ...any) {
    reqID := ctx.Value("reqID")
    if reqID != nil {
        args = append(args, slog.String("req_id", reqID.(string)))
    }
    slog.Info(msg, args...)
}

// Error logging with stack context
slog.Error("operation failed",
    slog.String("operation", "create_resource"),
    slog.String("resource_id", id),
    slog.Any("error", err),
)
```

---

## 10. Error Handling Patterns

### Domain Errors

```go
// pkg/domain/<domain>/errors.go
var (
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
)

// Wrap with context
func (s *Service) Get(ctx context.Context, id string) (*Resource, error) {
    resource, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("%w: id=%s", ErrNotFound, id)
        }
        return nil, fmt.Errorf("get resource: %w", err)
    }
    return resource, nil
}
```

### Context Error Handling

```go
func handleContextError(err error) error {
    switch {
    case errors.Is(err, context.Canceled):
        // Client disconnected OR shutdown signal
        return echo.NewHTTPError(http.StatusServiceUnavailable, "request cancelled")
    case errors.Is(err, context.DeadlineExceeded):
        // Timeout
        return echo.NewHTTPError(http.StatusGatewayTimeout, "request timeout")
    default:
        return err
    }
}
```

### Multi-Error Aggregation (Go 1.20+)

```go
func (s *Server) Shutdown(ctx context.Context) error {
    var errs []error

    if err := s.echo.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("HTTP: %w", err))
    }
    if err := s.workerPool.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("workers: %w", err))
    }
    if err := s.db.Close(); err != nil {
        errs = append(errs, fmt.Errorf("database: %w", err))
    }

    return errors.Join(errs...)  // Returns nil if slice is empty
}
```

---

## 11. Quick Reference

### Server Lifecycle Template

```go
func (s *Server) Run() {
    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer stop()

    s.Setup(ctx)

    serverErr := make(chan error, 1)
    go func() {
        if err := s.Start(); err != http.ErrServerClosed {
            serverErr <- err
        }
    }()

    select {
    case <-ctx.Done():
    case err := <-serverErr:
        slog.Error("server error", "error", err)
        return
    }

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    s.Shutdown(shutdownCtx)
}
```

### Context Check Pattern

```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // Continue processing
}
```

### Worker Pool Template

```go
for i := 0; i < workers; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case job := <-queue:
                process(job)
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

### Middleware Order

1. Recovery (panic → 500)
2. Context Propagation (root context)
3. CORS (preflight)
4. Request ID (tracing)
5. Transaction Logging (audit)
6. Authentication (per-group)

### Shutdown Order

1. HTTP Server (stop new requests)
2. Worker Pools (drain queues)
3. Background Tasks (cancel)
4. Database (close last)

---

## Further Reading

- [Go Context Package](https://pkg.go.dev/context)
- [Echo Framework Documentation](https://echo.labstack.com/docs)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Graceful Shutdown in Go](https://pkg.go.dev/os/signal#NotifyContext)
