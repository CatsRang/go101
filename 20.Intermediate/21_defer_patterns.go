package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// DEFER PATTERNS IN GO
// ============================================================================
// defer schedules a function call to run after the surrounding function returns
// Deferred calls execute in LIFO order (Last In, First Out)

// ============================================================================
// 1. BASIC DEFER
// ============================================================================

func basicDefer() {
	defer fmt.Println("This executes last")
	fmt.Println("This executes first")
	fmt.Println("This executes second")
}

// ============================================================================
// 2. DEFER ORDER (LIFO)
// ============================================================================

func deferOrder() {
	defer fmt.Println("Defer 1")
	defer fmt.Println("Defer 2")
	defer fmt.Println("Defer 3")
	fmt.Println("Function body")
	// Output: Function body, Defer 3, Defer 2, Defer 1
}

// ============================================================================
// 3. RESOURCE CLEANUP PATTERN
// ============================================================================

type File struct {
	name   string
	open   bool
	data   string
}

func (f *File) Open() error {
	if f.open {
		return fmt.Errorf("file already open")
	}
	f.open = true
	fmt.Printf("Opening file: %s\n", f.name)
	return nil
}

func (f *File) Write(data string) error {
	if !f.open {
		return fmt.Errorf("file not open")
	}
	f.data += data
	fmt.Printf("Writing to %s: %s\n", f.name, data)
	return nil
}

func (f *File) Close() error {
	if !f.open {
		return fmt.Errorf("file not open")
	}
	f.open = false
	fmt.Printf("Closing file: %s\n", f.name)
	return nil
}

func processFile() error {
	file := &File{name: "data.txt"}

	if err := file.Open(); err != nil {
		return err
	}
	// Defer ensures file is closed even if errors occur later
	defer file.Close()

	if err := file.Write("Hello, "); err != nil {
		return err
	}

	if err := file.Write("World!"); err != nil {
		return err
	}

	return nil
}

// ============================================================================
// 4. MUTEX UNLOCK PATTERN
// ============================================================================

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock() // Ensures unlock even if panic occurs

	sc.count++
	fmt.Printf("Count: %d\n", sc.count)
}

func (sc *SafeCounter) IncrementComplex() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Complex logic with multiple return paths
	if sc.count >= 10 {
		fmt.Println("Counter limit reached")
		return // defer ensures unlock happens
	}

	sc.count++

	if sc.count%2 == 0 {
		fmt.Println("Even count")
		return // defer ensures unlock happens
	}

	fmt.Println("Odd count")
	// defer ensures unlock happens
}

// ============================================================================
// 5. TIMING FUNCTION EXECUTION
// ============================================================================

func measureTime(name string) func() {
	start := time.Now()
	fmt.Printf("Starting %s\n", name)

	// Return a function that will be called when deferred
	return func() {
		fmt.Printf("Finished %s (took %v)\n", name, time.Since(start))
	}
}

func slowOperation() {
	defer measureTime("slowOperation")()

	// Simulate work
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Doing work...")
}

// ============================================================================
// 6. DEFER WITH NAMED RETURN VALUES
// ============================================================================

func divide(a, b int) (result int, err error) {
	defer func() {
		if err != nil {
			// Modify error before returning
			err = fmt.Errorf("division failed: %w", err)
		}
	}()

	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	return a / b, nil
}

// ============================================================================
// 7. DEFER WITH CLOSURES (CAPTURING VALUES)
// ============================================================================

// ❌ BAD: Captures variable reference
func badDeferLoop() {
	for i := 0; i < 3; i++ {
		defer fmt.Printf("Bad: %d\n", i)
	}
	// Prints: Bad: 2, Bad: 1, Bad: 0 (LIFO order, captures final value)
}

// ✓ GOOD: Pass value to avoid capture issues
func goodDeferLoop() {
	for i := 0; i < 3; i++ {
		i := i // Create new variable for each iteration
		defer func() {
			fmt.Printf("Good: %d\n", i)
		}()
	}
}

// Alternative: Use parameters
func betterDeferLoop() {
	for i := 0; i < 3; i++ {
		defer func(val int) {
			fmt.Printf("Better: %d\n", val)
		}(i)
	}
}

// ============================================================================
// 8. DEFER FOR ERROR HANDLING
// ============================================================================

func processWithErrorHandling() (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("Error occurred: %v\n", err)
			// Could log to file, send alert, etc.
		}
	}()

	// Simulate operations
	return fmt.Errorf("something went wrong")
}

// ============================================================================
// 9. DEFER WITH PANIC RECOVERY
// ============================================================================

func safeOperation() (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// This will panic
	panic("critical error")
	return nil
}

// ============================================================================
// 10. TRANSACTION PATTERN
// ============================================================================

type Transaction struct {
	active    bool
	committed bool
}

func (tx *Transaction) Begin() error {
	tx.active = true
	fmt.Println("Transaction started")
	return nil
}

func (tx *Transaction) Commit() error {
	if !tx.active {
		return fmt.Errorf("no active transaction")
	}
	tx.committed = true
	tx.active = false
	fmt.Println("Transaction committed")
	return nil
}

func (tx *Transaction) Rollback() error {
	if !tx.active {
		return nil // Already rolled back or committed
	}
	tx.active = false
	fmt.Println("Transaction rolled back")
	return nil
}

