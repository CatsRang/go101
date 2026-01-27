package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware is a function that wraps an http.Handler
// Common uses: logging, authentication, CORS, rate limiting, etc.

// Middleware type definition
type Middleware func(http.Handler) http.Handler

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /api/data", apiDataHandler)
	mux.HandleFunc("GET /api/protected", protectedHandler)

	// Apply middleware chain
	// Order matters: first applied = outermost wrapper
	handler := Chain(
		mux,
		LoggingMiddleware,
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
	)

	// For protected routes, you might apply auth middleware selectively
	// Example of route-specific middleware is shown in protectedHandler

	addr := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Println("\nTest with:")
	fmt.Println("  curl -v http://localhost:8080/")
	fmt.Println("  curl -v http://localhost:8080/api/data")
	fmt.Println("  curl -v http://localhost:8080/api/protected")
	fmt.Println("  curl -v -H 'Authorization: Bearer secret-token' http://localhost:8080/api/protected")

	log.Fatal(http.ListenAndServe(addr, handler))
}

// Chain applies middlewares in order
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply in reverse order so first middleware is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// === Logging Middleware ===
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log after request completes
		log.Printf(
			"[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// === Recovery Middleware (Panic Handler) ===
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// === CORS Middleware ===
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// === Request ID Middleware ===
type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate or extract request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// === Authentication Middleware (Example) ===
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// Simple token check (use proper JWT/OAuth in production)
		if token != "Bearer secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user", "authenticated-user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// === Rate Limiting Middleware (Simple In-Memory) ===
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Clean old requests
		now := time.Now()
		windowStart := now.Add(-rl.window)

		validRequests := []time.Time{}
		for _, t := range rl.requests[ip] {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}
		rl.requests[ip] = validRequests

		// Check limit
		if len(rl.requests[ip]) >= rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Record this request
		rl.requests[ip] = append(rl.requests[ip], now)

		next.ServeHTTP(w, r)
	})
}

// === Handlers ===

func homeHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(RequestIDKey)
	fmt.Fprintf(w, "Welcome! Request ID: %v\n", requestID)
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "This is public API data"}`)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Apply auth middleware inline for this specific route
	authHandler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Protected data", "user": "%v"}`, user)
	}))

	authHandler.ServeHTTP(w, r)
}
