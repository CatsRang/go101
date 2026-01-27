package main

import (
	"fmt"
	"strings"
)

/*
The Ellipsis (...) Notation in Go

The three-dot (...) notation is called the "ellipsis" and is used in three distinct contexts:

1. VARIADIC FUNCTION PARAMETERS: ...type
   - Allows a function to accept zero or more arguments of a type
   - Syntax: func name(param ...type)
   - Inside the function, param is a slice of that type

2. EXPANDING SLICES TO ARGUMENTS: varname...
   - Unpacks a slice into individual arguments when calling a variadic function
   - Syntax: functionName(sliceVar...)
   - Only allowed when passing to variadic parameters

3. ARRAY LITERALS WITH INFERRED LENGTH: [...]type{}
   - Allows Go to infer array length from the number of elements
   - Syntax: arr := [...]type{elem1, elem2, ...}
   - Used at definition time, not in function signatures

Summary Table:
┌──────────────┬────────────────────┬─────────────────────────────────┬──────────────────────┐
│ Notation     │ Context            │ Meaning/Usage                   │ Example              │
├──────────────┼────────────────────┼─────────────────────────────────┼──────────────────────┤
│ ...type      │ Function parameter │ Accept zero or more of type     │ func f(x ...int)     │
│ varname...   │ Function argument  │ Expand slice as individual args │ f(nums...)           │
│ [...]type{}  │ Array literal      │ Infer array length from elems   │ arr := [...]int{1,2} │
└──────────────┴────────────────────┴─────────────────────────────────┴──────────────────────┘
*/

// --- Example 1: Variadic function with a specific type (int) ---
// USE CASE 1: ...type in function parameter
// sumAll calculates the total of an arbitrary number of integers.
// Inside the function, numbers is a slice ([]int), even when called with individual arguments.
func sumAll(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// --- Example 2: Variadic function with another specific type (string) ---
// USE CASE 1: ...type in function parameter
// processItems joins a variable number of strings with a prefix.
func processItems(prefix string, items ...string) {
	processedString := strings.Join(items, ", ")
	fmt.Printf("[%s]: %s\n", prefix, processedString)
}

// --- Example 3: Variadic function with any type (`...any`) ---
// USE CASE 1: ...type in function parameter
// log prints a formatted message with a variable number of arguments of any type.
func log(prefix string, values ...any) {
	// Safely convert each value to its string representation for display.
	stringValues := make([]string, len(values))
	for i, v := range values {
		stringValues[i] = fmt.Sprint(v)
	}
	fmt.Printf("[%s] %s\n", prefix, strings.Join(stringValues, " "))
}

func main() {
	fmt.Println("=== Understanding Variadic Functions and Ellipsis (...) ===")
	fmt.Println()

	// ========================================================================
	// USE CASE 1: Variadic Function Parameters (...type)
	// ========================================================================
	fmt.Println("--- USE CASE 1: Variadic Function Parameters (...type) ---")

	// Call with individual int arguments
	fmt.Println("Sum 1:", sumAll(1, 2, 3)) // Output: 6

	// Call with zero arguments (valid for variadic functions)
	fmt.Println("Sum 2:", sumAll()) // Output: 0

	// Call with many arguments
	fmt.Println("Sum 3:", sumAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)) // Output: 55

	fmt.Println()

	// ========================================================================
	// USE CASE 2: Expanding Slices to Arguments (varname...)
	// ========================================================================
	fmt.Println("--- USE CASE 2: Expanding Slices (varname...) ---")

	// Instead of calling with individual arguments, we can expand a slice
	nums := []int{10, 20, 30}
	fmt.Println("Sum 4 (from slice):", sumAll(nums...)) // Expands nums as separate arguments
	// This is equivalent to: sumAll(10, 20, 30)

	// Example with strings
	vegetables := []string{"carrot", "potato", "onion"}
	processItems("VEGETABLES", vegetables...) // Expands the slice

	// Example with any type
	mixedValues := []any{"error code", 500, true}
	log("ERROR", mixedValues...) // Expands the slice

	fmt.Println()

	// ========================================================================
	// USE CASE 3: Array Literals with Inferred Length ([...]type{})
	// ========================================================================
	fmt.Println("--- USE CASE 3: Array Literals with Inferred Length ([...]type{}) ---")

	// Array with inferred length - Go counts the elements
	fruits := [...]string{"apple", "orange", "banana"} // Length is 3
	fmt.Printf("fruits array: %v, Type: %T, Length: %d\n", fruits, fruits, len(fruits))

	// Compare with explicit length
	colors := [3]string{"red", "green", "blue"}
	fmt.Printf("colors array: %v, Type: %T, Length: %d\n", colors, colors, len(colors))

	// Key difference: Arrays vs Slices
	// [...]int{1,2,3} creates an ARRAY with fixed length 3
	// []int{1,2,3} creates a SLICE with length 3 but can grow
	arrayExample := [...]int{1, 2, 3}
	sliceExample := []int{1, 2, 3}
	fmt.Printf("Array type: %T (fixed size)\n", arrayExample)
	fmt.Printf("Slice type: %T (can grow)\n", sliceExample)

	fmt.Println()

	// ========================================================================
	// COMBINING ALL THREE USE CASES
	// ========================================================================
	fmt.Println("--- Combining All Three Use Cases ---")

	// USE CASE 3: Create array with inferred length
	temperatures := [...]int{72, 68, 75, 70, 73}
	fmt.Printf("temperatures: %v\n", temperatures)

	// USE CASE 2: Expand array (converted to slice) to variadic function
	// Note: arrays can be converted to slices using [:]
	totalTemp := sumAll(temperatures[:]...)
	fmt.Printf("Total temperature: %d\n", totalTemp)

	// All three use cases in action:
	// 1. sumAll uses ...int parameter (USE CASE 1)
	// 2. temperatures[:] converts array to slice
	// 3. temperatures[:]... expands slice to arguments (USE CASE 2)
	// 4. temperatures was created with [...] (USE CASE 3)

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. ...type in parameters: variadic functions (accept 0+ args)")
	fmt.Println("2. varname... in calls: expand slice to individual arguments")
	fmt.Println("3. [...]type{} in literals: infer array length from elements")
}
