package main

import "fmt"

func main() {
	// ---- Slices
	{
		// Creating a slice from an array
		//primes := [6]int{2, 3, 5, 7, 11, 13} // a array
		primes := [...]int{2, 3, 5, 7, 11, 13} // a array
		var s []int = primes[1:4]              // Slices from index 1 to 3 (excluding index 4)
		fmt.Printf("> %T %v\n", s, s)

		// Short declaration syntax for slices
		numbers := []int{1, 2, 3, 4, 5}
		fmt.Printf("> %T %v\n", numbers, numbers)
	}

	// ---- Maps
	{
		// Creating a map with string keys and int values
		var ages map[string]int
		ages = make(map[string]int) // Initialize the map before using it

		ages["Alice"] = 30
		ages["Bob"] = 25

		fmt.Println(ages)

		// Short declaration syntax
		colors := map[string]string{"red": "#ff0000", "green": "#00ff00", "blue": "#0000ff"}
		fmt.Println(colors)

		// check if key exists
		c, exists := colors["red"]
		if exists {
			fmt.Printf("> 'red' color found, %v\n", c)
		}

		// iterate
		for k, v := range colors {
			fmt.Printf("> %v:%v\n", k, v)
		}

		// delete
		delete(colors, "blue")
		fmt.Printf("> %v\n", colors)
	}

	// ---- Pointers

	// ---- Functions

}
