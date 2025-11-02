package main

import "fmt"

func main() {
	// Go 1.21 introduced three new convenient built-in functions: `min`, `max`, and `clear`.

	// --- 1. `min` and `max` Functions ---
	// `min` and `max` compute the smallest or largest value of a fixed number of arguments.
	// They work with any ordered type, like integers, floats, and strings.
	fmt.Println("--- `min` and `max` ---")

	// Example with integers
	minInt := min(100, 42, 73)
	maxInt := max(100, 42, 73)
	fmt.Printf("For integers 100, 42, 73 -> min: %d, max: %d\n", minInt, maxInt)

	// Example with floating-point numbers
	minFloat := min(5.5, 3.14, 9.81)
	maxFloat := max(5.5, 3.14, 9.81)
	fmt.Printf("For floats 5.5, 3.14, 9.81 -> min: %f, max: %f\n", minFloat, maxFloat)

	// Example with strings (comparison is lexical)
	minString := min("apple", "orange", "banana")
	maxString := max("apple", "orange", "banana")
	fmt.Printf("For strings 'apple', 'orange', 'banana' -> min: '%s', max: '%s'\n\n", minString, maxString)

	// --- 2. `clear` Function ---
	// The `clear` function is used to clear the contents of a map or a slice.
	fmt.Println("--- `clear` ---")

	// --- `clear` on a map ---
	// For maps, `clear` deletes all key-value pairs, resulting in an empty map.
	fmt.Println("Clearing a map...")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("Original map: %v, Length: %d\n", m, len(m))
	clear(m)
	fmt.Printf("Map after `clear`: %v, Length: %d\n\n", m, len(m))

	// --- `clear` on a slice ---
	// For slices, `clear` sets all elements in the slice to their zero value.
	// IMPORTANT: It does NOT change the length or capacity of the slice.
	fmt.Println("Clearing a slice...")
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("Original slice: %v, Length: %d, Capacity: %d\n", s, len(s), cap(s))
	clear(s)
	fmt.Printf("Slice after `clear`: %v, Length: %d, Capacity: %d\n", s, len(s), cap(s))
}
