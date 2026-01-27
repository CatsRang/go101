package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// The context package provides a way to:
// - Carry deadlines, cancellation signals, and request-scoped values
// - Propagate cancellation across API boundaries and goroutines

// Define context key type at package level
type ctxKey string

const (
	userIDKey    ctxKey = "userID"
	requestIDKey ctxKey = "requestID"
)

func main() {
	// === context.Background() ===
	fmt.Println("=== context.Background() ===")
	ctx := context.Background()
	fmt.Printf("Background context: %v\n", ctx)

	// === context.WithCancel ===
	fmt.Println("\n=== context.WithCancel ===")

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker: received cancellation signal")
				return
			default:
				fmt.Println("Worker: working...")
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(150 * time.Millisecond)
	cancel() // Cancel the context
	time.Sleep(50 * time.Millisecond)

	// === context.WithTimeout ===
	fmt.Println("\n=== context.WithTimeout ===")

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // Always call cancel to release resources

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Operation completed")
	case <-ctx.Done():
		fmt.Printf("Context cancelled: %v\n", ctx.Err())
	}

	// === context.WithDeadline ===
	fmt.Println("\n=== context.WithDeadline ===")

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel = context.WithDeadline(context.Background(), deadline)
	defer cancel()

	dl, ok := ctx.Deadline()
	fmt.Printf("Deadline: %v, Set: %v\n", dl.Format(time.RFC3339Nano), ok)

	<-ctx.Done()
	fmt.Printf("Deadline exceeded: %v\n", ctx.Err())

	// === context.WithValue ===
	fmt.Println("\n=== context.WithValue ===")

	ctx = context.WithValue(context.Background(), userIDKey, "user-123")
	ctx = context.WithValue(ctx, requestIDKey, "req-456")

	processRequest(ctx)

	// === Go 1.21+: context.WithoutCancel ===
	fmt.Println("\n=== Go 1.21+: context.WithoutCancel ===")

	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Create a context that won't be cancelled when parent is cancelled
	// Useful for cleanup operations that should complete regardless
	detachedCtx := context.WithoutCancel(parentCtx)

	parentCancel() // Cancel parent

	fmt.Printf("Parent cancelled: %v\n", parentCtx.Err())
	fmt.Printf("Detached context cancelled: %v\n", detachedCtx.Err())

	// === Go 1.21+: context.AfterFunc ===
	fmt.Println("\n=== Go 1.21+: context.AfterFunc ===")

	ctx, cancel = context.WithCancel(context.Background())

	// Register a function to be called when context is done
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("AfterFunc: Context was cancelled, running cleanup!")
	})

	// Can optionally stop the AfterFunc from running
	_ = stop // stop() would prevent the AfterFunc from running

	cancel() // Triggers the AfterFunc
	time.Sleep(50 * time.Millisecond)

	// === Go 1.20+: context.WithCancelCause ===
	fmt.Println("\n=== Go 1.20+: context.WithCancelCause ===")

	ctx, cancelCause := context.WithCancelCause(context.Background())
	customErr := errors.New("custom cancellation reason")
	cancelCause(customErr)

	fmt.Printf("Context error: %v\n", ctx.Err())
	fmt.Printf("Cause: %v\n", context.Cause(ctx))

	// === Practical Example: HTTP-like Request with Timeout ===
	fmt.Println("\n=== Practical: Simulated HTTP Request ===")

	ctx, cancel = context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	result, err := simulateHTTPRequest(ctx, "https://api.example.com/data")
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", result)
	}

	fmt.Println("\nContext examples completed!")
}

func processRequest(ctx context.Context) {
	userID := ctx.Value(userIDKey)
	requestID := ctx.Value(requestIDKey)
	fmt.Printf("Processing request %v for user %v\n", requestID, userID)
}

func simulateHTTPRequest(ctx context.Context, url string) (string, error) {
	// Simulate network latency
	resultCh := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond) // Simulated work
		resultCh <- fmt.Sprintf("Data from %s", url)
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
