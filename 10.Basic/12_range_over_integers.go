package main

import "fmt"

func main() {
	// This example demonstrates the "range over integers" feature introduced in Go 1.22.
	// It provides a more concise way to write a loop that executes a fixed number of times.
	fmt.Println("--- Detailed Look at `for...range` over an integer ---")

	fmt.Println("Starting countdown...")
	for i := range 10 {
		fmt.Println(10 - i)
	}
	fmt.Println("go1.22 has lift-off!")
	fmt.Println()

	// This is functionally equivalent to the more traditional `for` loop:
	fmt.Println("--- Same countdown using a traditional `for` loop for comparison ---")
	for i := 0; i < 10; i++ {
		fmt.Println(10 - i)
	}
	fmt.Println("Traditional loop has lift-off!")

	// The `range` version is often considered more readable and less prone to off-by-one errors
	// in the loop condition.
}
