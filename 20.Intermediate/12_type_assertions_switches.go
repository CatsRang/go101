package main

import (
	"fmt"
	"strconv"
)

// ============================================================================
// TYPE ASSERTIONS AND TYPE SWITCHES
// ============================================================================

// ============================================================================
// 1. BASIC TYPE ASSERTIONS
// ============================================================================

// Type assertion syntax: value.(Type)
// Returns the underlying value if assertion succeeds, panics otherwise

func demonstrateBasicAssertion() {
	var i interface{} = "hello"

	// Single-value assertion (panics if wrong type)
	s := i.(string)
	fmt.Printf("String value: %s\n", s)

	// This would panic:
	// n := i.(int) // panic: interface conversion: interface {} is string, not int

	// Two-value assertion (safe, returns bool)
	s2, ok := i.(string)
	if ok {
		fmt.Printf("Successfully asserted as string: %s\n", s2)
	}

	// Check for wrong type
	n, ok := i.(int)
	if !ok {
		fmt.Printf("Failed to assert as int (value is %v, ok is %v)\n", n, ok)
	}
}

// ============================================================================
// 2. TYPE SWITCHES
// ============================================================================

// ProcessValue handles different types using type switch
func ProcessValue(value interface{}) {
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer: %d (doubled: %d)\n", v, v*2)
	case string:
		fmt.Printf("String: %q (length: %d)\n", v, len(v))
	case bool:
		fmt.Printf("Boolean: %v (negated: %v)\n", v, !v)
	case float64:
		fmt.Printf("Float64: %.2f (squared: %.2f)\n", v, v*v)
	case []int:
		fmt.Printf("Int slice: %v (sum: %d)\n", v, sumInts(v))
	case []string:
		fmt.Printf("String slice: %v (joined: %s)\n", v, joinStrings(v))
	case nil:
		fmt.Println("Nil value")
	default:
		fmt.Printf("Unknown type: %T (value: %v)\n", v, v)
	}
}

func sumInts(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// ============================================================================
// 3. INTERFACE TYPE ASSERTIONS
// ============================================================================

// Shape interface
type Shape interface {
	Area() float64
}

// Perimeter interface
type Perimeter interface {
	Perimeter() float64
}

// Rectangle implements both Shape and Perimeter
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle implements both Shape and Perimeter
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// ProcessShape demonstrates interface type assertions
func ProcessShape(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())

	// Check if shape also implements Perimeter
	if p, ok := s.(Perimeter); ok {
		fmt.Printf("Perimeter: %.2f\n", p.Perimeter())
	} else {
		fmt.Println("Shape doesn't implement Perimeter")
	}

	// Check specific type
	if rect, ok := s.(Rectangle); ok {
		fmt.Printf("It's a rectangle: %.2f x %.2f\n", rect.Width, rect.Height)
	} else if circle, ok := s.(Circle); ok {
		fmt.Printf("It's a circle with radius: %.2f\n", circle.Radius)
	}
}

// ============================================================================
// 4. TYPE ASSERTION FOR ERROR HANDLING
// ============================================================================

// CustomError represents a custom error type
type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// NetworkError represents network-specific errors
type NetworkError struct {
	Host string
	Port int
	Err  error
}

func (ne NetworkError) Error() string {
	return fmt.Sprintf("Network error at %s:%d - %v", ne.Host, ne.Port, ne.Err)
}

// HandleError demonstrates type assertion for error handling
func HandleError(err error) {
	if err == nil {
		fmt.Println("No error")
		return
	}

	// Type switch for different error types
	switch e := err.(type) {
	case CustomError:
		fmt.Printf("Custom error (code %d): %s\n", e.Code, e.Message)
		if e.Code >= 500 {
			fmt.Println("  → Server error, retrying...")
		}
	case NetworkError:
		fmt.Printf("Network error at %s:%d\n", e.Host, e.Port)
		fmt.Println("  → Check network connection")
	default:
		fmt.Printf("Generic error: %v\n", err)
	}
}

// ============================================================================
// 5. TYPE ASSERTIONS WITH MULTIPLE RETURN VALUES
// ============================================================================

// SafeConvert safely converts interface{} to various types
type SafeConvert struct{}

func (sc SafeConvert) ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

func (sc SafeConvert) ToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return fmt.Sprintf("%.2f", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", value)
	}
}

// ============================================================================
// 6. CHECKING MULTIPLE INTERFACES
// ============================================================================

// Reader interface
type Reader interface {
	Read() string
}

// Writer interface
type Writer interface {
	Write(data string)
}

// Closer interface
type Closer interface {
	Close() error
}

// File implements all three
type File struct {
	name   string
	data   string
	closed bool
}

func (f *File) Read() string {
	if f.closed {
		return ""
	}
	return f.data
}

func (f *File) Write(data string) {
	if !f.closed {
		f.data += data
	}
}

func (f *File) Close() error {
	f.closed = true
	return nil
}

