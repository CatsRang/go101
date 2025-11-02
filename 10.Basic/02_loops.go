package main

import "fmt"

func main() {
	// Go's `for` loop can be used in several ways.

	// --- 1. The most basic type, with a single condition (like a `while` loop in other languages) ---
	fmt.Println("--- 1. `for` with a single condition (like `while`) ---")
	i := 1
	for i <= 3 {
		fmt.Println("Current value of i:", i)
		i = i + 1
	}
	fmt.Println()

	// --- 2. The classic `for` loop with initial/condition/after statements ---
	fmt.Println("--- 2. Classic `for` loop ---")
	for j := 7; j <= 9; j++ {
		fmt.Println("Current value of j:", j)
	}
	fmt.Println()

	// --- 3. A `for` loop without a condition will loop forever (an infinite loop) ---
	// An `if` condition with a `break` statement can be used to exit the loop.
	// A `continue` statement can be used to skip to the next iteration.
	fmt.Println("--- 3. Infinite `for` loop with `break` and `continue` ---")
	k := 0
	for {
		k++
		if k%2 == 0 { // Skip even numbers
			continue
		}
		if k > 10 { // Exit the loop
			break
		}
		fmt.Println("Processing odd number:", k)
	}
	fmt.Println("Exited the infinite loop.")
	fmt.Println()

	// --- 4. `for` with a `range` clause (will be covered in more detail later) ---
	// This is used to iterate over elements in a variety of data structures.
	fmt.Println("--- 4. `for` with `range` (preview) ---")
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums { // we don't need the index, so we use the blank identifier `_`
		sum += num
	}
	fmt.Println("Sum of nums:", sum)
}
