package main

import "fmt"

func main() {

	/*
		- Execution Order:
			Deferred function calls are pushed onto a stack.
			When the surrounding function returns, the deferred calls are executed in last-in, first-out (LIFO) order.
		- Arguments Evaluation:
			The arguments of a deferred function are evaluated at the time the defer statement is executed,
			not when the deferred function is actually invoked.
	*/

	{
		x := 10
		defer fmt.Println("Deferred value of x:", x)

		x = 20
		fmt.Println("Current value of x:", x)
	}

	{
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()

		fmt.Println("Before panic")
		panic("Something went wrong!")
		fmt.Println("This will never be printed")
	}

	/*
		- Order of Execution:
			Deferred functions are executed in LIFO order (last in, first out).
		- Immediate Argument Evaluation:
			The arguments of a deferred function are evaluated when the defer statement is encountered, not when the function is actually executed.
		- Panic Handling:
			defer is often used with recover to handle panics gracefully.
	*/
}
