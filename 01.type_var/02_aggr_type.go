package main

import (
	"fmt"
	"reflect"
)

func main() {
	// ---- Arrays
	/*
		An array is a numbered sequence of elements of a single type,
		called the element type. The number of elements is a part of the array's type,
		which means arrays are of fixed size and cannot be resized. Arrays are useful
		when you know the exact number of elements you need to store ahead of time.
	*/

	var aint [5]int32 = [5]int32{1, 2, 3, 4, 5}
	var aint2 []int = []int{1, 2, 3, 4, 5, 6} // slice, not an array
	aint2 = append(aint2, 7)
	fmt.Printf("> %v, %T, %v,  %T, %v, %v\n", aint, aint, aint2, aint2, len(aint2), reflect.TypeOf(aint2))

	aint3 := [...]int32{1, 2, 3, 4, 5} // array creation
	// append(aing3, 6) // error
	fmt.Printf("> %v, %T\n", aint3, aint3)

	// array size is fixed, and it's length is part of the type
	// for example, [4]int and [5]int are distinct, incompatible types
	//aint = [6]int32{1, 2, 3, 4, 5, 6} // error: type mismatch

	// array do not have to be initialized with a value
	var aint4 [4]int32
	fmt.Printf("> %v, %v\n", aint4, len(aint4))

	// ---- Structs

	/*
		A struct is a composite type that groups together variables under a single name.
		These variables, known as fields, can be of different types.
		Structs are used to create complex data structures that represent objects with properties.
		They are especially useful for grouping data together to form records.
	*/
	type Person struct {
		Name string
		Age  int
	}

	var p1 Person
	p1.Name = "길동"
	p1.Age = 30

	fmt.Println(">", p1)

	p2 := Person{Name: "길동2", Age: 31}
	fmt.Println(">", p2)

	p3 := Person{"길동3", 32}
	fmt.Println(">", p3)
}
