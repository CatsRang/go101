package main

import (
	"fmt"
	"runtime/debug"
)

// ============================================================================
// PANIC AND RECOVER IN GO
// ============================================================================
// panic() stops normal execution and begins panicking
// recover() regains control of a panicking goroutine
// defer is required for recover() to work

// ============================================================================
// 1. BASIC PANIC
// ============================================================================

func demonstratePanic() {
	fmt.Println("Before panic")
	panic("something went wrong!")
	fmt.Println("After panic - this will NOT execute")
}

// ============================================================================
// 2. BASIC RECOVER
// ============================================================================

func safeFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println("About to panic...")
	panic("oops!")
	fmt.Println("This won't execute")
}

// ============================================================================
// 3. PANIC WITH DIFFERENT TYPES
// ============================================================================

func panicWithString() {
	panic("string panic")
}

func panicWithInt() {
	panic(42)
}

func panicWithError() {
	panic(fmt.Errorf("error panic"))
}

type CustomPanic struct {
	Message string
	Code    int
}

func panicWithStruct() {
	panic(CustomPanic{
		Message: "custom panic",
		Code:    500,
	})
}

func demonstratePanicTypes() {
	defer func() {
		if r := recover(); r != nil {
			// Type switch on panic value
			switch v := r.(type) {
			case string:
				fmt.Printf("String panic: %s\n", v)
			case int:
				fmt.Printf("Int panic: %d\n", v)
			case error:
				fmt.Printf("Error panic: %v\n", v)
			case CustomPanic:
				fmt.Printf("Custom panic: %s (code: %d)\n", v.Message, v.Code)
			default:
				fmt.Printf("Unknown panic type: %v\n", r)
			}
		}
	}()

	panicWithStruct()
}

// ============================================================================
// 4. RECOVER IN NESTED FUNCTIONS
// ============================================================================

func level3() {
	panic("panic at level 3")
}

func level2() {
	level3()
}

func level1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in level1: %v\n", r)
		}
	}()

	level2()
}

// ============================================================================
// 5. CONVERTING PANIC TO ERROR
// ============================================================================

func riskyOperation() (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	// Some operation that might panic
	panic("critical failure")
	return nil
}

// ============================================================================
// 6. PARTIAL RECOVERY - RE-PANIC
// ============================================================================

func handleSpecificPanics() {
	defer func() {
		if r := recover(); r != nil {
			// Only handle specific types of panics
			if str, ok := r.(string); ok && str == "expected panic" {
				fmt.Println("Handled expected panic")
				return
			}

			// Re-panic for unexpected panics
			fmt.Println("Unexpected panic, re-panicking...")
			panic(r)
		}
	}()

	panic("unexpected panic")
}

// ============================================================================
// 7. WEB SERVER PANIC RECOVERY
// ============================================================================

type Request struct {
	ID   int
	Path string
}

type Response struct {
	Status int
	Body   string
}

func handleRequest(req Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in request %d: %v\n", req.ID, r)

			// Log stack trace
			fmt.Printf("Stack trace:\n%s\n", debug.Stack())

			// Return 500 error response
			resp = Response{
				Status: 500,
				Body:   "Internal Server Error",
			}
		}
	}()

	// Simulate processing that might panic
	if req.Path == "/panic" {
		panic("intentional panic")
	}

	return Response{
		Status: 200,
		Body:   "OK",
	}
}

// ============================================================================
// 8. RESOURCE CLEANUP WITH PANIC
// ============================================================================

type Database struct {
	connected bool
}

func (db *Database) Connect() {
	fmt.Println("Database connected")
	db.connected = true
}

func (db *Database) Close() {
	if db.connected {
		fmt.Println("Database closed")
		db.connected = false
	}
}

func processWithCleanup() {
	db := &Database{}
	db.Connect()

	// Ensure cleanup happens even if panic occurs
	defer db.Close()

	// This will panic, but defer ensures Close() is called
	panic("database operation failed")
}

// ============================================================================
// 9. WHEN TO USE PANIC
// ============================================================================

// ✓ GOOD: Unrecoverable programming errors
func initializeConfig(configPath string) {
	if configPath == "" {
		// This is a programming error, not a runtime error
		panic("config path cannot be empty - this is a bug")
	}
	fmt.Printf("Loading config from %s\n", configPath)
}

// ❌ BAD: Using panic for normal errors
func badDivide(a, b int) int {
	if b == 0 {
		panic("division by zero") // Don't do this!
	}
	return a / b
}

// ✓ GOOD: Return error instead
func goodDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// ============================================================================
// 10. GOROUTINE PANIC HANDLING
// ============================================================================

func workerSafe(id int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Worker %d recovered from panic: %v\n", id, r)
		}
	}()

	if id == 2 {
		panic("worker 2 panicked")
	}

	fmt.Printf("Worker %d completed successfully\n", id)
}

