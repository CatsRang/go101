package main

import "fmt"

type DayOfWeek int

const (
	Sunday DayOfWeek = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	_ // skip number
)

func (w DayOfWeek) String() string {
	return [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}[w-1]
}

func main() {
	// ---- Constants
	/*
		A constant is a value that cannot be changed after it is declared.
		Constants are declared like variables, but with the 'const' keyword.
		Constants can be character, string, boolean, or numeric values.
		Constants cannot be declared using the ':=' syntax.
	*/
	{
		const Pi = 3.14
		const Pi2 float64 = 3.14 // with a specific type
		const World = "세계"
		const Truth = true

		fmt.Println("> ", Pi, World, Truth)
	}

	// ---- Enumerated Constants
	/*
		An enumerated constant is a set of named values called elements or enumerators of the constant type.
		Each element is assigned an integer
	*/

	// enum type
	// enum def
	// enum String function

	fmt.Printf("> %v %T %v %T\n", Monday, Monday, Friday, Friday)
	fmt.Println(">", Wednesday.String())
}
