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
	/*
		Map is a built-in data type that associates key and values.
		It's an unordered collection where you can store and retrieves values by keys.
		Maps are useful when you need to look up data quickly by a custome index.
	*/
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
	/*
		A pointer holds the memory address of a value.
		Pointers are used to access and modify the values of variables indirectly.
		They are particulary useful when you want to avoid copying large structs
		or when you need to modify the original variable.
	*/
	{
		v1 := 10
		p1 := &v1

		fmt.Println("> ", v1, p1, *p1)

		*p1 = 20
		fmt.Println("> ", v1, p1, *p1)
	}

	// ---- Functions
	/*
		Functions in Go can also be treated as values.
		You can assign them to variables, pass them as arguments to other functions, and even return them from functions.
		This makes Go functions highly versatile, allowing for functional programming patterns like closures.
	*/
	{
		fcompute := func(fn func(float64, float64) float64, a float64, b float64) float64 {
			return fn(a, b)
		}

		fsquare := func(x, y float64) float64 {
			return x*x + y*y
		}

		fmt.Println("> ", fcompute(fsquare, 2, 3))
	}

}
