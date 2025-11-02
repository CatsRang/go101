package main

import "fmt"

func main() {
	// In Go 1.22 and later, you can use a `for...range` loop over an integer.
	// The loop will iterate from 0 up to (but not including) the integer value.

	fmt.Println("--- Demonstrating `for...range` over an integer ---")
	fmt.Println("Counting up to 5 (exclusive):")
	for i := range 5 {
		fmt.Printf("Iteration %d\n", i)
	}
	fmt.Println()

	// This is a more concise and often clearer alternative to the classic for loop for simple iterations.
	// For example, the above is equivalent to:
	// for i := 0; i < 5; i++ { ... }

	// A common use case is a countdown.
	fmt.Println("--- Countdown Example ---")
	count := 5
	for i := range count {
		fmt.Println(count - i)
	}
	fmt.Println("Blast off!")
	fmt.Println()

	// Note: The loop variable `i` starts at 0 and increments by 1 in each iteration.
	// The value of the integer itself (e.g., `count`) is not changed.

	// If you don't need the loop variable, you can use the blank identifier `_`.
	fmt.Println("--- Using blank identifier `_` ---")
	fmt.Println("Printing 'Hello' 3 times:")
	for range 3 {
		fmt.Println("Hello")
	}
}
