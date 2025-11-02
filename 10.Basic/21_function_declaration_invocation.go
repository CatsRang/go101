package main

import "fmt"

// --- Function Declaration ---

// A function is a block of code that performs a specific task.
// The `func` keyword is used to declare a function.

// 1. A simple function with no parameters and no return value.
func sayHello() {
	fmt.Println("Hello, Go Functions!")
}

// 2. A function that takes parameters.
// `name` and `age` are parameters of type `string` and `int`.
func greetUser(name string, age int) {
	fmt.Printf("Greetings, %s! You are %d years old.\n", name, age)
}

// 3. A function that returns a value.
// The type of the return value is specified after the parameter list (`int`).
func add(a int, b int) int {
	// The `return` statement specifies the value to be returned.
	return a + b
}

// Note: When two or more consecutive named function parameters share a type,
// you can omit the type from all but the last one.
// The `add` function could also be written as: `func add(a, b int) int`

func main() {
	fmt.Println("--- Function Invocation ---")

	// --- Function Invocation ---
	// To call (or invoke) a function, you write its name followed by parentheses.

	// 1. Invoking the `sayHello` function.
	sayHello()

	// 2. Invoking the `greetUser` function with arguments.
	// "Alice" and 30 are the arguments passed to the function.
	greetUser("Alice", 30)

	// 3. Invoking the `add` function and storing the returned value in a variable.
	sum := add(5, 7)
	fmt.Printf("The sum of 5 and 7 is: %d\n", sum)

	// You can also use the returned value directly.
	fmt.Printf("The sum of 10 and 3 is: %d\n", add(10, 3))
}