// ProcessResource checks what operations are available
func ProcessResource(resource interface{}) {
	fmt.Printf("\nProcessing resource of type %T:\n", resource)

	// Check for Reader
	if r, ok := resource.(Reader); ok {
		fmt.Printf("  Can read: %s\n", r.Read())
	}

	// Check for Writer
	if w, ok := resource.(Writer); ok {
		w.Write(" [modified]")
		fmt.Println("  Can write: data written")
	}

	// Check for Closer
	if c, ok := resource.(Closer); ok {
		c.Close()
		fmt.Println("  Can close: resource closed")
	}
}

// ============================================================================
// 7. TYPE ASSERTION IN COMPLEX SCENARIOS
// ============================================================================

// Data represents generic data structure
type Data map[string]interface{}

// GetString safely retrieves string value
func (d Data) GetString(key string) string {
	if val, ok := d[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt safely retrieves int value
func (d Data) GetInt(key string) int {
	if val, ok := d[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetSlice safely retrieves slice
func (d Data) GetSlice(key string) []interface{} {
	if val, ok := d[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			return slice
		}
	}
	return nil
}

// ============================================================================
// 8. PERFORMANCE CONSIDERATIONS
// ============================================================================

// CompareAssertions demonstrates different assertion patterns
func CompareAssertions() {
	var value interface{} = 42

	// Pattern 1: Single assertion (panics on failure)
	// Use when you're certain of the type
	_ = value.(int)

	// Pattern 2: Safe assertion (returns bool)
	// Use when type is uncertain
	if v, ok := value.(int); ok {
		_ = v
	}

	// Pattern 3: Type switch (for multiple possibilities)
	// Use when handling multiple types
	switch value.(type) {
	case int:
		// handle int
	case string:
		// handle string
	}
}

func main() {
	fmt.Println("=== 1. Basic Type Assertions ===")
	demonstrateBasicAssertion()

	fmt.Println("\n=== 2. Type Switches ===")
	ProcessValue(42)
	ProcessValue("Hello, Go!")
	ProcessValue(true)
	ProcessValue(3.14159)
	ProcessValue([]int{1, 2, 3, 4, 5})
	ProcessValue([]string{"Go", "Python", "JavaScript"})
	ProcessValue(nil)
	ProcessValue(struct{ Name string }{"Unknown"})

	fmt.Println("\n=== 3. Interface Type Assertions ===")
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}

	fmt.Println("Rectangle:")
	ProcessShape(rect)

	fmt.Println("\nCircle:")
	ProcessShape(circle)

	fmt.Println("\n=== 4. Error Type Assertions ===")
	HandleError(nil)
	HandleError(CustomError{Code: 404, Message: "Not Found"})
	HandleError(CustomError{Code: 500, Message: "Internal Server Error"})
	HandleError(NetworkError{Host: "api.example.com", Port: 443, Err: fmt.Errorf("timeout")})
	HandleError(fmt.Errorf("generic error"))

	fmt.Println("\n=== 5. Safe Type Conversion ===")
	converter := SafeConvert{}

	values := []interface{}{42, int64(100), 3.14, "123", true, "not a number"}

	for _, val := range values {
		intVal, err := converter.ToInt(val)
		if err != nil {
			fmt.Printf("ToInt(%v): Error - %v\n", val, err)
		} else {
			fmt.Printf("ToInt(%v): %d\n", val, intVal)
		}
	}

	fmt.Println()

	for _, val := range values {
		strVal, err := converter.ToString(val)
		if err != nil {
			fmt.Printf("ToString(%v): Error - %v\n", val, err)
		} else {
			fmt.Printf("ToString(%v): %s\n", val, strVal)
		}
	}

	fmt.Println("\n=== 6. Multiple Interface Checks ===")
	file := &File{name: "test.txt", data: "Hello, World!"}
	ProcessResource(file)

	fmt.Println("\n=== 7. Complex Data Retrieval ===")
	data := Data{
		"name":   "Alice",
		"age":    30,
		"salary": 75000.50,
		"active": true,
		"tags":   []interface{}{"admin", "developer"},
	}

	fmt.Printf("Name: %s\n", data.GetString("name"))
	fmt.Printf("Age: %d\n", data.GetInt("age"))
	fmt.Printf("Salary (as int): %d\n", data.GetInt("salary"))
	fmt.Printf("Tags: %v\n", data.GetSlice("tags"))
	fmt.Printf("Missing key: %s\n", data.GetString("missing"))

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("✓ Use two-value assertion (v, ok) for safety")
	fmt.Println("✓ Use type switches for handling multiple types")
	fmt.Println("✓ Check for nil before asserting")
	fmt.Println("✓ Use type assertions to access additional methods")
	fmt.Println("✓ Document expected types in function signatures")
	fmt.Println("✗ Avoid single-value assertions unless type is guaranteed")
	fmt.Println("✗ Don't overuse empty interfaces")

	fmt.Println("\n=== When to Use Type Assertions ===")
	fmt.Println("• Converting from interface{} to concrete type")
	fmt.Println("• Checking if interface implements additional methods")
	fmt.Println("• Handling different error types")
	fmt.Println("• Working with heterogeneous data (JSON, etc.)")
	fmt.Println("• Implementing type-specific optimizations")
}
