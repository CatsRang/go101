package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ============================================================================
// ERROR WRAPPING IN GO
// ============================================================================
// Error wrapping adds context while preserving the original error.
// Use %w with fmt.Errorf() to wrap errors (Go 1.13+)

// ============================================================================
// 1. BASIC ERROR WRAPPING WITH %w
// ============================================================================

var ErrDatabase = errors.New("database error")
var ErrNetwork = errors.New("network error")
var ErrPermission = errors.New("permission denied")

// connectDB simulates database connection
func connectDB() error {
	// Simulate connection failure
	return ErrDatabase
}

// fetchUser wraps the database error with context
func fetchUser(id int) error {
	err := connectDB()
	if err != nil {
		// %w wraps the error, preserving it for errors.Is() and errors.As()
		return fmt.Errorf("failed to fetch user %d: %w", id, err)
	}
	return nil
}

// getUserData adds another layer of context
func getUserData(id int) error {
	err := fetchUser(id)
	if err != nil {
		return fmt.Errorf("getUserData: %w", err)
	}
	return nil
}

// ============================================================================
// 2. ERROR WRAPPING VS ERROR FORMATTING
// ============================================================================

// ❌ BAD: Using %v loses the original error
func badWrap(id int) error {
	err := ErrDatabase
	// %v converts error to string - original error is lost!
	return fmt.Errorf("failed to connect: %v", err)
}

// ✓ GOOD: Using %w preserves the original error
func goodWrap(id int) error {
	err := ErrDatabase
	// %w wraps error - original is preserved
	return fmt.Errorf("failed to connect: %w", err)
}

func demonstrateWrapVsFormat() {
	// Bad wrapping
	badErr := badWrap(1)
	fmt.Printf("Bad wrap - errors.Is(ErrDatabase): %v\n", errors.Is(badErr, ErrDatabase))

	// Good wrapping
	goodErr := goodWrap(1)
	fmt.Printf("Good wrap - errors.Is(ErrDatabase): %v\n", errors.Is(goodErr, ErrDatabase))
}

// ============================================================================
// 3. UNWRAPPING ERRORS
// ============================================================================

// CustomError with Unwrap method
type CustomError struct {
	Message string
	Err     error
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", ce.Message, ce.Err)
}

// Unwrap returns the wrapped error
func (ce *CustomError) Unwrap() error {
	return ce.Err
}

// processData demonstrates custom error wrapping
func processData(data string) error {
	if data == "" {
		return &CustomError{
			Message: "invalid data",
			Err:     errors.New("data cannot be empty"),
		}
	}
	return nil
}

// ============================================================================
// 4. ERROR CHAINS
// ============================================================================

// Simulate a call stack with error wrapping
func layer1() error {
	return errors.New("root cause: file not found")
}

func layer2() error {
	err := layer1()
	if err != nil {
		return fmt.Errorf("layer2: failed to read file: %w", err)
	}
	return nil
}

func layer3() error {
	err := layer2()
	if err != nil {
		return fmt.Errorf("layer3: data processing failed: %w", err)
	}
	return nil
}

func layer4() error {
	err := layer3()
	if err != nil {
		return fmt.Errorf("layer4: request handling failed: %w", err)
	}
	return nil
}

// ============================================================================
// 5. CHECKING WRAPPED ERRORS WITH errors.Is()
// ============================================================================

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal error")
)

func getResource(id int) error {
	if id == 0 {
		return fmt.Errorf("invalid ID: %w", ErrBadRequest)
	}
	if id == 404 {
		return fmt.Errorf("resource %d: %w", id, ErrNotFound)
	}
	return nil
}

func handleRequest(id int) error {
	err := getResource(id)
	if err != nil {
		return fmt.Errorf("handleRequest failed: %w", err)
	}
	return nil
}

func demonstrateErrorsIs() {
	err := handleRequest(404)

	// errors.Is() unwraps the error chain to find ErrNotFound
	if errors.Is(err, ErrNotFound) {
		fmt.Println("✓ Detected NotFound error in chain")
	}

	// Check for other errors
	if errors.Is(err, ErrUnauthorized) {
		fmt.Println("Is unauthorized")
	} else {
		fmt.Println("✓ Correctly not unauthorized")
	}
}

