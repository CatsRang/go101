package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ============================================================================
// BASIC ERROR HANDLING IN GO
// ============================================================================
// Go uses explicit error handling with the error interface.
// Errors are values, not exceptions.

// ============================================================================
// 1. THE ERROR INTERFACE
// ============================================================================
// type error interface {
//     Error() string
// }

// ============================================================================
// 2. CREATING ERRORS
// ============================================================================

// Method 1: errors.New() - simplest way
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Method 2: fmt.Errorf() - formatted error messages
func validateAge(age int) error {
	if age < 0 {
		return fmt.Errorf("age cannot be negative: %d", age)
	}
	if age > 150 {
		return fmt.Errorf("age too large: %d", age)
	}
	return nil
}

// Method 3: Returning nil for no error
func processName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	// Success case - return nil
	return nil
}

// ============================================================================
// 3. ERROR HANDLING PATTERNS
// ============================================================================

// Pattern 1: Immediate error check (most common)
func parseAndDouble(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string: %v", err)
	}
	return num * 2, nil
}

// Pattern 2: Multiple error checks
func processUser(name string, age int, email string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if age < 18 {
		return errors.New("user must be 18 or older")
	}
	if email == "" {
		return errors.New("email is required")
	}
	// All validations passed
	return nil
}

// Pattern 3: Named return values for cleaner code
func findUser(id int) (name string, email string, err error) {
	if id <= 0 {
		err = errors.New("invalid user ID")
		return // returns zero values for name and email, plus the error
	}

	// Simulate database lookup
	if id == 1 {
		name = "Alice"
		email = "alice@example.com"
		return // err is nil by default
	}

	err = fmt.Errorf("user not found: %d", id)
	return
}

// ============================================================================
// 4. ERROR CHECKING AND HANDLING
// ============================================================================

// handleDivision demonstrates basic error checking
func handleDivision() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %.2f\n", result)

	// This will error
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %.2f\n", result)
}

// ============================================================================
// 5. SENTINEL ERRORS - Predefined Error Values
// ============================================================================

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")
	ErrTimeout       = errors.New("operation timed out")
	ErrAlreadyExists = errors.New("already exists")
)

// lookup simulates a database lookup
func lookup(id int) (string, error) {
	if id == 0 {
		return "", ErrInvalidInput
	}
	if id == 999 {
		return "", ErrNotFound
	}
	return fmt.Sprintf("User-%d", id), nil
}

// Using sentinel errors allows callers to check for specific errors
func demonstrateSentinelErrors() {
	_, err := lookup(999)
	if err == ErrNotFound {
		fmt.Println("Handling not found error - maybe create the user?")
	} else if err == ErrInvalidInput {
		fmt.Println("Handling invalid input - ask user to retry")
	} else if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	}
}

// ============================================================================
// 6. ERRORS.IS() - Checking Error Identity
// ============================================================================

func getUserData(id int) (string, error) {
	if id == 0 {
		return "", ErrInvalidInput
	}
	if id < 0 {
		return "", fmt.Errorf("negative ID not allowed: %w", ErrInvalidInput)
	}
	return "User data", nil
}

func demonstrateErrorsIs() {
	_, err := getUserData(-5)

	// errors.Is() checks if err is or wraps ErrInvalidInput
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("Detected invalid input error (even if wrapped)")
	}
}

// ============================================================================
// 7. ERROR HANDLING GUIDELINES
// ============================================================================

// ❌ BAD: Ignoring errors
func badErrorHandling() {
	parseAndDouble("not a number") // Error ignored!
	// This is dangerous and can cause bugs
}

// ✓ GOOD: Always check errors
func goodErrorHandling() {
	_, err := parseAndDouble("not a number")
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		// Handle appropriately: log, return, retry, etc.
	}
}

// ❌ BAD: Generic error messages
func badErrorMessage(filename string) error {
	// Opening file fails...
	return errors.New("error")
}

// ✓ GOOD: Descriptive error messages
func goodErrorMessage(filename string) error {
	// Opening file fails...
	return fmt.Errorf("failed to open file %q: file does not exist", filename)
}

// ============================================================================
// 8. EARLY RETURN PATTERN
// ============================================================================

