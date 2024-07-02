/*
Generics in the Go programming language are a feature that allows functions, types, and data structures to be written with placeholders for the types they operate on. This enables more flexible and reusable code. Before generics, Go developers had to either write specific functions for each type or use interface types with runtime type assertions or type switches, which were less type-safe and often less performant.

Key Concepts of Generics in Go
Type Parameters: Type parameters are used to define a generic type or function. They are specified using square brackets [] and define the types that can be passed to the generic function or type.

Type Constraints: Type constraints specify what capabilities a type parameter must have, such as certain methods that must be implemented. These are defined using interfaces.

Type Inference: When using generic functions, Go can often infer the specific type from the context, minimizing the need for explicit type declarations.
*/

/**
[References]
- https://go.dev/blog/intro-generics
- https://medium.com/the-godev-corner/go-generics-everything-you-need-to-know-52dd3796d8a1
- https://gfw.go101.org/generics/555-type-constraints-and-parameters.html
*/

package main

import "fmt"
import "golang.org/x/exp/constraints"

func Swap[T any](a, b T) (T, T) {
	return b, a
}

// Swap02 Define a generic function with a type parameter T.
func Swap02[T any, U any](a T, b U) (U, T) {
	return b, a
}

func add[T constraints.Integer | constraints.Float](a, b T) T { // ❶
	return a + b
}

func GMin[T constraints.Ordered](x, y T) T {
	if x > y {
		return y
	} else {
		return x
	}
}

type ExInt interface {
	int | int8 | int16 | int32 | int64
}

// ExInt02 ~T means the set of all types whose underlying type is T
type ExInt02 interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func main() {
	// Example with integers.
	a, b := Swap(1, 2)
	fmt.Println(a, b) // Output: 2 1

	// Example with strings.
	x, y := Swap("hello", "world")
	fmt.Println(x, y) // Output: world hello

	fmt.Println("> add", add(10, 20))
	fmt.Println("> add[int32]", add[int32](10, 20))

	fmt.Println("> GMin", GMin(10, 20))
	fmt.Println("> GMin[int32] ", GMin[int32](32, 11))
}