// ============================================================================
// 6. TYPE ASSERTIONS WITH errors.As()
// ============================================================================

// ValidationError represents validation failures
type ValidationError struct {
	Field string
	Value interface{}
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field=%s, value=%v", ve.Field, ve.Value)
}

func validateInput(email string) error {
	if !strings.Contains(email, "@") {
		err := &ValidationError{
			Field: "email",
			Value: email,
		}
		return fmt.Errorf("input validation: %w", err)
	}
	return nil
}

func processForm(email string) error {
	err := validateInput(email)
	if err != nil {
		return fmt.Errorf("form processing: %w", err)
	}
	return nil
}

func demonstrateErrorsAs() {
	err := processForm("invalid-email")

	// errors.As() finds ValidationError in the chain
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("✓ Found ValidationError: field=%s, value=%v\n", ve.Field, ve.Value)
	}
}

// ============================================================================
// 7. MULTIPLE LEVELS OF WRAPPING
// ============================================================================

type RepositoryError struct {
	Op  string
	Err error
}

func (re *RepositoryError) Error() string {
	return fmt.Sprintf("repository.%s: %v", re.Op, re.Err)
}

func (re *RepositoryError) Unwrap() error {
	return re.Err
}

type ServiceError struct {
	Service string
	Err     error
}

func (se *ServiceError) Error() string {
	return fmt.Sprintf("service.%s: %v", se.Service, se.Err)
}

func (se *ServiceError) Unwrap() error {
	return se.Err
}

// Repository layer
func repoGetUser(id int) error {
	return &RepositoryError{
		Op:  "GetUser",
		Err: ErrNotFound,
	}
}

// Service layer
func serviceGetUser(id int) error {
	err := repoGetUser(id)
	if err != nil {
		return &ServiceError{
			Service: "UserService",
			Err:     err,
		}
	}
	return nil
}

// Controller layer
func controllerGetUser(id int) error {
	err := serviceGetUser(id)
	if err != nil {
		return fmt.Errorf("controller: failed to get user %d: %w", id, err)
	}
	return nil
}

// ============================================================================
// 8. ADDING CONTEXT AT EACH LAYER
// ============================================================================

func readFile(filename string) error {
	return io.EOF // Simulate EOF
}

func parseConfig(filename string) error {
	err := readFile(filename)
	if err != nil {
		return fmt.Errorf("parse config from %q: %w", filename, err)
	}
	return nil
}

func loadConfig() error {
	err := parseConfig("config.json")
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	return nil
}

func initializeApp() error {
	err := loadConfig()
	if err != nil {
		return fmt.Errorf("application initialization: %w", err)
	}
	return nil
}

// ============================================================================
// 9. PRACTICAL EXAMPLE: HTTP HANDLER ERROR CHAIN
// ============================================================================

type HTTPError struct {
	StatusCode int
	Err        error
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %v", he.StatusCode, he.Err)
}

func (he *HTTPError) Unwrap() error {
	return he.Err
}

func dbQuery(query string) error {
	return errors.New("connection timeout")
}

func getUserByID(id int) error {
	err := dbQuery("SELECT * FROM users")
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}
	return nil
}

func handleHTTPRequest(userID int) error {
	err := getUserByID(userID)
	if err != nil {
		return &HTTPError{
			StatusCode: 500,
			Err:        fmt.Errorf("failed to get user %d: %w", userID, err),
		}
	}
	return nil
}

// ============================================================================
// 10. BEST PRACTICES
// ============================================================================

// ✓ GOOD: Add context at each layer
func goodErrorPropagation() error {
	err := errors.New("disk full")
	if err != nil {
		// Add context about what we were doing
		return fmt.Errorf("failed to write log file: %w", err)
	}
	return nil
}

// ❌ BAD: Just returning the error without context
func badErrorPropagation() error {
	err := errors.New("disk full")
	// No context - user doesn't know what failed
	return err
}