// ✓ GOOD: Early returns make code cleaner
func validateUserInput(name string, age int, email string) error {
	if name == "" {
		return errors.New("name is required")
	}

	if age < 18 {
		return errors.New("must be 18 or older")
	}

	if email == "" {
		return errors.New("email is required")
	}

	// Happy path - all validations passed
	fmt.Println("User input is valid")
	return nil
}

// ❌ BAD: Nested if statements (avoid this)
func nestedValidation(name string, age int, email string) error {
	if name != "" {
		if age >= 18 {
			if email != "" {
				fmt.Println("User input is valid")
				return nil
			} else {
				return errors.New("email is required")
			}
		} else {
			return errors.New("must be 18 or older")
		}
	}
	return errors.New("name is required")
}

// ============================================================================
// 9. WORKING WITH STANDARD LIBRARY ERRORS
// ============================================================================

func demonstrateStdlibErrors() {
	// strconv errors
	_, err := strconv.Atoi("not a number")
	if err != nil {
		fmt.Printf("strconv.Atoi error: %v\n", err)
		fmt.Printf("Error type: %T\n", err)
	}

	// fmt.Errorf creates formatted errors
	userID := 123
	err = fmt.Errorf("user %d not found in database", userID)
	fmt.Printf("Formatted error: %v\n", err)
}

// ============================================================================
// 10. MULTIPLE RETURN VALUES WITH ERRORS
// ============================================================================

// Go convention: error is always the last return value
func calculate(a, b int, operation string) (int, error) {
	switch operation {
	case "add":
		return a + b, nil
	case "subtract":
		return a - b, nil
	case "multiply":
		return a * b, nil
	case "divide":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", operation)
	}
}

func main() {
	fmt.Println("=== 1. Basic Error Handling ===")
	handleDivision()

	fmt.Println("\n=== 2. Validation Errors ===")
	if err := validateAge(-5); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}
	if err := validateAge(200); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}
	if err := validateAge(25); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Println("Age is valid")
	}

	fmt.Println("\n=== 3. Parse and Double ===")
	result, err := parseAndDouble("42")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %d\n", result)
	}

	result, err = parseAndDouble("not a number")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== 4. User Processing ===")
	err = processUser("Alice", 25, "alice@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("User is valid")
	}

	err = processUser("", 25, "alice@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	err = processUser("Bob", 16, "bob@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== 5. Find User (Named Returns) ===")
	name, email, err := findUser(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found user: %s (%s)\n", name, email)
	}

	name, email, err = findUser(999)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== 6. Sentinel Errors ===")
	demonstrateSentinelErrors()

	user, err := lookup(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found: %s\n", user)
	}

	fmt.Println("\n=== 7. errors.Is() ===")
	demonstrateErrorsIs()

	fmt.Println("\n=== 8. Early Return Pattern ===")
	validateUserInput("Alice", 25, "alice@example.com")
	validateUserInput("", 25, "alice@example.com")

	fmt.Println("\n=== 9. Standard Library Errors ===")
	demonstrateStdlibErrors()

	fmt.Println("\n=== 10. Calculator Example ===")
	operations := []string{"add", "subtract", "multiply", "divide", "modulo"}
	for _, op := range operations {
		result, err := calculate(10, 5, op)
		if err != nil {
			fmt.Printf("%s: Error - %v\n", op, err)
		} else {
			fmt.Printf("%s: 10 %s 5 = %d\n", op, op, result)
		}
	}

	// Division by zero
	result, err = calculate(10, 0, "divide")
	if err != nil {
		fmt.Printf("Division by zero: %v\n", err)
	}

	fmt.Println("\n=== ERROR HANDLING PRINCIPLES ===")
	fmt.Println("✓ Errors are values, not exceptions")
	fmt.Println("✓ Always check returned errors")
	fmt.Println("✓ Error is conventionally the last return value")
	fmt.Println("✓ Return nil to indicate no error")
	fmt.Println("✓ Use errors.New() for simple error messages")
	fmt.Println("✓ Use fmt.Errorf() for formatted error messages")
	fmt.Println("✓ Define sentinel errors for common error cases")
	fmt.Println("✓ Use errors.Is() to check for specific errors")
	fmt.Println("✓ Provide context in error messages")
	fmt.Println("✓ Use early returns to avoid nested if statements")
	fmt.Println("✗ Never ignore errors (always assign to _)")
	fmt.Println("✗ Don't panic for normal error conditions")
}
