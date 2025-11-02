package main

import (
	"errors"
	"fmt"
)

// --- Function with Multiple Return Values ---

// 1. This function takes two integers and returns two results: their sum and their product.
// The return types are listed in parentheses: `(int, int)`.
func calculateSumAndProduct(a, b int) (int, int) {
	sum := a + b
	product := a * b
	return sum, product
}

// 2. A very common pattern in Go is for a function to return a result and an error value.
// This function attempts to divide two numbers. It returns the result and an error.
// If the division is successful, the error is `nil`. If not, it returns an error.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// It's idiomatic to return a zero value for the result when an error occurs.
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil // `nil` means no error occurred.
}

func main() {
	// --- Handling Multiple Return Values ---

	// 1. We can assign the returned values to two variables.
	s, p := calculateSumAndProduct(5, 6)
	fmt.Printf("For 5 and 6 -> Sum: %d, Product: %d\n", s, p)
	fmt.Println()

	// 2. Handling a function that can return an error.
	fmt.Println("--- Error Handling Example ---")
	result, err := divide(10.0, 2.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result of 10.0 / 2.0 is:", result)
	}

	// Now let's try the case that causes an error.
	result, err = divide(10.0, 0.0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result of 10.0 / 0.0 is:", result)
	}
	fmt.Println()

	// If you only care about one of the returned values, you can use the
	// blank identifier `_` to ignore the other.
	fmt.Println("--- Using the Blank Identifier `_` ---")
	sumOnly, _ := calculateSumAndProduct(3, 4)
	fmt.Println("Ignoring the product, the sum of 3 and 4 is:", sumOnly)
}
