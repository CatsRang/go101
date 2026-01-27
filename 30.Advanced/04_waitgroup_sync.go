package main

import (
	"fmt"
	"sync"
	"time"
)

// sync package provides synchronization primitives for goroutines:
// - WaitGroup: Wait for a collection of goroutines to finish
// - Mutex: Mutual exclusion lock for protecting shared state
// - RWMutex: Reader/Writer mutual exclusion lock
// - Once: Run something exactly once

func main() {
	// === WaitGroup ===
	fmt.Println("=== WaitGroup ===")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment counter before starting goroutine

		go func(id int) {
			defer wg.Done() // Decrement counter when goroutine finishes

			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Duration(id*50) * time.Millisecond)
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	wg.Wait() // Block until counter becomes 0
	fmt.Println("All workers completed!")

	// === Go 1.25+ WaitGroup.Go Method ===
	// Note: This feature is available in Go 1.25+
	// var wg2 sync.WaitGroup
	// wg2.Go(func() {
	//     // No need to manually call Add/Done
	//     fmt.Println("Task using wg.Go()")
	// })
	// wg2.Wait()

	// === Mutex ===
	fmt.Println("\n=== Mutex (Protecting Shared State) ===")

	var counter int
	var mu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()         // Acquire lock
			counter++         // Critical section
			mu.Unlock()       // Release lock
		}()
	}

	wg.Wait()
	fmt.Printf("Counter value: %d (should be 100)\n", counter)

	// === RWMutex ===
	fmt.Println("\n=== RWMutex (Multiple Readers, Single Writer) ===")

	var data = map[string]string{"key": "initial"}
	var rwmu sync.RWMutex

	// Writer
	go func() {
		rwmu.Lock()
		defer rwmu.Unlock()
		data["key"] = "updated"
		fmt.Println("Writer: updated data")
	}()

	// Multiple readers can read simultaneously
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond) // Let writer go first

			rwmu.RLock()
			defer rwmu.RUnlock()
			fmt.Printf("Reader %d: data = %s\n", id, data["key"])
		}(i)
	}

	wg.Wait()

	// === sync.Once ===
	fmt.Println("\n=== sync.Once (Execute Once) ===")

	var once sync.Once
	initFunc := func() {
		fmt.Println("Initialization: This runs only once!")
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d calling Once.Do\n", id)
			once.Do(initFunc)
		}(i)
	}

	wg.Wait()

	// === sync.Map (Thread-safe Map) ===
	fmt.Println("\n=== sync.Map (Thread-safe Map) ===")

	var sm sync.Map

	// Store values
	sm.Store("key1", "value1")
	sm.Store("key2", "value2")

	// Load value
	if val, ok := sm.Load("key1"); ok {
		fmt.Printf("Loaded: %v\n", val)
	}

	// LoadOrStore
	actual, loaded := sm.LoadOrStore("key3", "value3")
	fmt.Printf("LoadOrStore key3: value=%v, loaded=%v\n", actual, loaded)

	// Range over map
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true // Continue iteration
	})

	// === sync.Cond (Condition Variable) ===
	fmt.Println("\n=== sync.Cond (Condition Variable) ===")

	var ready bool
	var condMu sync.Mutex
	cond := sync.NewCond(&condMu)

	// Waiter goroutine
	go func() {
		condMu.Lock()
		for !ready {
			cond.Wait() // Releases lock and waits
		}
		fmt.Println("Waiter: condition is ready!")
		condMu.Unlock()
	}()

	time.Sleep(50 * time.Millisecond)

	// Signal the condition
	condMu.Lock()
	ready = true
	cond.Signal() // Wake up one waiter
	condMu.Unlock()

	time.Sleep(50 * time.Millisecond)
	fmt.Println("\nSync primitives examples completed!")
}
