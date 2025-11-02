package main

import (
	"fmt"
	"strings"
)

// --- Example 1: Variadic function with a specific type (int) ---
// sumAll calculates the total of an arbitrary number of integers.
func sumAll(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// --- Example 2: Variadic function with another specific type (string) ---
// processItems joins a variable number of strings with a prefix.
func processItems(prefix string, items ...string) {
	processedString := strings.Join(items, ", ")
	fmt.Printf("[%s]: %s\n", prefix, processedString)
}

// --- Example 3: Variadic function with any type (`...any`) ---
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
	fmt.Println("--- Integer Variadic Example ---")
	// Call with individual int arguments.
	fmt.Println("Sum 1:", sumAll(1, 2, 3))
	// Call by passing a slice of ints.
	nums := []int{10, 20, 30}
	fmt.Println("Sum 2:", sumAll(nums...))
	fmt.Println()

	fmt.Println("--- String Variadic Example ---")
	// Call with individual string arguments.
	processItems("FRUITS", "apple", "orange")
	// Call by passing a slice of strings.
	vegetables := []string{"carrot", "potato"}
	processItems("VEGETABLES", vegetables...)
	fmt.Println()

	fmt.Println("--- Any Type Variadic Example ---")
	log("INFO", "Application started on port", 8080)
	log("ERROR", "User login failed for ID:", 123, "Success:", false)
}
