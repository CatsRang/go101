package main

import (
	"errors"
	"fmt"
	"strings"
)

// This file contains functions to be tested.
// See 06_testing_example_test.go for the corresponding tests.

// Add returns the sum of two integers
func Add(a, b int) int {
	return a + b
}

// Divide returns the quotient of a/b or error if b is zero
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// IsPalindrome checks if a string is a palindrome (case-insensitive)
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// Fibonacci returns the nth Fibonacci number
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciIterative is an optimized version for benchmarking comparison
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// User represents a user entity
type User struct {
	ID    int
	Name  string
	Email string
}

// Validate checks if the user data is valid
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

// StringSliceContains checks if a slice contains a string
func StringSliceContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("This file contains functions to be tested.")
	fmt.Println("Run tests with: go test -v ./30.Advanced/")
	fmt.Println()

	// Demo the functions
	fmt.Printf("Add(2, 3) = %d\n", Add(2, 3))

	result, err := Divide(10, 2)
	if err == nil {
		fmt.Printf("Divide(10, 2) = %.2f\n", result)
	}

	fmt.Printf("IsPalindrome(\"radar\") = %v\n", IsPalindrome("radar"))
	fmt.Printf("Fibonacci(10) = %d\n", Fibonacci(10))

	user := &User{ID: 1, Name: "John", Email: "john@example.com"}
	if err := user.Validate(); err == nil {
		fmt.Printf("User %s is valid\n", user.Name)
	}
}
