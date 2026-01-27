package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Go 1.21+ introduced log/slog for structured logging
// Benefits: structured data, log levels, multiple handlers, better performance

func main() {
	fmt.Println("=== Go 1.21+ Structured Logging with slog ===\n")

	// === Default Logger (Text Handler) ===
	fmt.Println("--- Default Logger ---")

	slog.Info("Application started")
	slog.Debug("Debug message (not shown with default level)")
	slog.Warn("This is a warning")
	slog.Error("An error occurred")

	// === Logger with Attributes ===
	fmt.Println("\n--- Logger with Attributes ---")

	slog.Info("User action",
		slog.String("user_id", "user-123"),
		slog.String("action", "login"),
		slog.Int("attempt", 1),
	)

	// === JSON Handler ===
	fmt.Println("\n--- JSON Handler ---")

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Show all levels including Debug
	})
	jsonLogger := slog.New(jsonHandler)

	jsonLogger.Info("Processing request",
		slog.String("method", "GET"),
		slog.String("path", "/api/users"),
		slog.Duration("latency", 45*time.Millisecond),
	)

	// === Text Handler with Options ===
	fmt.Println("\n--- Text Handler with Source ---")

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // Include source file and line number
	})
	textLogger := slog.New(textHandler)

	textLogger.Debug("Debugging with source info")

	// === Grouped Attributes ===
	fmt.Println("\n--- Grouped Attributes ---")

	jsonLogger.Info("HTTP request",
		slog.Group("request",
			slog.String("method", "POST"),
			slog.String("path", "/api/orders"),
			slog.String("ip", "192.168.1.1"),
		),
		slog.Group("response",
			slog.Int("status", 201),
			slog.Duration("duration", 123*time.Millisecond),
		),
	)

	// === Logger with Default Attributes ===
	fmt.Println("\n--- Logger with Default Attributes ---")

	// Create logger with default attributes added to every log
	serviceLogger := jsonLogger.With(
		slog.String("service", "order-service"),
		slog.String("version", "1.2.3"),
	)

	serviceLogger.Info("Order created", slog.String("order_id", "ord-456"))
	serviceLogger.Info("Order shipped", slog.String("order_id", "ord-456"))

	// === Context-Aware Logging ===
	fmt.Println("\n--- Context-Aware Logging ---")

	// Create a context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "req-789")

	// Log with context
	jsonLogger.InfoContext(ctx, "Processing with context",
		slog.String("user", "alice"),
	)

	// === Error Logging with Details ===
	fmt.Println("\n--- Error Logging ---")

	err := errors.New("connection refused")
	jsonLogger.Error("Database connection failed",
		slog.String("host", "localhost"),
		slog.Int("port", 5432),
		slog.String("error", err.Error()),
		slog.Int("retry_count", 3),
	)

	// === Custom Log Levels ===
	fmt.Println("\n--- Log Levels ---")

	// Create handler that shows all levels
	verboseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	verboseLogger := slog.New(verboseHandler)

	verboseLogger.Debug("Debug level - detailed debugging info")
	verboseLogger.Info("Info level - general information")
	verboseLogger.Warn("Warn level - warning conditions")
	verboseLogger.Error("Error level - error conditions")

	// === Practical Example: HTTP Handler Logging ===
	fmt.Println("\n--- Practical: HTTP Handler Logging Example ---")

	// Simulate HTTP request logging
	simulateHTTPLogging(jsonLogger)

	// === Set Default Logger ===
	fmt.Println("\n--- Setting Default Logger ---")

	// Set as default logger for slog package functions
	slog.SetDefault(jsonLogger)

	// Now slog.Info uses jsonLogger
	slog.Info("Using new default logger")

	// === LogValuer Interface ===
	fmt.Println("\n--- Custom LogValuer ---")

	user := &LoggableUser{
		ID:       1,
		Username: "alice",
		Password: "secret123", // Should not be logged
		Email:    "alice@example.com",
	}

	jsonLogger.Info("User logged in", slog.Any("user", user))

	fmt.Println("\nStructured logging examples completed!")
}

// LoggableUser implements slog.LogValuer to control what gets logged
type LoggableUser struct {
	ID       int
	Username string
	Password string // Sensitive - should not be logged
	Email    string
}

// LogValue returns the value to log, excluding sensitive fields
func (u *LoggableUser) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", u.ID),
		slog.String("username", u.Username),
		slog.String("email", u.Email),
		// Password is intentionally omitted
	)
}

// simulateHTTPLogging demonstrates logging in HTTP handlers
func simulateHTTPLogging(logger *slog.Logger) {
	// Simulate a request
	req := &http.Request{
		Method:     "GET",
		RequestURI: "/api/products/123",
		RemoteAddr: "192.168.1.100:54321",
	}

	start := time.Now()

	// Log request start
	logger.Info("Request received",
		slog.Group("request",
			slog.String("method", req.Method),
			slog.String("uri", req.RequestURI),
			slog.String("remote_addr", req.RemoteAddr),
		),
	)

	// Simulate processing
	time.Sleep(50 * time.Millisecond)

	// Log request completion
	logger.Info("Request completed",
		slog.String("method", req.Method),
		slog.String("uri", req.RequestURI),
		slog.Int("status", 200),
		slog.Duration("duration", time.Since(start)),
	)
}
