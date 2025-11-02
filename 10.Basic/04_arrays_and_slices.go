package main

import "fmt"

func main() {
	// --- ARRAYS ---
	// An array has a fixed size. The type `[n]T` is an array of `n` values of type `T`.
	fmt.Println("--- ARRAYS ---")
	var a [5]int // Declares an array of 5 integers, initialized to zero.
	fmt.Println("Empty array 'a':", a)

	a[4] = 100 // Set a value at an index.
	fmt.Println("Set a[4]=100:", a)
	fmt.Println("Get a[4]:", a[4])
	fmt.Println("Length of 'a':", len(a))

	// You can declare and initialize an array in one line.
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("Initialized array 'b':", b)
	fmt.Println()

	// --- SLICES ---
	// A slice is a dynamically-sized, flexible view into the elements of an array.
	// In practice, slices are much more common than arrays.
	fmt.Println("--- SLICES ---")
	var s []int // A slice of integers. Unlike an array, the length is not specified.
	fmt.Println("Uninitialized slice 's':", s, "is nil?", s == nil)

	// The `make` built-in function allocates a zeroed array and returns a slice that refers to that array.
	s = make([]int, 5) // Creates a slice of length 5 (and capacity 5).
	fmt.Println("Slice 's' after make:", s)
	s[4] = 100
	fmt.Println("Set s[4]=100:", s)

	// Slices can be created with the `[]T{}` literal syntax.
	t := []int{10, 20, 30, 40, 50}
	fmt.Println("Initialized slice 't':", t)

	// Slices can be "sliced" themselves using the `[low:high]` syntax.
	// This selects a half-open range which includes the first element, but excludes the last one.
	sliceOfT := t[1:4] // Includes elements t[1], t[2], t[3]
	fmt.Println("Slice of 't' [1:4]:", sliceOfT)

	sliceOfT[0] = 999 // IMPORTANT: Modifying the slice modifies the underlying array!
	fmt.Println("Modified sliceOfT[0]:", sliceOfT)
	fmt.Println("Original slice 't' is now:", t, "<- See the change!")

	// The `append` function is used to add new elements to a slice.
	// It may create a new, larger underlying array if the original is not large enough.
	fmt.Println("\n--- Appending to a slice ---")
	var u []int
	fmt.Println("Initial slice 'u':", u)
	u = append(u, 1)
	fmt.Println("After append(u, 1):", u)
	u = append(u, 2, 3, 4)
	fmt.Println("After append(u, 2, 3, 4):", u)
}
