package main

import "fmt"

func main() {
	// --- Regular Variables ---
	// When we declare a variable, the program allocates a piece of memory to store its value.
	// We can think of this piece of memory as having an address.
	i := 42
	fmt.Println("--- Regular Variable ---")
	fmt.Printf("The value of 'i' is: %d\n", i)

	// --- Getting a Memory Address with `&` ---
	// The `&` operator gives you the memory address where the variable's value is stored.
	// This address is the "pointer" to the variable.
	fmt.Printf("The memory address of 'i' is: %p\n", &i) // %p is the format specifier for pointers
	fmt.Println()

	// --- Pointers ---
	// A pointer is a variable that stores the memory address of another variable.

	// --- Declaring a Pointer ---
	// The type `*T` is a pointer to a `T` value. Its zero value is `nil`.
	// `p` is a pointer to an integer.
	var p *int
	fmt.Println("--- Pointer Basics ---")
	fmt.Printf("The value of uninitialized pointer 'p' is: %v\n", p) // A nil pointer

	// We can assign the memory address of `i` to the pointer `p`.
	p = &i
	fmt.Printf("Pointer 'p' now stores the memory address of 'i': %p\n", p)

	// So, `p` now "points to" `i`.
	// The value of the pointer itself is the memory address.
	// The type of `p` is `*int` (pointer to an int).
	fmt.Printf("The type of 'p' is %T\n", p)
	fmt.Println()

	// --- The main idea ---
	// Pointers are useful because they allow you to share data without making copies of it.
	// When you pass a pointer to a function, the function can modify the original variable's value.
	// This is known as "pass-by-reference" (though technically, Go is always pass-by-value;
	// it's the pointer's value (the address) that gets passed).
	// We will see this in the next example on dereferencing.
}
