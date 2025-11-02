package main

import "fmt"

// --- Anonymous Functions ---
// An anonymous function is a function that is declared without a name.

// --- Closures ---
// A closure is a function value that references variables from outside its body.
// The function may access and assign to the referenced variables; in this sense
// the function is "bound" to the variables.

// This function `intSequence` returns another function.
// The returned function is a closure because it "closes over" the variable `i`.
func intSequence() func() int {
	i := 0
	// The anonymous function being returned here has access to `i`, which is in its parent's scope.
	return func() int {
		i++
		return i
	}
}

func main() {
	// --- Anonymous Function Example ---
	// You can define and call an anonymous function in one go.
	// This is useful for short, one-off tasks.
	func(msg string) {
		fmt.Println("Message from anonymous function:", msg)
	}("Hello, anonymous!") // The `()` here invokes the function immediately.
	fmt.Println()

	// You can also assign an anonymous function to a variable.
	add := func(a, b int) int {
		return a + b
	}
	// And then call it like a regular function.
	sum := add(3, 4)
	fmt.Printf("Sum from assigned anonymous function: %d\n\n", sum)

	// --- Closure Example ---
	// Here, we call `intSequence`, and we get a function back.
	// This function (`nextInt`) holds its own reference to its own `i`.
	nextInt := intSequence()

	// Each time we call `nextInt`, it updates its own `i` variable.
	fmt.Println("--- Closure Example `nextInt` ---")
	fmt.Println(nextInt()) // Output: 1
	fmt.Println(nextInt()) // Output: 2
	fmt.Println(nextInt()) // Output: 3
	fmt.Println()

	// To prove that `nextInt` has its own state, we can create another one.
	newInts := intSequence()
	fmt.Println("--- New Closure `newInts` ---")
	fmt.Println(newInts()) // Output: 1 (This starts its own sequence)
	fmt.Println(newInts()) // Output: 2
}
