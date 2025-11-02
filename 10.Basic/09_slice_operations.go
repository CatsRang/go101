package main

import "fmt"

func main() {
	// --- 1. Slicing ---
	// Slicing creates a new slice header that points to the same underlying array.
	fmt.Println("--- 1. Slicing ---")
	original := []string{"a", "b", "c", "d", "e"}
	fmt.Printf("Original: %v\n", original)

	// Create a slice containing elements from index 1 up to (not including) 4
	slice1 := original[1:4] // contains "b", "c", "d"
	fmt.Printf("Slice [1:4]: %v\n", slice1)

	// Omitting the low bound implies 0
	slice2 := original[:3] // contains "a", "b", "c"
	fmt.Printf("Slice [:3]: %v\n", slice2)

	// Omitting the high bound implies len(s)
	slice3 := original[2:] // contains "c", "d", "e"
	fmt.Printf("Slice [2:]: %v\n\n", slice3)

	// --- 2. Append ---
	// The `append` function adds elements to the end of a slice.
	fmt.Println("--- 2. Append ---")
	var s []int
	s = append(s, 1)
	s = append(s, 2, 3)
	fmt.Printf("Appended single elements: %v\n", s)

	// You can also append another slice using the `...` syntax to expand it.
	s2 := []int{4, 5, 6}
	s = append(s, s2...)
	fmt.Printf("Appended another slice: %v\n\n", s)

	// --- 3. Copy ---
	// The `copy` function copies elements from a source slice into a destination slice.
	// It's important because it creates a true copy, not just a view of the same underlying array.
	// The number of elements copied is the minimum of len(src) and len(dst).
	fmt.Println("--- 3. Copy ---")
	src := []int{10, 20, 30}
	// The destination slice must be created first.
	dst := make([]int, len(src))

	numCopied := copy(dst, src)

	fmt.Printf("Source:      %v\n", src)
	fmt.Printf("Destination: %v\n", dst)
	fmt.Printf("Number of elements copied: %d\n", numCopied)

	// Let's prove they are independent by modifying the source.
	src[0] = 999
	fmt.Printf("Modified source: %v\n", src)
	fmt.Printf("Destination remains unchanged: %v\n", dst)
}
