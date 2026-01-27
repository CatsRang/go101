package main

import (
	"fmt"
	"time"
)

// Goroutines are lightweight threads managed by the Go runtime.
// They are much cheaper than OS threads (starting at ~2KB stack size).

func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Hello, %s! (iteration %d)\n", name, i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== Basic Goroutines ===")

	// Start a goroutine with the 'go' keyword
	go sayHello("Alice")
	go sayHello("Bob")

	// Main goroutine continues executing
	fmt.Println("Main goroutine is running...")

	// Wait for goroutines to finish (crude way - use sync.WaitGroup in production)
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n=== Anonymous Goroutines ===")

	// Anonymous function as goroutine
	go func() {
		fmt.Println("Anonymous goroutine running!")
	}()

	// Go 1.22+ Loop Variable Fix
	// Before Go 1.22, loop variables were shared across iterations
	// Now each iteration gets its own copy
	fmt.Println("\n=== Go 1.22+ Loop Variable Fix ===")
	values := []string{"a", "b", "c"}
	for _, v := range values {
		go func() {
			// Go 1.22+: v is now a fresh variable per iteration
			// This is safe and prints a, b, c (in any order)
			fmt.Printf("Value: %s\n", v)
		}()
	}

	// Goroutine with parameters
	fmt.Println("\n=== Goroutines with Parameters ===")
	for i := 0; i < 3; i++ {
		go func(num int) {
			fmt.Printf("Goroutine %d is running\n", num)
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Println("\nMain goroutine finished!")
}
