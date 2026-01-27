package main

import (
	"fmt"
	"time"
)

// The select statement lets a goroutine wait on multiple communication operations.
// It blocks until one of its cases can proceed, then executes that case.

func main() {
	// === Basic Select ===
	fmt.Println("=== Basic Select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "message from channel 1"
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "message from channel 2"
	}()

	// Select waits for the first available message
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}

	// === Select with Timeout ===
	fmt.Println("\n=== Select with Timeout ===")

	slowCh := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond)
		slowCh <- "slow response"
	}()

	select {
	case msg := <-slowCh:
		fmt.Println("Received:", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout! No response received")
	}

	// === Non-blocking Select with Default ===
	fmt.Println("\n=== Non-blocking Select with Default ===")

	nonBlockCh := make(chan int)

	select {
	case val := <-nonBlockCh:
		fmt.Println("Received:", val)
	default:
		fmt.Println("No value available (non-blocking)")
	}

	// === Select for Sending ===
	fmt.Println("\n=== Select for Sending ===")

	sendCh := make(chan int, 1)
	sendCh <- 1 // Buffer is full

	select {
	case sendCh <- 2:
		fmt.Println("Sent value 2")
	default:
		fmt.Println("Channel buffer is full, cannot send")
	}

	// === Select in a Loop (Common Pattern) ===
	fmt.Println("\n=== Select in a Loop (Worker Pattern) ===")

	jobs := make(chan int, 5)
	done := make(chan bool)

	// Worker goroutine
	go func() {
		for {
			select {
			case job := <-jobs:
				fmt.Printf("Processing job: %d\n", job)
			case <-done:
				fmt.Println("Worker shutting down")
				return
			}
		}
	}()

	// Send some jobs
	for i := 1; i <= 3; i++ {
		jobs <- i
	}

	time.Sleep(50 * time.Millisecond)
	done <- true // Signal worker to stop

	// === Multiple Ready Channels (Random Selection) ===
	fmt.Println("\n=== Random Selection when Multiple Ready ===")

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2

	// When multiple cases are ready, select chooses one randomly
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch:
			fmt.Printf("Got: %d\n", v)
		}
	}

	fmt.Println("\nSelect examples completed!")
}
