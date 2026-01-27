package main

import (
	"fmt"
	"time"
)

// Channels are the pipes that connect concurrent goroutines.
// They allow goroutines to communicate and synchronize.

func main() {
	// === Unbuffered Channels ===
	fmt.Println("=== Unbuffered Channels ===")

	// Create an unbuffered channel
	ch := make(chan string)

	go func() {
		ch <- "Hello from goroutine!" // Send blocks until receiver is ready
	}()

	msg := <-ch // Receive blocks until sender sends
	fmt.Println(msg)

	// === Buffered Channels ===
	fmt.Println("\n=== Buffered Channels ===")

	// Create a buffered channel with capacity 3
	bufferedCh := make(chan int, 3)

	// Can send up to 3 values without blocking
	bufferedCh <- 1
	bufferedCh <- 2
	bufferedCh <- 3
	// bufferedCh <- 4 // This would block!

	fmt.Printf("Channel length: %d, capacity: %d\n", len(bufferedCh), cap(bufferedCh))
	fmt.Println("Received:", <-bufferedCh)
	fmt.Println("Received:", <-bufferedCh)
	fmt.Println("Received:", <-bufferedCh)

	// === Channel Directions ===
	fmt.Println("\n=== Channel Directions ===")

	// Channels can be restricted to send-only or receive-only
	dataCh := make(chan int)

	go producer(dataCh) // Send-only
	go consumer(dataCh) // Receive-only

	time.Sleep(100 * time.Millisecond)

	// === Closing Channels ===
	fmt.Println("\n=== Closing Channels ===")

	numbers := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers) // Close when done sending
	}()

	// Range over channel until closed
	for num := range numbers {
		fmt.Printf("Received: %d\n", num)
	}

	// === Check if Channel is Closed ===
	fmt.Println("\n=== Check if Channel is Closed ===")

	testCh := make(chan int, 1)
	testCh <- 42
	close(testCh)

	val, ok := <-testCh
	fmt.Printf("Value: %d, Channel open: %v\n", val, ok)

	val, ok = <-testCh // Reading from closed empty channel
	fmt.Printf("Value: %d, Channel open: %v\n", val, ok)
}

// Send-only channel parameter
func producer(ch chan<- int) {
	ch <- 100
	fmt.Println("Producer sent: 100")
}

// Receive-only channel parameter
func consumer(ch <-chan int) {
	val := <-ch
	fmt.Printf("Consumer received: %d\n", val)
}
