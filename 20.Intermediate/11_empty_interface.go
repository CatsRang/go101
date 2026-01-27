package main

import (
	"fmt"
)

// ============================================================================
// EMPTY INTERFACE: interface{} and any
// ============================================================================
// The empty interface can hold values of any type because every type
// implements at least zero methods.

// In Go 1.18+, 'any' is an alias for 'interface{}'

// ============================================================================
// 1. BASIC EMPTY INTERFACE USAGE
// ============================================================================

// PrintAnything accepts any type
func PrintAnything(value interface{}) {
	fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// PrintAny accepts any type (using 'any' alias from Go 1.18+)
func PrintAny(value any) {
	fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// ============================================================================
// 2. STORING DIFFERENT TYPES
// ============================================================================

// Container holds any type of value
type Container struct {
	Value interface{}
}

// NewContainer creates a container
func NewContainer(value interface{}) *Container {
	return &Container{Value: value}
}

// GetValue returns the stored value
func (c *Container) GetValue() interface{} {
	return c.Value
}

// SetValue updates the stored value
func (c *Container) SetValue(value interface{}) {
	c.Value = value
}

// ============================================================================
// 3. COLLECTIONS WITH MIXED TYPES
// ============================================================================

// Store creates a generic key-value store
type Store struct {
	data map[string]interface{}
}

// NewStore creates a new store
func NewStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

// Set stores a value of any type
func (s *Store) Set(key string, value interface{}) {
	s.data[key] = value
}

// Get retrieves a value
func (s *Store) Get(key string) (interface{}, bool) {
	value, exists := s.data[key]
	return value, exists
}

// GetAll returns all stored values
func (s *Store) GetAll() map[string]interface{} {
	return s.data
}

// ============================================================================
// 4. FUNCTIONS RETURNING interface{}
// ============================================================================

// GetConfigValue simulates getting configuration values
func GetConfigValue(key string) interface{} {
	config := map[string]interface{}{
		"port":       8080,
		"host":       "localhost",
		"debug":      true,
		"timeout":    30.5,
		"allowed_ips": []string{"127.0.0.1", "192.168.1.1"},
	}
	return config[key]
}

// ============================================================================
// 5. VARIADIC FUNCTIONS WITH interface{}
// ============================================================================

// PrintAll prints all arguments of any type
func PrintAll(values ...interface{}) {
	fmt.Println("Printing all values:")
	for i, v := range values {
		fmt.Printf("  %d: %v (type: %T)\n", i, v, v)
	}
}

// Sum attempts to sum any numeric values
func Sum(values ...interface{}) float64 {
	total := 0.0
	for _, v := range values {
		switch val := v.(type) {
		case int:
			total += float64(val)
		case int64:
			total += float64(val)
		case float64:
			total += val
		case float32:
			total += float64(val)
		default:
			fmt.Printf("Warning: Cannot sum type %T\n", v)
		}
	}
	return total
}

// ============================================================================
// 6. REFLECTION WITH EMPTY INTERFACE
// ============================================================================

// Describe provides detailed information about a value
func Describe(value interface{}) {
	fmt.Printf("\n=== Describing Value ===\n")
	fmt.Printf("Value: %v\n", value)
	fmt.Printf("Type: %T\n", value)

	// Check for nil
	if value == nil {
		fmt.Println("Value is nil")
		return
	}

	// Type-specific information
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer value: %d\n", v)
		fmt.Printf("Is positive: %v\n", v > 0)
	case string:
		fmt.Printf("String value: %q\n", v)
		fmt.Printf("Length: %d\n", len(v))
	case []interface{}:
		fmt.Printf("Slice with %d elements\n", len(v))
	case map[string]interface{}:
		fmt.Printf("Map with %d keys\n", len(v))
	default:
		fmt.Printf("Other type: %T\n", v)
	}
}

// ============================================================================
// 7. JSON-LIKE DATA STRUCTURES
// ============================================================================

// JSON simulates JSON-like nested data
type JSON map[string]interface{}

// NewJSON creates a new JSON object
func NewJSON() JSON {
	return make(JSON)
}

// Set sets a key-value pair
func (j JSON) Set(key string, value interface{}) {
	j[key] = value
}

// Get retrieves a value
func (j JSON) Get(key string) interface{} {
	return j[key]
}

// GetString retrieves a string value
func (j JSON) GetString(key string) string {
	if v, ok := j[key].(string); ok {
		return v
	}
	return ""
}

// GetInt retrieves an int value
func (j JSON) GetInt(key string) int {
	if v, ok := j[key].(int); ok {
		return v
	}
	return 0
}

// ============================================================================
// 8. DANGERS OF EMPTY INTERFACE
// ============================================================================

// BadAdd shows the danger of empty interface - no type safety!
func BadAdd(a, b interface{}) interface{} {
	// This compiles but will panic at runtime if wrong types
	// return a.(int) + b.(int) // Dangerous!

	// Better: handle with type assertions
	aInt, aOk := a.(int)
	bInt, bOk := b.(int)

	if aOk && bOk {
		return aInt + bInt
	}
	return nil
}

// GoodAdd uses concrete types - type safe!
func GoodAdd(a, b int) int {
	return a + b
}

// ============================================================================
// 9. WORKING WITH SLICES OF interface{}
// ============================================================================

// Filter filters a slice based on a predicate
func Filter(slice []interface{}, predicate func(interface{}) bool) []interface{} {
	result := make([]interface{}, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map applies a function to each element
func Map(slice []interface{}, mapper func(interface{}) interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

func main() {
	fmt.Println("=== 1. Basic Empty Interface ===")

	// Can accept any type
	PrintAnything(42)
	PrintAnything("Hello")
	PrintAnything(3.14)
	PrintAnything(true)
	PrintAnything([]int{1, 2, 3})

	fmt.Println("\n=== 2. Using 'any' (Go 1.18+) ===")

	// 'any' is clearer and more idiomatic
	PrintAny(100)
	PrintAny("World")
	PrintAny(map[string]int{"age": 30})

	fmt.Println("\n=== 3. Container Example ===")

	// Container can hold different types
	intContainer := NewContainer(42)
	fmt.Printf("Container value: %v\n", intContainer.GetValue())

	// Change to different type
	intContainer.SetValue("Now I'm a string!")
	fmt.Printf("Container value: %v\n", intContainer.GetValue())

	intContainer.SetValue([]int{1, 2, 3, 4, 5})
	fmt.Printf("Container value: %v\n", intContainer.GetValue())

	fmt.Println("\n=== 4. Generic Store ===")

	store := NewStore()

	// Store different types
	store.Set("name", "Alice")
	store.Set("age", 30)
	store.Set("salary", 75000.50)
	store.Set("active", true)
	store.Set("skills", []string{"Go", "Python", "JavaScript"})

	// Retrieve values
	name, _ := store.Get("name")
	age, _ := store.Get("age")

	fmt.Printf("Name: %v (type: %T)\n", name, name)
	fmt.Printf("Age: %v (type: %T)\n", age, age)

	// Print all
	fmt.Println("\nAll stored values:")
	for key, value := range store.GetAll() {
		fmt.Printf("  %s: %v (type: %T)\n", key, value, value)
	}

	fmt.Println("\n=== 5. Configuration Example ===")

	port := GetConfigValue("port")
	host := GetConfigValue("host")
	debug := GetConfigValue("debug")

	fmt.Printf("Port: %v (type: %T)\n", port, port)
	fmt.Printf("Host: %v (type: %T)\n", host, host)
	fmt.Printf("Debug: %v (type: %T)\n", debug, debug)

	fmt.Println("\n=== 6. Variadic Functions ===")

	PrintAll(1, "two", 3.0, true, []int{4, 5})

	// Sum different numeric types
	total := Sum(10, 20.5, int64(30), float32(15.5), "ignored")
	fmt.Printf("Total: %.2f\n", total)

	fmt.Println("\n=== 7. Describe Function ===")

	Describe(42)
	Describe("Hello, Go!")
	Describe([]interface{}{1, "two", 3.0})
	Describe(map[string]interface{}{"key": "value"})
	Describe(nil)

	fmt.Println("\n=== 8. JSON-Like Data ===")

	user := NewJSON()
	user.Set("id", 123)
	user.Set("name", "John Doe")
	user.Set("email", "john@example.com")
	user.Set("active", true)
	user.Set("tags", []string{"admin", "developer"})

	// Nested object
	address := NewJSON()
	address.Set("city", "San Francisco")
	address.Set("state", "CA")
	user.Set("address", address)

	fmt.Println("User data:")
	for key, value := range user {
		fmt.Printf("  %s: %v (type: %T)\n", key, value, value)
	}

	// Type-safe retrieval
	name = user.GetString("name")
	id := user.GetInt("id")
	fmt.Printf("\nName: %s, ID: %d\n", name, id)

	fmt.Println("\n=== 9. Type Safety Comparison ===")

	// Bad: Using interface{} - no compile-time safety
	result := BadAdd(10, 20)
	fmt.Printf("BadAdd result: %v\n", result)

	// This would panic at runtime:
	// BadAdd("10", "20") // Panic!

	// Good: Using concrete types - type safe
	result2 := GoodAdd(10, 20)
	fmt.Printf("GoodAdd result: %d\n", result2)

	// This won't compile:
	// GoodAdd("10", "20") // ✗ Compile error!

	fmt.Println("\n=== 10. Slice Operations ===")

	mixed := []interface{}{1, "two", 3, "four", 5, "six", 7}

	// Filter for integers only
	isInt := func(v interface{}) bool {
		_, ok := v.(int)
		return ok
	}

	integers := Filter(mixed, isInt)
	fmt.Printf("Integers only: %v\n", integers)

	// Map: double all integers
	doubler := func(v interface{}) interface{} {
		if num, ok := v.(int); ok {
			return num * 2
		}
		return v
	}

	doubled := Map(integers, doubler)
	fmt.Printf("Doubled: %v\n", doubled)

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("✓ Prefer concrete types over interface{}/any when possible")
	fmt.Println("✓ Use interface{}/any for truly generic operations")
	fmt.Println("✓ Always use type assertions or type switches when extracting values")
	fmt.Println("✓ Document what types are expected/returned")
	fmt.Println("✓ In Go 1.18+, prefer 'any' over 'interface{}' for clarity")
	fmt.Println("✓ Consider using generics (Go 1.18+) instead of interface{}")
	fmt.Println("✗ Don't use interface{} just to avoid thinking about types")
	fmt.Println("✗ Be careful: interface{} removes type safety!")

	fmt.Println("\n=== When to Use interface{}/any ===")
	fmt.Println("Good use cases:")
	fmt.Println("  • JSON encoding/decoding")
	fmt.Println("  • Configuration systems")
	fmt.Println("  • Logging frameworks (fmt.Println, log.Printf)")
	fmt.Println("  • Generic data structures (before Go 1.18)")
	fmt.Println("  • Working with reflection")
	fmt.Println("\nAvoid for:")
	fmt.Println("  • Business logic")
	fmt.Println("  • When types are known at compile time")
	fmt.Println("  • Performance-critical code")
}
