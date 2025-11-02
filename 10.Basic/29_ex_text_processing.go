package main

import (
	"fmt"
	"strings"
	"unicode"
)

// --- Text Processing Library ---

// WordCount counts the occurrences of each word in a given text.
// It returns a map where keys are words and values are their frequencies.
func WordCount(text string) map[string]int {
	// Normalize the text to lower case to make counting case-insensitive.
	text = strings.ToLower(text)
	// Split the text into words based on spaces and punctuation.
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Use a map to store the frequency of each word.
	freqs := make(map[string]int)
	for _, word := range words {
		freqs[word]++
	}
	return freqs
}

// CharFrequency counts the occurrences of each character (rune) in a text.
func CharFrequency(text string) map[rune]int {
	freqs := make(map[rune]int)
	for _, r := range text {
		freqs[r]++
	}
	return freqs
}

// FindMinMaxWordLength finds the shortest and longest word lengths in a text.
// This function demonstrates the use of the Go 1.21+ `min` and `max` built-ins.
func FindMinMaxWordLength(text string) (minLength, maxLength int) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0, 0
	}

	// Initialize minLength to a very large value to ensure the first word length is smaller.
	minLength = 1e9 // A large number
	maxLength = 0

	for _, word := range words {
		length := len(word)
		minLength = min(minLength, length) // Use the `min` built-in
		maxLength = max(maxLength, length) // Use the `max` built-in
	}
	return minLength, maxLength
}

// TransformText provides a simple example of text transformation.
// It takes a pointer to a string and modifies it in place to be uppercase.
func TransformTextToUpper(text *string) {
	// Dereference the pointer to get the value, transform it,
	// and then use the pointer to update the original variable.
	*text = strings.ToUpper(*text)
}

func main() {
	text := "Hello world! This is a simple, simple world of Go."

	fmt.Printf("--- Original Text ---\n\"%s\"\n\n", text)

	// --- 1. Word Counting ---
	fmt.Println("--- Word Count ---")
	wordCounts := WordCount(text)
	for word, count := range wordCounts {
		fmt.Printf("'%s': %d\n", word, count)
	}
	fmt.Println()

	// --- 2. Character Frequency ---
	fmt.Println("--- Character Frequency ---")
	charCounts := CharFrequency(text)
	// Example of clearing a map (Go 1.21+)
	fmt.Printf("Char frequency map has %d unique characters.\n", len(charCounts))
	clear(charCounts)
	fmt.Printf("After `clear`, map has %d unique characters.\n\n", len(charCounts))

	// --- 3. Min/Max Word Length ---
	fmt.Println("--- Min/Max Word Length ---")
	minLen, maxLen := FindMinMaxWordLength(text)
	fmt.Printf("Shortest word length: %d, Longest word length: %d\n\n", minLen, maxLen)

	// --- 4. Text Transformation (using a pointer) ---
	fmt.Println("--- Text Transformation ---")
	fmt.Printf("Original text variable before transformation: \"%s\"\n", text)
	// We pass the memory address of the `text` variable.
	TransformTextToUpper(&text)
	fmt.Printf("Original text variable after transformation: \"%s\"\n", text)
}
