package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Strings in Go are immutable sequences of bytes.
	// They are typically interpreted as UTF-8 encoded text.
	s := "Hello, Go World!"
	fmt.Printf("Original string: '%s'\n\n", s)

	// --- Basic String Operations ---
	fmt.Println("--- Basic Operations ---")
	fmt.Println("Length (bytes):", len(s))
	// To get the number of characters (runes), use the `strings` package
	fmt.Println("Length (runes/characters):", strings.Count(s, "")-1) // A common trick for rune count
	fmt.Println("Character at index 4 (byte):", string(s[4]))
	fmt.Println()

	// --- Common functions from the `strings` package ---
	fmt.Println("--- `strings` Package Examples ---")
	fmt.Println("Contains 'Go':", strings.Contains(s, "Go"))
	fmt.Println("Count of 'o':", strings.Count(s, "o"))
	fmt.Println("Has prefix 'Hello':", strings.HasPrefix(s, "Hello"))
	fmt.Println("Has suffix 'World!':", strings.HasSuffix(s, "World!"))
	fmt.Println("Index of 'Go':", strings.Index(s, "Go"))
	fmt.Println("Join ['a','b','c'] with '-':", strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Println("Repeat 'a' 5 times:", strings.Repeat("a", 5))
	fmt.Println("Replace 'o' with '0' (first one):", strings.Replace(s, "o", "0", 1))
	fmt.Println("Replace 'o' with '0' (all):", strings.Replace(s, "o", "0", -1))
	fmt.Println("Split on ', ':", strings.Split(s, ", "))
	fmt.Println("To lower case:", strings.ToLower(s))
	fmt.Println("To upper case:", strings.ToUpper(s))
	fmt.Println()

	// --- String Concatenation ---
	s1 := "First"
	s2 := "Second"
	concat := s1 + " " + s2
	fmt.Println("Concatenated string:", concat)
	fmt.Println()

	// --- String Conversion with `strconv` ---
	fmt.Println("--- `strconv` Package Examples ---")
	// Convert an integer to a string
	num := 123
	strNum := strconv.Itoa(num)
	fmt.Printf("Integer %d converted to string '%s'\n", num, strNum)

	// Convert a string to an integer
	strInt := "456"
	numInt, err := strconv.Atoi(strInt)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	} else {
		fmt.Printf("String '%s' converted to integer %d\n", strInt, numInt)
	}

	// Parsing other types
	strBool := "true"
	b, _ := strconv.ParseBool(strBool)
	fmt.Printf("Parsed bool: %v\n", b)

	strFloat := "3.14159"
	f, _ := strconv.ParseFloat(strFloat, 64)
	fmt.Printf("Parsed float: %f\n", f)
}