func performTransaction(shouldCommit bool) error {
	tx := &Transaction{}
	tx.Begin()

	// Defer rollback - only happens if not committed
	defer func() {
		if tx.active {
			tx.Rollback()
		}
	}()

	// Do work
	fmt.Println("Performing database operations...")

	if !shouldCommit {
		return fmt.Errorf("operation failed")
	}

	// Explicitly commit
	return tx.Commit()
}

// ============================================================================
// 11. CONTEXT CLEANUP
// ============================================================================

type Context struct {
	name      string
	resources []string
}

func (c *Context) AddResource(resource string) {
	c.resources = append(c.resources, resource)
	fmt.Printf("[%s] Acquired: %s\n", c.name, resource)
}

func (c *Context) Cleanup() {
	fmt.Printf("[%s] Cleaning up %d resources\n", c.name, len(c.resources))
	for i := len(c.resources) - 1; i >= 0; i-- {
		fmt.Printf("[%s] Releasing: %s\n", c.name, c.resources[i])
	}
	c.resources = nil
}

func processWithContext() {
	ctx := &Context{name: "MyProcess"}
	defer ctx.Cleanup()

	ctx.AddResource("Database connection")
	ctx.AddResource("File handle")
	ctx.AddResource("Network socket")

	fmt.Println("Processing...")
}

// ============================================================================
// 12. DEFER EVALUATION TIME
// ============================================================================

func deferEvaluation() {
	x := 10
	defer fmt.Printf("Defer sees x = %d\n", x) // x is evaluated now (10)

	x = 20
	fmt.Printf("x is now %d\n", x)
	// When defer executes, it uses the value captured earlier (10)
}

func deferWithClosure() {
	x := 10
	defer func() {
		fmt.Printf("Closure sees x = %d\n", x) // Captures variable, not value
	}()

	x = 20
	fmt.Printf("x is now %d\n", x)
	// When defer executes, it sees current value (20)
}

// ============================================================================
// 13. MULTIPLE RESOURCES
// ============================================================================

func processMultipleResources() error {
	file1 := &File{name: "file1.txt"}
	if err := file1.Open(); err != nil {
		return err
	}
	defer file1.Close()

	file2 := &File{name: "file2.txt"}
	if err := file2.Open(); err != nil {
		return err // file1 will still be closed
	}
	defer file2.Close()

	file3 := &File{name: "file3.txt"}
	if err := file3.Open(); err != nil {
		return err // file1 and file2 will be closed
	}
	defer file3.Close()

	// All files will be closed in reverse order: file3, file2, file1
	fmt.Println("All files opened successfully")
	return nil
}

func main() {
	fmt.Println("=== 1. Basic Defer ===")
	basicDefer()

	fmt.Println("\n=== 2. Defer Order (LIFO) ===")
	deferOrder()

	fmt.Println("\n=== 3. Resource Cleanup ===")
	processFile()

	fmt.Println("\n=== 4. Mutex Unlock Pattern ===")
	counter := &SafeCounter{}
	counter.Increment()
	counter.Increment()
	counter.IncrementComplex()

	fmt.Println("\n=== 5. Timing Function Execution ===")
	slowOperation()

	fmt.Println("\n=== 6. Defer with Named Returns ===")
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== 7. Defer in Loops ===")
	fmt.Println("Bad defer loop:")
	badDeferLoop()

	fmt.Println("\nGood defer loop:")
	goodDeferLoop()

	fmt.Println("\nBetter defer loop:")
	betterDeferLoop()

	fmt.Println("\n=== 8. Defer for Error Handling ===")
	processWithErrorHandling()

	fmt.Println("\n=== 9. Defer with Panic Recovery ===")
	err = safeOperation()
	fmt.Printf("Returned error: %v\n", err)

	fmt.Println("\n=== 10. Transaction Pattern ===")
	fmt.Println("Successful transaction:")
	performTransaction(true)

	fmt.Println("\nFailed transaction (will rollback):")
	performTransaction(false)

	fmt.Println("\n=== 11. Context Cleanup ===")
	processWithContext()

	fmt.Println("\n=== 12. Defer Evaluation Time ===")
	deferEvaluation()
	fmt.Println()
	deferWithClosure()

	fmt.Println("\n=== 13. Multiple Resources ===")
	processMultipleResources()

	fmt.Println("\n=== DEFER BEST PRACTICES ===")
	fmt.Println("✓ Use defer for cleanup (close, unlock, rollback)")
	fmt.Println("✓ Place defer immediately after acquiring resource")
	fmt.Println("✓ Remember LIFO order for multiple defers")
	fmt.Println("✓ Use defer with named return values for error handling")
	fmt.Println("✓ Be careful with defer in loops (can accumulate)")
	fmt.Println("✓ Arguments are evaluated immediately, not at defer time")
	fmt.Println("✓ Closures capture variables, not values")
	fmt.Println("✓ defer is perfect for unlock, close, rollback patterns")
	fmt.Println("✗ Don't use defer for performance-critical code")
	fmt.Println("✗ Don't defer in tight loops (memory accumulation)")

	fmt.Println("\n=== COMMON DEFER PATTERNS ===")
	fmt.Println("1. Resource cleanup: defer file.Close()")
	fmt.Println("2. Mutex unlock: defer mu.Unlock()")
	fmt.Println("3. Panic recovery: defer func() { recover() }()")
	fmt.Println("4. Transaction rollback: defer tx.Rollback()")
	fmt.Println("5. Timer measurement: defer measureTime()()")
	fmt.Println("6. Error logging: defer logError()")
	fmt.Println("7. Context cleanup: defer ctx.Cleanup()")
}
