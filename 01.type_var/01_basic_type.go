package main

import "fmt"

func main() {
	// ---- Integers
	// int, int8, int16, int32(rune), int64
	// uint, uint8(byte), uint16, uint32, uint64
	{
		var n64 int64
		n64 = 200000

		fmt.Printf("> n64 = (%T)%v\n", n64, n64)

		n1 := -100
		fmt.Printf("> n1 = (%T)%v\n", n1, n1)
	}

	// ---- Floating Point Nubmers
	// float32, float64
	{

	}

	// ---- Complex Numbers
	// complex64, complex128
	{

	}

	// ---- Boolean
	// bool

	// ---- String, Rune
	// string
	{
		for i, r := range "Hello" {
			fmt.Printf("[%v] %c (%v)\n", i, r, r)
		}
	}
}
