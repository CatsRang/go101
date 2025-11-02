package main

import "fmt"

func main() {
	// The `for...range` construct is a powerful way to iterate over elements in Go.
	// It works on slices, arrays, strings, maps, and channels.

	// --- 1. `range` on a slice/array ---
	// When ranging over a slice or array, `range` returns two values: the index and a copy of the element at that index.
	fmt.Println("--- `range` on a slice ---")
	nums := []int{2, 3, 4}
	sum := 0
	for index, num := range nums {
		fmt.Printf("Index: %d, Value: %d\n", index, num)
		sum += num
	}
	fmt.Println("Sum of nums:", sum)
	fmt.Println()

	// If you only need the value, you can ignore the index with the blank identifier `_`.
	fmt.Println("--- `range` on a slice (value only) ---")
	sum = 0 // reset sum
	for _, num := range nums {
		sum += num
	}
	fmt.Println("Sum calculated again:", sum)
	fmt.Println()

	// If you only need the index, you can omit the second variable.
	fmt.Println("--- `range` on a slice (index only) ---")
	for index := range nums {
		fmt.Printf("Found index: %d\n", index)
	}
	fmt.Println()

	// --- 2. `range` on a map ---
	// When ranging over a map, `range` returns the key and the value of each pair.
	// The order of iteration over a map is not guaranteed.
	fmt.Println("--- `range` on a map ---")
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	// You can also just iterate over the keys.
	fmt.Println("\n--- `range` on a map (keys only) ---")
	for k := range kvs {
		fmt.Println("key:", k)
	}
	fmt.Println()

	// --- 3. `range` on a string ---
	// When ranging over a string, `range` iterates over Unicode code points (runes).
	// The first value is the starting byte index of the rune and the second is the rune itself.
	fmt.Println("--- `range` on a string ---")
	str := "Go-lang"
	for index, runeValue := range str {
		fmt.Printf("Index: %d, Rune: %c\n", index, runeValue)
	}
}
