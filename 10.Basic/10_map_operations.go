package main

import "fmt"

func main() {
	// --- Declaring and Initializing a Map ---
	// Maps must be initialized before use.
	ages := make(map[string]int)

	// --- 1. Insertion (or Update) ---
	// Use `map[key] = value` syntax to add a new key-value pair or update an existing one.
	fmt.Println("--- 1. Insertion / Update ---")
	ages["alice"] = 31
	ages["bob"] = 34
	ages["charlie"] = 28
	fmt.Println("Map after insertions:", ages)

	// Update an existing value
	ages["alice"] = 32
	fmt.Println("Map after updating Alice's age:", ages)
	fmt.Println()

	// --- 2. Deletion ---
	// The `delete` built-in function removes an element from the map by key.
	fmt.Println("--- 2. Deletion ---")
	fmt.Println("Deleting Bob...")
	delete(ages, "bob")
	fmt.Println("Map after deletion:", ages)

	// Deleting a key that is not present does not cause an error.
	delete(ages, "david") // No error
	fmt.Println("Map after attempting to delete a non-existent key:", ages)
	fmt.Println()

	// --- 3. Retrieval and Checking for Existence ---
	// Access a value with `map[key]`.
	fmt.Println("--- 3. Retrieval ---")
	fmt.Println("Alice's age is:", ages["alice"])

	// If you retrieve a key that doesn't exist, you get the zero value for the value type.
	fmt.Println("David's age (non-existent) is:", ages["david"]) // Prints 0, the zero value for int

	// To distinguish a zero value from a non-existent key, use the two-value assignment.
	age, ok := ages["charlie"]
	if ok {
		fmt.Println("Charlie's age is", age)
	} else {
		fmt.Println("Charlie is not in the map.")
	}

	age, ok = ages["bob"]
	if ok {
		fmt.Println("Bob's age is", age)
	} else {
		fmt.Println("Bob is not in the map (as expected).")
	}
	fmt.Println()

	// --- 4. Iteration ---
	// Use a `for...range` loop to iterate over the keys and values of a map.
	// The order of iteration is not guaranteed.
	fmt.Println("--- 4. Iteration ---")
	ages["david"] = 45 // Add one more for iteration
	for name, age := range ages {
		fmt.Printf("%s is %d years old.\n", name, age)
	}
}
