package main

import "fmt"

// This function takes an integer value.
// It receives a COPY of the original value.
func zeroVal(ival int) {
	ival = 0 // This change only affects the copy, not the original.
}

// This function takes a pointer to an integer.
// It receives a copy of the memory address.
func zeroPtr(iptr *int) {
	// The `*` operator here is the dereferencing operator.
	// It gives us access to the value that the pointer `iptr` points to.
	// By changing `*iptr`, we are changing the original variable's value.
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("Initial value of 'i':", i)
	fmt.Println()

	// --- Reading a value through a pointer ---
	p := &i // p now points to i
	fmt.Println("--- Reading via Pointer ---")
	fmt.Printf("Value of 'i' read directly: %d\n", i)
	fmt.Printf("Value of 'i' read through pointer 'p' (dereferencing): %d\n", *p) // The * operator accesses the value
	fmt.Println()

	// --- Modifying a value through a pointer ---
	fmt.Println("--- Modifying via Pointer ---")
	fmt.Println("Modifying the value at the pointer's address to 21...")
	*p = 21 // This is the same as `i = 21`
	fmt.Printf("The new value of 'i' is: %d\n", i)
	fmt.Println()

	// --- Pointers in Functions ---
	// This is the most common use case for pointers: modifying arguments in a function.
	fmt.Println("--- Pointers in Functions ---")
	num := 42
	fmt.Println("Value of 'num' before calling functions:", num)

	// Call `zeroVal` with the value of `num`.
	// A copy of `num`'s value (42) is passed to the function.
	zeroVal(num)
	fmt.Println("Value of 'num' after `zeroVal`:", num) // The original `num` is unchanged.

	// Call `zeroPtr` with the memory address of `num`.
	// A copy of `num`'s address is passed to the function.
	zeroPtr(&num)
	fmt.Println("Value of 'num' after `zeroPtr`:", num) // The original `num` is changed!
	fmt.Println()

	// You can also see the value of the pointer itself.
	fmt.Println("The memory address of 'num' is:", &num)
}
