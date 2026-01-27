package main

import (
	"fmt"
	"strings"
)

func sumAll(numbers ...int) int {
	sum := 0

	for idx, num := range numbers {
		fmt.Printf("[%d]: %d\n", idx, num)
		sum += num
	}

	return sum
}

func concatStrings(strs ...string) string {
	return strings.Join(strs, "-")
}

func main() {
	fmt.Println(sumAll(1, 2, 3))
	fmt.Println(concatStrings("hello", "world", "this is", "go"))

	// ---- other eclipse notation
	{
		// array length infer
		arr_1 := [...]int{1, 2, 3}
		fmt.Printf("> %T, %v\n", arr_1, arr_1)
	}
}