func demonstrateGoroutinePanic() {
	// Launch multiple goroutines with panic protection
	for i := 1; i <= 3; i++ {
		go workerSafe(i)
	}

	// Give goroutines time to complete
	// In real code, use sync.WaitGroup
	fmt.Println("Main goroutine continues...")
}

// ============================================================================
// 11. PANIC IN DEFER
// ============================================================================

func panicInDefer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in defer: %v\n", r)
		}
	}()

	defer func() {
		panic("panic in defer")
	}()

	fmt.Println("Function body")
}

// ============================================================================
// 12. MULTIPLE DEFERS WITH PANIC
// ============================================================================

func multipleDeferWithPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Final recovery: %v\n", r)
		}
	}()

	defer fmt.Println("Defer 3: This executes")
	defer fmt.Println("Defer 2: This executes")
	defer fmt.Println("Defer 1: This executes")

	panic("panic occurs")

	fmt.Println("This does NOT execute")
}

// ============================================================================
// 13. STACK TRACE
// ============================================================================

func captureStackTrace() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %v\n", r)
			fmt.Println("\nStack trace:")
			fmt.Printf("%s\n", debug.Stack())
		}
	}()

	deepFunction()
}

func deepFunction() {
	evenDeeperFunction()
}

func evenDeeperFunction() {
	panic("deep panic")
}

func main() {
	fmt.Println("=== 1. Basic Panic (Commented to Avoid Crash) ===")
	// demonstratePanic() // This would crash the program
	fmt.Println("(Skipped to prevent crash)")

	fmt.Println("\n=== 2. Basic Recover ===")
	safeFunction()
	fmt.Println("Program continues after recovery")

	fmt.Println("\n=== 3. Panic with Different Types ===")
	demonstratePanicTypes()

	fmt.Println("\n=== 4. Nested Function Panic ===")
	level1()

	fmt.Println("\n=== 5. Converting Panic to Error ===")
	err := riskyOperation()
	fmt.Printf("Error: %v\n", err)

	fmt.Println("\n=== 6. Partial Recovery (Re-panic) ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Caught re-panicked error: %v\n", r)
			}
		}()
		handleSpecificPanics()
	}()

	fmt.Println("\n=== 7. Web Server Panic Recovery ===")
	requests := []Request{
		{ID: 1, Path: "/api/users"},
		{ID: 2, Path: "/panic"},
		{ID: 3, Path: "/api/data"},
	}

	for _, req := range requests {
		resp := handleRequest(req)
		fmt.Printf("Request %d → Status: %d, Body: %s\n", req.ID, resp.Status, resp.Body)
	}

	fmt.Println("\n=== 8. Resource Cleanup ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered: %v (cleanup happened)\n", r)
			}
		}()
		processWithCleanup()
	}()

	fmt.Println("\n=== 9. Good vs Bad Panic Usage ===")
	initializeConfig("/etc/config.json")

	// Good: using errors
	result, err := goodDivide(10, 0)
	if err != nil {
		fmt.Printf("Division error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("\n=== 10. Goroutine Panic (Output May Vary) ===")
	demonstrateGoroutinePanic()

	fmt.Println("\n=== 11. Panic in Defer ===")
	panicInDefer()

	fmt.Println("\n=== 12. Multiple Defers with Panic ===")
	multipleDeferWithPanic()

	fmt.Println("\n=== 13. Stack Trace ===")
	captureStackTrace()

	fmt.Println("\n=== PANIC AND RECOVER RULES ===")
	fmt.Println("✓ recover() only works inside defer")
	fmt.Println("✓ recover() only catches panics in the same goroutine")
	fmt.Println("✓ Each goroutine should have its own panic recovery")
	fmt.Println("✓ Use defer to ensure cleanup happens")
	fmt.Println("✓ Convert panics to errors at API boundaries")
	fmt.Println("✓ Log stack traces for debugging")
	fmt.Println("✓ Re-panic if you can't handle it")
	fmt.Println("✗ Don't use panic for normal error handling")
	fmt.Println("✗ Don't ignore recovered panics (at least log them)")
	fmt.Println("✗ Don't rely on panic for control flow")

	fmt.Println("\n=== WHEN TO USE PANIC ===")
	fmt.Println("Use panic for:")
	fmt.Println("  • Unrecoverable errors during initialization")
	fmt.Println("  • Programming errors (bugs in the code)")
	fmt.Println("  • Impossible situations (should never happen)")
	fmt.Println("  • Invariant violations")
	fmt.Println("\nUse errors for:")
	fmt.Println("  • Expected error conditions")
	fmt.Println("  • Network/IO failures")
	fmt.Println("  • Validation failures")
	fmt.Println("  • Business logic errors")
	fmt.Println("  • User input errors")
	fmt.Println("\nRecover when:")
	fmt.Println("  • In web servers (don't crash on request panic)")
	fmt.Println("  • In worker pools (don't kill whole pool)")
	fmt.Println("  • At application boundaries")
	fmt.Println("  • Converting library panics to errors")
}
