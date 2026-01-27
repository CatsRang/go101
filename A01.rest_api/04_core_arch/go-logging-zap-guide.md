# Comprehensive Go Zap Logger Guide

An in-depth guide to mastering Uber's Zap logger - the blazing fast, structured logging library for Go applications. This document covers everything from basic usage to advanced production patterns.

---

## Table of Contents

1. [Introduction](#introduction)
2. [Installation & Setup](#installation--setup)
3. [Core Concepts](#core-concepts)
4. [Logger vs SugaredLogger](#logger-vs-sugaredlogger)
5. [Log Levels](#log-levels)
6. [Configuration](#configuration)
7. [Structured Logging](#structured-logging)
8. [Contextual Logging](#contextual-logging)
9. [Custom Encoders](#custom-encoders)
10. [Output Configuration](#output-configuration)
11. [Log Rotation](#log-rotation)
12. [Sampling](#sampling)
13. [Error Handling](#error-handling)
14. [Production Patterns](#production-patterns)
15. [Testing](#testing)
16. [Performance Optimization](#performance-optimization)
17. [Best Practices](#best-practices)

---

## Introduction

### What is Zap?

Zap is a high-performance, structured logging library developed by Uber for Go applications. According to benchmarks, it's 4-10x faster than other structured logging packages and significantly faster than the standard library.

### Key Features

- **Blazing Fast**: Zero-allocation design for hot paths
- **Structured Logging**: JSON and console output formats
- **Type Safety**: Strongly-typed fields prevent errors
- **Leveled Logging**: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL
- **Flexible Configuration**: Production and development presets
- **Contextual Fields**: Add persistent context to loggers
- **Sampling**: Reduce log volume in high-frequency scenarios
- **Atomic Level Changes**: Adjust log levels without restarts

### Performance Comparison

| Package | Time | Time % to zap | Objects Allocated |
|---------|------|---------------|-------------------|
| **zap** | 193 ns/op | +0% | 0 allocs/op |
| **zap (sugared)** | 227 ns/op | +18% | 1 allocs/op |
| zerolog | 81 ns/op | -58% | 0 allocs/op |
| slog | 322 ns/op | +67% | 0 allocs/op |
| go-kit | 5377 ns/op | +2686% | 56 allocs/op |
| logrus | 21997 ns/op | +11297% | 68 allocs/op |

---

## Installation & Setup

### Installation

```bash
go get -u go.uber.org/zap
```

### Basic Usage

```go
package main

import (
    "go.uber.org/zap"
)

func main() {
    // Create a production logger
    logger, err := zap.NewProduction()
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // Flushes buffer, if any
    
    // Log a message
    logger.Info("Hello from Zap logger!")
}
```

**Output:**
```json
{"level":"info","ts":1684092708.7246346,"caller":"main.go:12","msg":"Hello from Zap logger!"}
```

### Import Statements

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)
```

---

## Core Concepts

### Logger Lifecycle

1. **Creation**: Initialize logger with `NewProduction()` or `NewDevelopment()`
2. **Usage**: Log messages with appropriate level methods
3. **Cleanup**: Always `defer logger.Sync()` to flush buffered logs

### Production vs Development

**Production Logger:**
```go
logger, _ := zap.NewProduction()
// - JSON format
// - Info level and above
// - Logs to stderr
```

**Output:**
```json
{"level":"info","ts":1684092708.7246346,"caller":"main.go:12","msg":"Hello"}
```

**Development Logger:**
```go
logger, _ := zap.NewDevelopment()
// - Console format
// - Debug level and above
// - Logs to stdout
// - Stacktraces on warnings
```

**Output:**
```
2023-05-14T20:42:39.137+0100 INFO main.go:12 Hello
```

### Environment-Based Configuration

```go
logger := zap.Must(zap.NewProduction())
if os.Getenv("APP_ENV") == "development" {
    logger = zap.Must(zap.NewDevelopment())
}
defer logger.Sync()
```

---

## Logger vs SugaredLogger

### Logger (Strongly Typed)

The `Logger` type is designed for maximum performance with zero allocations.

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("User logged in",
    zap.String("username", "johndoe"),
    zap.Int("userid", 123456),
    zap.String("provider", "google"),
    zap.Duration("latency", time.Millisecond*150),
)
```

**Output:**
```json
{"level":"info","ts":1684094903.7353888,"caller":"main.go:17","msg":"User logged in","username":"johndoe","userid":123456,"provider":"google","latency":0.15}
```

### SugaredLogger (Flexible)

The `SugaredLogger` provides a more ergonomic API at a small performance cost (10-50% slower).

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

sugar := logger.Sugar()

// Multiple APIs available
sugar.Info("Hello from Zap logger!")
sugar.Infoln("Hello from Zap logger!")
sugar.Infof("Hello from Zap logger! The time is %s", time.Now().Format("03:04 AM"))

// Loosely-typed key-value pairs
sugar.Infow("User logged in",
    "username", "johndoe",
    "userid", 123456,
    "provider", "google",
)
```

### Converting Between Types

```go
// Logger -> SugaredLogger
sugar := logger.Sugar()

// SugaredLogger -> Logger
logger := sugar.Desugar()
```

### When to Use Which

- **Logger**: Use in performance-critical paths, hot loops, high-frequency logging
- **SugaredLogger**: Use for general application logging where ergonomics matter

---

## Log Levels

Zap provides seven log levels in increasing order of severity:

### Log Level Reference

| Level | Integer | Purpose | Behavior |
|-------|---------|---------|----------|
| **DEBUG** | -1 | Detailed diagnostic information | Normal logging |
| **INFO** | 0 | General informational messages | Normal logging |
| **WARN** | 1 | Warning messages for unusual situations | Normal logging |
| **ERROR** | 2 | Error events that might still allow the app to continue | Includes stacktrace |
| **DPANIC** | 3 | Development panic (panics in dev, error in prod) | Panic in dev mode |
| **PANIC** | 4 | Severe error, then panics | Calls `panic()` |
| **FATAL** | 5 | Critical error, then exits | Calls `os.Exit(1)` |

### Usage Examples

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Debug("This is a debug message")
logger.Info("Application started successfully")
logger.Warn("High memory usage detected", zap.Int("percent", 85))
logger.Error("Failed to connect to database", zap.Error(err))
logger.DPanic("This should never happen") // Panics in development
logger.Panic("Critical system failure")    // Always panics
logger.Fatal("Cannot start application")   // Exits with status 1
```

### Setting Log Level

```go
config := zap.NewProductionConfig()
config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
logger, _ := config.Build()
```

### Dynamic Level Changes

```go
atom := zap.NewAtomicLevel()
atom.SetLevel(zap.WarnLevel) // Change level at runtime

config := zap.NewProductionConfig()
config.Level = atom
logger, _ := config.Build()

// Later, change level without restarting
atom.SetLevel(zap.DebugLevel)
```

---

## Configuration

### Using zap.Config

```go
func createLogger() *zap.Logger {
    encoderCfg := zap.NewProductionEncoderConfig()
    encoderCfg.TimeKey = "timestamp"
    encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
    
    config := zap.Config{
        Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
        Development:       false,
        DisableCaller:     false,
        DisableStacktrace: false,
        Sampling:          nil,
        Encoding:          "json",
        EncoderConfig:     encoderCfg,
        OutputPaths:       []string{"stdout", "logs/app.log"},
        ErrorOutputPaths:  []string{"stderr"},
        InitialFields: map[string]interface{}{
            "pid":     os.Getpid(),
            "version": "1.0.0",
        },
    }
    
    return zap.Must(config.Build())
}
```

**Output:**
```json
{"level":"info","timestamp":"2023-05-15T12:40:16.647+0100","caller":"main.go:42","msg":"Hello from Zap!","pid":2329946,"version":"1.0.0"}
```

### Encoder Configuration

```go
encoderConfig := zapcore.EncoderConfig{
    TimeKey:        "timestamp",
    LevelKey:       "level",
    NameKey:        "logger",
    CallerKey:      "caller",
    MessageKey:     "message",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.CapitalLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.StringDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
}

config := zap.Config{
    Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
    Development:   true,
    Encoding:      "console",
    EncoderConfig: encoderConfig,
    OutputPaths:   []string{"stdout"},
}

logger, _ := config.Build()
```

### Advanced Configuration with zapcore

```go
func createLogger() *zap.Logger {
    stdout := zapcore.AddSync(os.Stdout)
    file := zapcore.AddSync(&lumberjack.Logger{
        Filename:   "logs/app.log",
        MaxSize:    10, // megabytes
        MaxBackups: 3,
        MaxAge:     7, // days
    })
    
    level := zap.NewAtomicLevelAt(zap.InfoLevel)
    
    productionCfg := zap.NewProductionEncoderConfig()
    productionCfg.TimeKey = "timestamp"
    productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
    
    developmentCfg := zap.NewDevelopmentEncoderConfig()
    developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
    
    consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
    fileEncoder := zapcore.NewJSONEncoder(productionCfg)
    
    core := zapcore.NewTee(
        zapcore.NewCore(consoleEncoder, stdout, level),
        zapcore.NewCore(fileEncoder, file, level),
    )
    
    return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
```

---

## Structured Logging

### Adding Fields

Zap provides strongly-typed field constructors:

```go
logger.Info("User registered",
    zap.String("username", "alice"),
    zap.Int("age", 25),
    zap.Bool("verified", true),
    zap.Float64("balance", 99.99),
    zap.Duration("elapsed", time.Second*2),
    zap.Time("created_at", time.Now()),
    zap.Strings("roles", []string{"user", "admin"}),
)
```

### Common Field Types

```go
// Strings and primitives
zap.String("key", "value")
zap.Int("count", 42)
zap.Int64("id", 1234567890)
zap.Float64("price", 19.99)
zap.Bool("active", true)

// Time and duration
zap.Time("timestamp", time.Now())
zap.Duration("latency", time.Millisecond*150)

// Complex types
zap.Error(err)
zap.Any("custom", complexStruct)
zap.Reflect("reflected", anyValue)

// Arrays
zap.Strings("tags", []string{"go", "logging"})
zap.Ints("scores", []int{10, 20, 30})

// Binary
zap.Binary("data", []byte{0x01, 0x02})

// Namespaced
zap.Namespace("user")
zap.String("id", "123")
zap.String("name", "Alice")
```

### Nested Objects

```go
logger.Info("Request completed",
    zap.Namespace("request"),
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Int("status", 200),
    
    zap.Namespace("user"),
    zap.String("id", "user-123"),
    zap.String("ip", "192.168.1.1"),
)
```

**Output:**
```json
{
  "level": "info",
  "msg": "Request completed",
  "request": {
    "method": "GET",
    "path": "/api/users",
    "status": 200
  },
  "user": {
    "id": "user-123",
    "ip": "192.168.1.1"
  }
}
```

---

## Contextual Logging

### Using With()

Create child loggers with persistent fields:

```go
logger, _ := zap.NewProduction()
defer logger.Sync()

// Create child logger with context
requestLogger := logger.With(
    zap.String("request_id", "abc123"),
    zap.String("user_id", "42"),
)

// All logs from requestLogger include these fields
requestLogger.Info("Processing request")
requestLogger.Warn("Request took too long")
requestLogger.Error("Failed to process request", zap.Error(err))
```

**Output:**
```json
{"level":"info","ts":1721169926.2557518,"msg":"Processing request","request_id":"abc123","user_id":"42"}
{"level":"warn","ts":1721169926.255857,"msg":"Request took too long","request_id":"abc123","user_id":"42"}
{"level":"error","ts":1721169926.2558908,"msg":"Failed to process request","request_id":"abc123","user_id":"42","error":"timeout"}
```

### Context Pattern in HTTP Handlers

```go
type contextKey string

const loggerKey contextKey = "logger"

// Middleware to add logger to context
func requestLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger := zap.L().With(
            zap.String("request_id", uuid.New().String()),
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
        )
        
        ctx := context.WithValue(r.Context(), loggerKey, logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Retrieve logger from context
func loggerFromContext(ctx context.Context) *zap.Logger {
    if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
        return logger
    }
    return zap.L() // fallback to global logger
}

// Use in handler
func myHandler(w http.ResponseWriter, r *http.Request) {
    logger := loggerFromContext(r.Context())
    logger.Info("Handler called")
}
```

### Global Context Fields

```go
func createLogger() *zap.Logger {
    buildInfo, _ := debug.ReadBuildInfo()
    
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        os.Stdout,
        zap.InfoLevel,
    )
    
    return zap.New(core.With([]zapcore.Field{
        zap.String("go_version", buildInfo.GoVersion),
        zap.Int("pid", os.Getpid()),
        zap.String("host", hostname),
    }))
}
```

---

## Custom Encoders

### Creating a Custom Encoder

```go
type SensitiveFieldEncoder struct {
    zapcore.Encoder
    cfg zapcore.EncoderConfig
}

func (e *SensitiveFieldEncoder) EncodeEntry(
    entry zapcore.Entry,
    fields []zapcore.Field,
) (*buffer.Buffer, error) {
    filtered := make([]zapcore.Field, 0, len(fields))
    
    for _, field := range fields {
        // Redact sensitive fields
        if field.Key == "password" || field.Key == "secret" {
            field.String = "[REDACTED]"
        }
        filtered = append(filtered, field)
    }
    
    return e.Encoder.EncodeEntry(entry, filtered)
}

func (e *SensitiveFieldEncoder) Clone() zapcore.Encoder {
    return &SensitiveFieldEncoder{
        Encoder: e.Encoder.Clone(),
        cfg:     e.cfg,
    }
}

func NewSensitiveFieldsEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
    encoder := zapcore.NewJSONEncoder(config)
    return &SensitiveFieldEncoder{encoder, config}
}
```

### Hiding Sensitive Data with Stringer

```go
type User struct {
    ID    string
    Name  string
    Email string
}

// Implement fmt.Stringer to control logging output
func (u User) String() string {
    return u.ID // Only log ID, hide email
}

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    user := User{
        ID:    "USR-12345",
        Name:  "John Doe",
        Email: "john.doe@example.com",
    }
    
    logger.Info("user login", zap.Any("user", user))
    // Output: {"level":"info","msg":"user login","user":"USR-12345"}
}
```

---

## Output Configuration

### Multiple Output Destinations

```go
config := zap.NewProductionConfig()
config.OutputPaths = []string{
    "stdout",
    "/var/log/myapp.log",
    "kafka://localhost:9092/logs", // Custom sink
}
config.ErrorOutputPaths = []string{
    "stderr",
    "/var/log/myapp-errors.log",
}

logger, _ := config.Build()
```

### Tee to Multiple Outputs

```go
func createLogger() *zap.Logger {
    stdout := zapcore.AddSync(os.Stdout)
    file := zapcore.AddSync(&lumberjack.Logger{
        Filename: "app.log",
        MaxSize:  10,
    })
    
    consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
    jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
    
    core := zapcore.NewTee(
        zapcore.NewCore(consoleEncoder, stdout, zap.DebugLevel),
        zapcore.NewCore(jsonEncoder, file, zap.InfoLevel),
    )
    
    return zap.New(core)
}
```

---

## Log Rotation

### Using Lumberjack

```go
import "gopkg.in/natefinch/lumberjack.v2"

func createLogger() *zap.Logger {
    w := zapcore.AddSync(&lumberjack.Logger{
        Filename:   "/var/log/myapp/app.log",
        MaxSize:    100, // megabytes
        MaxBackups: 3,
        MaxAge:     28,  // days
        Compress:   true,
    })
    
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        w,
        zap.InfoLevel,
    )
    
    return zap.New(core)
}
```

### External Log Rotation

For production, use external tools like `logrotate`:

```
/var/log/myapp/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    postrotate
        killall -SIGUSR1 myapp
    endscript
}
```

---

## Sampling

Reduce log volume by sampling similar log entries.

### Basic Sampling

```go
func createLogger() *zap.Logger {
    stdout := zapcore.AddSync(os.Stdout)
    level := zap.NewAtomicLevelAt(zap.InfoLevel)
    
    jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
    jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)
    
    samplingCore := zapcore.NewSamplerWithOptions(
        jsonOutCore,
        time.Second,  // interval
        3,            // log first 3 entries
        0,            // log zero entries thereafter within interval
    )
    
    return zap.New(samplingCore)
}

func main() {
    logger := createLogger()
    defer logger.Sync()
    
    for i := 1; i <= 10; i++ {
        logger.Info("an info message")
        logger.Warn("a warning")
    }
    // Only first 3 iterations produce output
}
```

### Advanced Sampling

```go
samplingCore := zapcore.NewSamplerWithOptions(
    core,
    time.Second,
    100,  // First 100 entries per second
    100,  // Then 100 more (10% sample rate after first 100)
    zapcore.SamplerHook(func(entry zapcore.Entry, decision zapcore.SamplingDecision) {
        if decision == zapcore.LogDropped {
            // Track dropped logs
        }
    }),
)
```

---

## Error Handling

### Logging Errors

```go
logger.Error("Failed to perform an operation",
    zap.String("operation", "someOperation"),
    zap.Error(err),
    zap.Int("retryAttempts", 3),
    zap.String("user", "john.doe"),
)
```

**Output:**
```json
{
  "level":"error",
  "ts":1684164638.0570025,
  "caller":"main.go:47",
  "msg":"Failed to perform an operation",
  "operation":"someOperation",
  "error":"something happened",
  "retryAttempts":3,
  "user":"john.doe",
  "stacktrace":"main.main\n\t/home/user/main.go:47\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:250"
}
```

### Wrapping Errors

```go
func processData(id string) error {
    if err := fetchData(id); err != nil {
        logger.Error("Failed to fetch data",
            zap.String("id", id),
            zap.Error(err),
        )
        return fmt.Errorf("process data: %w", err)
    }
    return nil
}
```

### Fatal and Panic Handling

```go
// Fatal: logs and calls os.Exit(1)
logger.Fatal("Cannot start application",
    zap.String("reason", "port already in use"),
    zap.Int("port", 8080),
)

// Panic: logs and calls panic()
logger.Panic("Critical system failure",
    zap.Error(err),
)

// DPanic: panics in dev, error in production
logger.DPanic("This should never happen",
    zap.String("state", "invalid"),
)
```

---

## Production Patterns

### Singleton Logger

```go
package logger

import (
    "sync"
    "go.uber.org/zap"
)

var (
    logger *zap.Logger
    once   sync.Once
)

func Get() *zap.Logger {
    once.Do(func() {
        var err error
        if os.Getenv("ENV") == "production" {
            logger, err = zap.NewProduction()
        } else {
            logger, err = zap.NewDevelopment()
        }
        if err != nil {
            panic(err)
        }
    })
    return logger
}
```

### Global Logger Setup

```go
func init() {
    logger := createLogger()
    zap.ReplaceGlobals(logger)
}

func main() {
    defer zap.L().Sync()
    
    // Use global logger anywhere
    zap.L().Info("Application started")
    zap.S().Infow("User action", "action", "login")
}
```

### Graceful Shutdown

```go
func main() {
    logger := logger.Get()
    defer func() {
        if err := logger.Sync(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
        }
    }()
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        logger.Info("Shutting down gracefully")
        logger.Sync()
        os.Exit(0)
    }()
    
    // Application code
}
```

### HTTP Middleware Example

```go
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Wrap response writer to capture status
            ww := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            // Add logger to context
            requestID := uuid.New().String()
            reqLogger := logger.With(
                zap.String("request_id", requestID),
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
            )
            
            ctx := context.WithValue(r.Context(), loggerKey, reqLogger)
            
            // Process request
            next.ServeHTTP(ww, r.WithContext(ctx))
            
            // Log completion
            reqLogger.Info("Request completed",
                zap.Int("status", ww.statusCode),
                zap.Duration("duration", time.Since(start)),
            )
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

---

## Testing

### Capturing Logs in Tests

```go
import (
    "testing"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "go.uber.org/zap/zaptest/observer"
)

func TestLogging(t *testing.T) {
    // Create test logger with observer
    core, recorded := observer.New(zapcore.InfoLevel)
    logger := zap.New(core)
    
    // Code under test
    logger.Info("test message", zap.String("key", "value"))
    logger.Error("error message", zap.Error(errors.New("test error")))
    
    // Assertions
    logs := recorded.All()
    if len(logs) != 2 {
        t.Errorf("Expected 2 logs, got %d", len(logs))
    }
    
    if logs[0].Message != "test message" {
        t.Errorf("Expected 'test message', got '%s'", logs[0].Message)
    }
    
    if logs[0].Level != zapcore.InfoLevel {
        t.Errorf("Expected InfoLevel, got %v", logs[0].Level)
    }
    
    // Check fields
    fields := logs[0].Context
    if fields[0].Key != "key" || fields[0].String != "value" {
        t.Error("Field mismatch")
    }
}
```

### Testing with Custom Writer

```go
func TestLogOutput(t *testing.T) {
    var buf bytes.Buffer
    
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        zapcore.AddSync(&buf),
        zapcore.InfoLevel,
    )
    
    logger := zap.New(core)
    logger.Info("test message")
    
    output := buf.String()
    if !strings.Contains(output, "test message") {
        t.Errorf("Log output missing expected message: %s", output)
    }
}
```

### Mock Logger

```go
type MockLogger struct {
    Entries []Entry
}

type Entry struct {
    Level   string
    Message string
    Fields  map[string]interface{}
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
    m.Entries = append(m.Entries, Entry{
        Level:   "info",
        Message: msg,
        Fields:  fieldsToMap(fields),
    })
}
```

---

## Performance Optimization

### Benchmark Your Logging

```go
func BenchmarkZapLogger(b *testing.B) {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.Info("Benchmark log message",
            zap.Int("iteration", i),
            zap.String("status", "running"),
        )
    }
}
```

### Buffered Logging

```go
func createBufferedLogger() *zap.Logger {
    bufferedWriter := &zapcore.BufferedWriteSyncer{
        WS:   zapcore.AddSync(os.Stdout),
        Size: 256 * 1024, // 256KB buffer
    }
    
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        bufferedWriter,
        zap.InfoLevel,
    )
    
    return zap.New(core)
}
```

### Avoid Expensive Operations

```go
// BAD: Always evaluates expensive function
logger.Debug("User data", zap.Any("data", fetchExpensiveUserData()))

// GOOD: Only evaluates when debug is enabled
if ce := logger.Check(zap.DebugLevel, "User data"); ce != nil {
    ce.Write(zap.Any("data", fetchExpensiveUserData()))
}
```

---

## Best Practices

### 1. Use Appropriate Log Levels

```go
logger.Debug("Detailed diagnostic info")    // Development only
logger.Info("Normal operations")            // Important events
logger.Warn("Unusual but handled")          // Attention needed
logger.Error("Errors that need fixing")     // Requires action
logger.Fatal("Cannot continue")             // Immediate exit
```

### 2. Add Context, Not Clutter

```go
// GOOD: Relevant context
logger.Info("Payment processed",
    zap.String("order_id", orderID),
    zap.Float64("amount", amount),
    zap.String("currency", "USD"),
)

// BAD: Too much detail
logger.Info("Payment processed",
    zap.Any("entire_request_object", request),
    zap.Any("entire_database_record", record),
)
```

### 3. Never Log Sensitive Data

```go
// NEVER
logger.Info("User login", zap.String("password", password))

// ALWAYS redact or omit
logger.Info("User login", zap.String("user_id", userID))
```

### 4. Use Structured Fields

```go
// GOOD: Structured
logger.Error("Database error",
    zap.Error(err),
    zap.String("query", query),
    zap.Duration("elapsed", elapsed),
)

// BAD: Unstructured
logger.Errorf("Database error: %v (query: %s, took: %v)", err, query, elapsed)
```

### 5. Log at Boundaries

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    logger.Info("Request received")
    
    // Process...
    
    logger.Info("Request completed", zap.Int("status", 200))
}
```

### 6. Include Request IDs

```go
requestLogger := logger.With(zap.String("request_id", requestID))
requestLogger.Info("Processing started")
// All subsequent logs include request_id
```

### 7. Don't Log in Tight Loops

```go
// BAD
for i := 0; i < 1000000; i++ {
    logger.Debug("Processing", zap.Int("i", i))
}

// GOOD: Sample or aggregate
if i%10000 == 0 {
    logger.Debug("Progress", zap.Int("processed", i))
}
```

### 8. Always Sync on Shutdown

```go
func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync() // Essential!
    
    // Application code
}
```

---

## Resources

### Official Documentation
- [GitHub Repository](https://github.com/uber-go/zap)
- [Package Documentation](https://pkg.go.dev/go.uber.org/zap)
- [FAQ](https://github.com/uber-go/zap/blob/master/FAQ.md)

### Guides
- [BetterStack Zap Guide](https://betterstack.com/community/guides/logging/go/zap/)
- [SigNoz Zap Guide](https://signoz.io/guides/zap-logger/)
- [Dash0 Zap Guide](https://www.dash0.com/guides/logging-in-go-with-zap)

---

## Conclusion

Zap is the gold standard for logging in Go applications, offering:

- **Performance**: 4-10x faster than alternatives
- **Type Safety**: Compile-time guarantees with strongly-typed fields
- **Flexibility**: From development to production configurations
- **Observability**: Structured logs for modern monitoring systems

By following this guide and adopting Zap's best practices, you'll build Go applications with world-class logging that aids debugging, monitoring, and operational excellence.

Happy logging! ⚡
