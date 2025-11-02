package main

import "fmt"

func main() {
	// --- Length and Capacity of a Slice ---
	// The length of a slice is the number of elements it contains.
	// The capacity of a slice is the number of elements in the underlying array,
	// counting from the first element in the slice.

	// Let's create a slice with a specified length and capacity.
	// make([]T, length, capacity)
	s := make([]int, 3, 5)
	printSliceInfo("s", s)

	// The length is 3, so we can access s[0], s[1], s[2].
	s[0], s[1], s[2] = 10, 20, 30
	fmt.Println("After assigning values:", s)

	// Accessing s[3] would cause a panic because it's out of the length bounds.
	// e.g., s[3] = 40 // panic: runtime error: index out of range [3] with length 3

	// --- Appending to a slice within its capacity ---
	// When you append elements and the capacity is sufficient, the length increases,
	// but a new underlying array is NOT allocated. The capacity remains the same.
	fmt.Println("\nAppending within capacity...")
	s = append(s, 40)
	printSliceInfo("s", s)
	s = append(s, 50)
	printSliceInfo("s", s)
	fmt.Println("Slice 's' is now full (len == cap).")

	// --- Appending to a slice beyond its capacity ---
	// If you append more elements, Go will allocate a new, larger underlying array.
	// The elements from the old array are copied to the new one.
	// The capacity will typically double (for smaller slices).
	fmt.Println("\nAppending beyond capacity...")
	s = append(s, 60)
	printSliceInfo("s", s)
	fmt.Println("A new, larger underlying array was allocated.")
	fmt.Println()

	// --- Slicing and its effect on capacity ---
	// Slicing a slice creates a new slice that shares the same underlying array.
	fmt.Println("--- Slicing and Capacity ---")
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	printSliceInfo("a", a)

	b := a[2:5] // Slice from index 2 up to (not including) 5
	printSliceInfo("b", b)
	// The capacity of 'b' is calculated from its starting index (2) to the end of the underlying array 'a'.
	// Capacity = len(a) - startIndex = 10 - 2 = 8.

	c := a[5:] // Slice from index 5 to the end
	printSliceInfo("c", c)
	// Capacity = len(a) - startIndex = 10 - 5 = 5.
}

// Helper function to print slice information
func printSliceInfo(name string, slice []int) {
	fmt.Printf("Slice '%s': len=%d, cap=%d, values=%v\n", name, len(slice), cap(slice), slice)
}
