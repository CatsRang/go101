package main

import (
	"fmt"
	"time"
)

func main() {
	// --- if, else if, else ---
	fmt.Println("--- Demonstrating if/else if/else ---")
	number := 10

	if number > 15 {
		fmt.Println("Number is greater than 15")
	} else if number > 5 {
		fmt.Printf("Number (%d) is greater than 5 but not greater than 15\n", number)
	} else {
		fmt.Println("Number is 5 or less")
	}

	// An `if` statement can also include a short statement to execute before the condition.
	// Variables declared by the statement are only in scope until the end of the `if`.
	if n := 9; n < 0 {
		fmt.Println(n, "is negative")
	} else if n < 10 {
		fmt.Println(n, "has 1 digit")
	} else {
		fmt.Println(n, "has multiple digits")
	}
	fmt.Println()

	// --- switch ---
	fmt.Println("--- Demonstrating switch ---")
	day := time.Now().Weekday()

	fmt.Println("When's Saturday?")
	switch time.Saturday {
	case day + 0:
		fmt.Println("Today.")
	case day + 1:
		fmt.Println("Tomorrow.")
	case day + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	// A switch without an expression is an alternate way to express if/else logic.
	// It can also be used to compare types.
	hour := time.Now().Hour()
	fmt.Println("\nWhat time of day is it?")
	switch {
	case hour < 12:
		fmt.Println("Good morning!")
	case hour < 17:
		fmt.Println("Good afternoon!")
	default:
		fmt.Println("Good evening!")
	}
}