// ✓ GOOD: Preserve error for checking
func goodErrorCheck(id int) error {
	if id == 0 {
		return fmt.Errorf("invalid user ID: %w", ErrBadRequest)
	}
	return nil
}

// ❌ BAD: Using %v loses the original error
func badErrorCheck(id int) error {
	if id == 0 {
		return fmt.Errorf("invalid user ID: %v", ErrBadRequest)
	}
	return nil
}

func main() {
	fmt.Println("=== 1. Basic Error Wrapping ===")
	err := getUserData(123)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is ErrDatabase: %v\n", errors.Is(err, ErrDatabase))

	fmt.Println("\n=== 2. Wrap vs Format ===")
	demonstrateWrapVsFormat()

	fmt.Println("\n=== 3. Custom Error with Unwrap ===")
	err = processData("")
	fmt.Printf("Error: %v\n", err)
	if errors.Is(err, errors.New("data cannot be empty")) {
		fmt.Println("Contains 'data cannot be empty'")
	}

	fmt.Println("\n=== 4. Error Chain ===")
	err = layer4()
	fmt.Printf("Full error: %v\n", err)

	fmt.Println("\n=== 5. errors.Is() with Wrapped Errors ===")
	demonstrateErrorsIs()

	fmt.Println("\n=== 6. errors.As() with Wrapped Errors ===")
	demonstrateErrorsAs()

	fmt.Println("\n=== 7. Multiple Layers ===")
	err = controllerGetUser(123)
	fmt.Printf("Error chain: %v\n", err)
	fmt.Printf("Is ErrNotFound: %v\n", errors.Is(err, ErrNotFound))

	// Check each layer
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		fmt.Printf("✓ Found RepositoryError: Op=%s\n", repoErr.Op)
	}

	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		fmt.Printf("✓ Found ServiceError: Service=%s\n", svcErr.Service)
	}

	fmt.Println("\n=== 8. Context at Each Layer ===")
	err = initializeApp()
	fmt.Printf("Error with full context: %v\n", err)
	fmt.Printf("Is io.EOF: %v\n", errors.Is(err, io.EOF))

	fmt.Println("\n=== 9. HTTP Handler Example ===")
	err = handleHTTPRequest(123)
	fmt.Printf("HTTP Error: %v\n", err)

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		fmt.Printf("✓ HTTP Status Code: %d\n", httpErr.StatusCode)
	}

	fmt.Println("\n=== 10. Good vs Bad Error Handling ===")
	goodErr := goodErrorPropagation()
	badErr := badErrorPropagation()

	fmt.Printf("Good error (with context): %v\n", goodErr)
	fmt.Printf("Bad error (no context): %v\n", badErr)

	fmt.Println("\n=== ERROR WRAPPING BEST PRACTICES ===")
	fmt.Println("✓ Use %w to wrap errors, not %v")
	fmt.Println("✓ Add context at each layer of the call stack")
	fmt.Println("✓ Use errors.Is() to check for specific errors")
	fmt.Println("✓ Use errors.As() to extract custom error types")
	fmt.Println("✓ Implement Unwrap() for custom error types")
	fmt.Println("✓ Wrap errors from external libraries")
	fmt.Println("✓ Include relevant information (IDs, names, operations)")
	fmt.Println("✓ Build error chains from bottom to top")
	fmt.Println("✗ Don't use %v when you need to check the error later")
	fmt.Println("✗ Don't wrap errors unnecessarily (avoid deep chains)")
	fmt.Println("✗ Don't include sensitive data in error messages")

	fmt.Println("\n=== WHEN TO WRAP ERRORS ===")
	fmt.Println("Wrap when:")
	fmt.Println("  • Adding context about what operation failed")
	fmt.Println("  • Passing errors up the call stack")
	fmt.Println("  • You need to preserve the original error for checking")
	fmt.Println("\nDon't wrap when:")
	fmt.Println("  • The error message is already clear")
	fmt.Println("  • You're at the final error handler")
	fmt.Println("  • Wrapping adds no value")
}
