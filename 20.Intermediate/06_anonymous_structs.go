package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("=== Basic Anonymous Struct ===")

	// Define and initialize in one statement
	person := struct {
		Name string
		Age  int
		City string
	}{
		Name: "John Doe",
		Age:  30,
		City: "San Francisco",
	}

	fmt.Printf("Person: %+v\n", person)
	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)

	fmt.Println("\n=== Anonymous Struct in Function Return ===")

	// Function that returns an anonymous struct
	getConfig := func() struct {
		Host string
		Port int
	} {
		return struct {
			Host string
			Port int
		}{
			Host: "localhost",
			Port: 8080,
		}
	}

	config := getConfig()
	fmt.Printf("Config: %+v\n", config)

	fmt.Println("\n=== Anonymous Struct for JSON Response ===")

	// Common pattern for API responses
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}{
		Status:  "success",
		Message: "User created successfully",
		Data: struct {
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			UserID:   123,
			Username: "johndoe",
			Email:    "john@example.com",
		},
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("JSON Response:\n%s\n", string(jsonData))

	fmt.Println("\n=== Slice of Anonymous Structs ===")

	// Common for test data or temporary collections
	users := []struct {
		ID   int
		Name string
		Role string
	}{
		{1, "Alice", "Admin"},
		{2, "Bob", "User"},
		{3, "Charlie", "Moderator"},
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("  ID: %d, Name: %s, Role: %s\n", user.ID, user.Name, user.Role)
	}

	fmt.Println("\n=== Map with Anonymous Struct Values ===")

	// Anonymous struct as map value type
	inventory := map[string]struct {
		Quantity int
		Price    float64
		Location string
	}{
		"laptop": {
			Quantity: 15,
			Price:    999.99,
			Location: "Warehouse A",
		},
		"mouse": {
			Quantity: 50,
			Price:    29.99,
			Location: "Warehouse B",
		},
		"keyboard": {
			Quantity: 30,
			Price:    79.99,
			Location: "Warehouse A",
		},
	}

	fmt.Println("Inventory:")
	for item, details := range inventory {
		fmt.Printf("  %s: Qty=%d, Price=$%.2f, Location=%s\n",
			item, details.Quantity, details.Price, details.Location)
	}

	fmt.Println("\n=== Anonymous Struct for Table-Driven Tests ===")

	// Common pattern in Go testing
	testCases := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero", 0, 0},
		{"positive", 5, 25},
		{"negative", -3, 9},
		{"large", 10, 100},
	}

	square := func(x int) int { return x * x }

	fmt.Println("Test cases for square function:")
	for _, tc := range testCases {
		result := square(tc.input)
		status := "PASS"
		if result != tc.expected {
			status = "FAIL"
		}
		fmt.Printf("  [%s] %s: square(%d) = %d (expected %d)\n",
			status, tc.name, tc.input, result, tc.expected)
	}

	fmt.Println("\n=== Anonymous Struct for Grouping Related Data ===")

	// Temporary grouping without defining a named type
	result := struct {
		Sum     int
		Average float64
		Min     int
		Max     int
		Count   int
	}{
		Sum:     150,
		Average: 30.0,
		Min:     10,
		Max:     50,
		Count:   5,
	}

	fmt.Printf("Statistics: %+v\n", result)

	fmt.Println("\n=== Anonymous Struct in Channel ===")

	// Channel of anonymous structs
	events := make(chan struct {
		Type      string
		Timestamp string
		Data      interface{}
	}, 3)

	// Send events
	events <- struct {
		Type      string
		Timestamp string
		Data      interface{}
	}{
		Type:      "user.created",
		Timestamp: "2024-01-15T10:30:00Z",
		Data:      map[string]interface{}{"user_id": 123},
	}

	events <- struct {
		Type      string
		Timestamp string
		Data      interface{}
	}{
		Type:      "user.updated",
		Timestamp: "2024-01-15T10:35:00Z",
		Data:      map[string]interface{}{"user_id": 123, "field": "email"},
	}

	close(events)

	fmt.Println("Events:")
	for event := range events {
		fmt.Printf("  Type: %s, Time: %s, Data: %v\n",
			event.Type, event.Timestamp, event.Data)
	}

	fmt.Println("\n=== Anonymous Struct for Options ===")

	// Configuration or options pattern
	processData := func(data []int, opts struct {
		Reverse   bool
		Filter    func(int) bool
		Transform func(int) int
	}) []int {
		result := make([]int, 0)

		for _, v := range data {
			// Apply filter if provided
			if opts.Filter != nil && !opts.Filter(v) {
				continue
			}

			// Apply transformation if provided
			value := v
			if opts.Transform != nil {
				value = opts.Transform(v)
			}

			result = append(result, value)
		}

		// Reverse if requested
		if opts.Reverse {
			for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
				result[i], result[j] = result[j], result[i]
			}
		}

		return result
	}

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	processed := processData(numbers, struct {
		Reverse   bool
		Filter    func(int) bool
		Transform func(int) int
	}{
		Reverse:   true,
		Filter:    func(n int) bool { return n%2 == 0 }, // Even numbers only
		Transform: func(n int) int { return n * n },      // Square them
	})

	fmt.Printf("Original: %v\n", numbers)
	fmt.Printf("Processed (even, squared, reversed): %v\n", processed)

	fmt.Println("\n=== Anonymous Struct Embedding ===")

	// Combining anonymous struct with regular fields
	employee := struct {
		ID         int
		Department string
		Person     struct {
			Name  string
			Email string
		}
		Salary float64
	}{
		ID:         1001,
		Department: "Engineering",
		Person: struct {
			Name  string
			Email string
		}{
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
		Salary: 85000.00,
	}

	fmt.Printf("Employee: %+v\n", employee)
	fmt.Printf("Name: %s, Department: %s\n", employee.Person.Name, employee.Department)

	fmt.Println("\n=== When to Use Anonymous Structs ===")
	fmt.Println("Good use cases:")
	fmt.Println("  - One-off data structures")
	fmt.Println("  - JSON encoding/decoding for API responses")
	fmt.Println("  - Test tables and test data")
	fmt.Println("  - Grouping function parameters or return values")
	fmt.Println("  - Temporary data grouping")
	fmt.Println("")
	fmt.Println("Avoid anonymous structs when:")
	fmt.Println("  - The structure is used in multiple places")
	fmt.Println("  - You need to attach methods to the type")
	fmt.Println("  - The type has clear business meaning (use named type)")
	fmt.Println("  - Code readability would benefit from a named type")
}
