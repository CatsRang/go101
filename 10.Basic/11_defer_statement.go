package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("--- Simple Defer Example ---")
	simpleDefer()
	fmt.Println()

	fmt.Println("--- Defer for Resource Cleanup ---")
	// A common and important use of defer is to ensure resources are cleaned up.
	// For example, making sure a file is closed after it has been opened.
	if err := createFileAndWrite("my_deferred_file.txt"); err != nil {
		fmt.Println("Error:", err)
	}
	// The file will be closed by the defer statement in the function,
	// whether the function returns normally or with an error.
	fmt.Println()

	fmt.Println("--- Defer Stacking (LIFO Order) ---")
	// When you have multiple defer statements, they are pushed onto a stack.
	// They are executed in Last-In, First-Out (LIFO) order when the function returns.
	stackingDefer()
}

func simpleDefer() {
	// The `defer` statement defers the execution of a function until the surrounding function returns.
	defer fmt.Println("world") // This will be printed second.

	fmt.Println("hello") // This will be printed first.
}

func createFileAndWrite(filename string) error {
	fmt.Printf("Creating file: %s\n", filename)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	// `defer f.Close()` is called right after the successful file creation.
	// This guarantees that `f.Close()` will be called when the `createFileAndWrite` function exits.
	defer f.Close()
	defer fmt.Println("File cleanup deferred.") // This will run before f.Close() due to LIFO

	_, err = f.WriteString("Hello from a deferred write!")
	if err != nil {
		// Even if we return early due to an error, the deferred calls will still execute.
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("Successfully wrote to file.")
	return nil
}

func stackingDefer() {
	fmt.Println("counting")

	for i := 0; i < 4; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("done")
	// When this function returns, the deferred prints will execute in the reverse order
	// they were deferred. Output will be: 3, 2, 1, 0.
}
