package main

import "fmt"

// ---- Interface
type Speaker interface {
	Speak(msg string)
}

type Human struct {
	name string
}

func (h Human) Speak(msg string) {
	fmt.Printf("%v> %v\n", h.name, msg)
}

func main() {

	/*
		In Go, an interface type is a set of method signatures,
		Interfaces are used to define the behavior of an object

		- Decoupling
		- Plexibility
		- Polymorphism
	*/
	{
		h := Human{name: "길동"}
		h.Speak("안녕하세요")
	}

	// ---- any, interface{} type
	/*
		The empty interface type(interface{}) has no methods,
		so all types implement it by default.
		It can hold values of any type, simliar to 'void*' in C or 'Object' in Java.
		However, using 'interface{}' sacrifices type safety for flexibility and should be done cautiously.
	*/

	{
		f01 := func(v interface{}) {
			fmt.Printf("> %v, %T\n", v, v)
		}

		f01(10)
		f01("Hello")
	}

	// ---- interface type assertions and Type Switches
	/*
		- Type Assertion: Provide access to an interface value's underlying concrete value.
		- Type Switches: Allow you to test against multiple types withing the same interface
	*/
	{
		var i interface{} = "Hello"

		s := i.(string)
		fmt.Println(">", s)

		// ----
		// var.(type)은 switch문에서만 사용 가능
		i = 0.34
		switch v := i.(type) {
		case int:
			fmt.Println("Int>", v)
		case string:
			fmt.Println("String>", v)
		case Human:
			fmt.Println("Human", v)
		default:
			fmt.Printf("%T> %v\n", v, v)
		}
	}
}
