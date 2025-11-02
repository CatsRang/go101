package main

import "fmt"

// --- Function with Named Return Values ---

// In this function, the return values are named `sum` and `product`.
// They are declared in the function signature and can be used like local variables.
func calculateNamed(a, b int) (sum int, product int) {
	// `sum` and `product` are initialized to their zero values (0 for int).
	sum = a + b
	product = a * b

	// A `return` statement without arguments is known as a "naked" return.
	// It returns the current values of the named return variables.
	return
}

// Another example for clarity.
func split(input int) (x, y int) {
	x = input * 4 / 9
	y = input - x
	// Naked return. Returns the current values of x and y.
	return
}

func main() {
	// --- Using Named Return Values ---

	// The function is called in the same way.
	s, p := calculateNamed(7, 8)
	fmt.Printf("For 7 and 8 -> Sum: %d, Product: %d\n\n", s, p)

	// Another example.
	fmt.Println("Splitting 17:")
	xVal, yVal := split(17)
	fmt.Printf("x = %d, y = %d\n", xVal, yVal)

	// While "naked" returns are allowed, they can sometimes make code less clear,
	// especially in longer functions. It's often better to be explicit about what
	// you are returning. For example: `return sum, product`.
	// However, for short functions like these, they can be useful.
}
