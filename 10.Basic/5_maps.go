package main

import "fmt"

func main() {
	// A map maps keys to values.
	// The zero value of a map is `nil`. A `nil` map has no keys, nor can keys be added.
	var m map[string]int
	fmt.Println("Uninitialized map 'm':", m, "is nil?", m == nil)

	// The `make` function returns an initialized map of a given type.
	m = make(map[string]int)
	fmt.Println("Map 'm' after make:", m, "is nil?", m == nil)

	// Set key/value pairs using typical `name[key] = val` syntax.
	m["k1"] = 7
	m["k2"] = 13
	fmt.Println("Map 'm' after setting values:", m)

	// Get a value for a key with `name[key]`.
	v1 := m["k1"]
	fmt.Println("Value of m[\"k1\"]:", v1)

	// The `len` function returns the number of key/value pairs.
	fmt.Println("Length of 'm':", len(m))

	// The `delete` function removes a key/value pair from a map.
	delete(m, "k2")
	fmt.Println("Map 'm' after deleting 'k2':", m)

	// When getting a value from a map, you can optionally get a second return value
	// which is a boolean indicating if the key was present in the map.
	// This is useful to distinguish between a missing key and a key with a zero value (like 0 or "").
	_, prs := m["k2"]
	fmt.Println("Is 'k2' present in 'm'?", prs)

	val, present := m["k1"]
	fmt.Println("Is 'k1' present in 'm'?", present, ", value is:", val)
	fmt.Println()

	// You can also declare and initialize a new map in one line with a map literal.
	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("Initialized map 'n':", n)

	// Note: Maps are not ordered. When you iterate over a map,
	// the order of iteration is not guaranteed to be the same every time.
	fmt.Println("\nIterating over map 'n':")
	for key, value := range n {
		fmt.Printf("Key: %s -> Value: %d\n", key, value)
	}
}
