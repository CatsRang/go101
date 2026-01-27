package main

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// CUSTOM ERROR TYPES IN GO
// ============================================================================
// Custom errors allow you to add context and enable type-safe error handling

// ============================================================================
// 1. SIMPLE CUSTOM ERROR
// ============================================================================

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", ve.Field, ve.Message)
}

// validate demonstrates using custom errors
func validate(email string, age int) error {
	if email == "" {
		return ValidationError{
			Field:   "email",
			Message: "cannot be empty",
		}
	}
	if age < 0 {
		return ValidationError{
			Field:   "age",
			Message: "cannot be negative",
		}
	}
	return nil
}

// ============================================================================
// 2. ERROR WITH ADDITIONAL CONTEXT
// ============================================================================

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       interface{}
	Timestamp time.Time
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%v' not found (at %s)",
		nfe.Resource, nfe.ID, nfe.Timestamp.Format("15:04:05"))
}

// findResource simulates resource lookup
func findResource(resourceType string, id int) error {
	if id == 999 {
		return NotFoundError{
			Resource:  resourceType,
			ID:        id,
			Timestamp: time.Now(),
		}
	}
	return nil
}

// ============================================================================
// 3. ERROR WITH ERROR CODE
// ============================================================================

// APIError represents an API error with HTTP-like codes
type APIError struct {
	Code    int
	Message string
	Details string
}

func (ae APIError) Error() string {
	if ae.Details != "" {
		return fmt.Sprintf("[%d] %s - %s", ae.Code, ae.Message, ae.Details)
	}
	return fmt.Sprintf("[%d] %s", ae.Code, ae.Message)
}

// HTTP-like error constructors
func BadRequest(message string) APIError {
	return APIError{Code: 400, Message: message}
}

func Unauthorized(message string) APIError {
	return APIError{Code: 401, Message: message}
}

func NotFound(message string) APIError {
	return APIError{Code: 404, Message: message}
}

func InternalServerError(details string) APIError {
	return APIError{Code: 500, Message: "Internal Server Error", Details: details}
}

// processRequest demonstrates API errors
func processRequest(userID int, token string) error {
	if token == "" {
		return Unauthorized("missing authentication token")
	}
	if userID == 0 {
		return BadRequest("user ID is required")
	}
	if userID == 999 {
		return NotFound("user not found")
	}
	return nil
}

// ============================================================================
// 4. ERROR WITH WRAPPED ERROR
// ============================================================================

// DatabaseError wraps database errors with context
type DatabaseError struct {
	Operation string
	Query     string
	Err       error
}

func (de DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v\nQuery: %s",
		de.Operation, de.Err, de.Query)
}

// Unwrap allows errors.Is and errors.As to work
func (de DatabaseError) Unwrap() error {
	return de.Err
}

// executeQuery simulates database operations
func executeQuery(query string) error {
	// Simulate an error
	if query == "" {
		return DatabaseError{
			Operation: "execute",
			Query:     query,
			Err:       errors.New("empty query"),
		}
	}
	return nil
}

// ============================================================================
// 5. MULTIPLE FIELDS ERROR
// ============================================================================

// MultiFieldError represents errors for multiple fields
type MultiFieldError struct {
	Errors map[string]string
}

func (mfe MultiFieldError) Error() string {
	if len(mfe.Errors) == 0 {
		return "validation error: no errors"
	}
	msg := "validation errors:\n"
	for field, err := range mfe.Errors {
		msg += fmt.Sprintf("  - %s: %s\n", field, err)
	}
	return msg
}

// validateUser validates a user with multiple fields
func validateUser(name, email string, age int) error {
	errs := make(map[string]string)

	if name == "" {
		errs["name"] = "required"
	}
	if len(name) < 3 {
		errs["name"] = "must be at least 3 characters"
	}
	if email == "" {
		errs["email"] = "required"
	}
	if age < 18 {
		errs["age"] = "must be 18 or older"
	}
	if age > 150 {
		errs["age"] = "invalid age"
	}

	if len(errs) > 0 {
		return MultiFieldError{Errors: errs}
	}
	return nil
}

// ============================================================================
// 6. ERROR WITH RETRY INFORMATION
// ============================================================================

// RetryableError indicates an error that can be retried
type RetryableError struct {
	Err        error
	RetryAfter time.Duration
	Attempts   int
}

func (re RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempt %d): %v (retry after %v)",
		re.Attempts, re.Err, re.RetryAfter)
}

func (re RetryableError) Unwrap() error {
	return re.Err
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var re RetryableError
	return errors.As(err, &re)
}

// fetchData simulates a network call
func fetchData(attempt int) error {
	if attempt < 3 {
		return RetryableError{
			Err:        errors.New("network timeout"),
			RetryAfter: time.Second * 2,
			Attempts:   attempt,
		}
	}
	return nil
}

// ============================================================================
// 7. TYPE ASSERTIONS WITH CUSTOM ERRORS
// ============================================================================

// handleError demonstrates type-specific error handling
func handleError(err error) {
	if err == nil {
		return
	}

	// Type assertion for ValidationError
	if ve, ok := err.(ValidationError); ok {
		fmt.Printf("Validation error on field %s: %s\n", ve.Field, ve.Message)
		return
	}

	// Type assertion for NotFoundError
	if nfe, ok := err.(NotFoundError); ok {
		fmt.Printf("Resource not found: %s (ID: %v)\n", nfe.Resource, nfe.ID)
		return
	}

	// Type assertion for APIError
	if ae, ok := err.(APIError); ok {
		fmt.Printf("API Error [%d]: %s\n", ae.Code, ae.Message)
		// Handle based on code
		if ae.Code >= 500 {
			fmt.Println("  → Server error, will retry")
		} else if ae.Code >= 400 {
			fmt.Println("  → Client error, check request")
		}
		return
	}

	// Generic error
	fmt.Printf("Unknown error: %v\n", err)
}

// ============================================================================
// 8. ERRORS.AS() FOR TYPE-SAFE ERROR CHECKING
// ============================================================================

func demonstrateErrorsAs() {
	err := processRequest(999, "valid-token")

	// errors.As() finds the first error in err's chain that matches target
	var apiErr APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("Found APIError: code=%d, message=%s\n", apiErr.Code, apiErr.Message)

		// Access fields directly
		if apiErr.Code == 404 {
			fmt.Println("Handling 404 error...")
		}
	}
}

// ============================================================================
// 9. CUSTOM ERROR WITH METHODS
// ============================================================================

// PermissionError represents a permission denied error
type PermissionError struct {
	User     string
	Resource string
	Action   string
}

func (pe PermissionError) Error() string {
	return fmt.Sprintf("permission denied: user '%s' cannot %s '%s'",
		pe.User, pe.Action, pe.Resource)
}

// IsPermissionDenied checks if error is permission related
func (pe PermissionError) IsPermissionDenied() bool {
	return true
}

// SuggestAction suggests how to resolve the error
func (pe PermissionError) SuggestAction() string {
	return fmt.Sprintf("Contact admin to grant '%s' permission to user '%s'",
		pe.Action, pe.User)
}

// checkPermission simulates permission check
func checkPermission(user, resource, action string) error {
	if user != "admin" {
		return PermissionError{
			User:     user,
			Resource: resource,
			Action:   action,
		}
	}
	return nil
}

// ============================================================================
// 10. ERROR BUILDER PATTERN
// ============================================================================

// ErrorBuilder helps build complex errors
type ErrorBuilder struct {
	code    int
	message string
	details []string
	cause   error
}

func NewErrorBuilder() *ErrorBuilder {
	return &ErrorBuilder{
		details: make([]string, 0),
	}
}

func (eb *ErrorBuilder) WithCode(code int) *ErrorBuilder {
	eb.code = code
	return eb
}

func (eb *ErrorBuilder) WithMessage(message string) *ErrorBuilder {
	eb.message = message
	return eb
}

func (eb *ErrorBuilder) WithDetail(detail string) *ErrorBuilder {
	eb.details = append(eb.details, detail)
	return eb
}

func (eb *ErrorBuilder) WithCause(err error) *ErrorBuilder {
	eb.cause = err
	return eb
}

func (eb *ErrorBuilder) Build() error {
	msg := eb.message
	if len(eb.details) > 0 {
		msg += "\nDetails:"
		for _, d := range eb.details {
			msg += fmt.Sprintf("\n  - %s", d)
		}
	}
	if eb.cause != nil {
		msg += fmt.Sprintf("\nCaused by: %v", eb.cause)
	}
	return fmt.Errorf("[%d] %s", eb.code, msg)
}

func main() {
	fmt.Println("=== 1. Simple Custom Error ===")
	err := validate("", 25)
	handleError(err)

	err = validate("user@example.com", -5)
	handleError(err)

	fmt.Println("\n=== 2. Not Found Error ===")
	err = findResource("User", 999)
	handleError(err)

	fmt.Println("\n=== 3. API Errors ===")
	errors := []error{
		processRequest(0, "token"),
		processRequest(1, ""),
		processRequest(999, "token"),
	}

	for i, err := range errors {
		if err != nil {
			fmt.Printf("%d. ", i+1)
			handleError(err)
		}
	}

	fmt.Println("\n=== 4. Database Error with Unwrap ===")
	err = executeQuery("")
	fmt.Printf("Error: %v\n", err)

	// Check wrapped error
	if errors.Is(err, errors.New("empty query")) {
		fmt.Println("Contains empty query error")
	}

	fmt.Println("\n=== 5. Multi-Field Errors ===")
	err = validateUser("AB", "", 16)
	if err != nil {
		fmt.Print(err.Error())
	}

	err = validateUser("Alice Johnson", "alice@example.com", 25)
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Println("User validation passed!")
	}

	fmt.Println("\n=== 6. Retryable Error ===")
	for i := 1; i <= 4; i++ {
		err := fetchData(i)
		if err != nil {
			fmt.Printf("Attempt %d: %v\n", i, err)
			if IsRetryable(err) {
				fmt.Println("  → Will retry...")
			}
		} else {
			fmt.Printf("Attempt %d: Success!\n", i)
			break
		}
	}

	fmt.Println("\n=== 7. errors.As() ===")
	demonstrateErrorsAs()

	fmt.Println("\n=== 8. Permission Error with Methods ===")
	err = checkPermission("john", "database", "delete")
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Use type assertion to access custom methods
		if pe, ok := err.(PermissionError); ok {
			fmt.Printf("Suggestion: %s\n", pe.SuggestAction())
		}
	}

	fmt.Println("\n=== 9. Error Builder Pattern ===")
	complexErr := NewErrorBuilder().
		WithCode(500).
		WithMessage("Failed to process request").
		WithDetail("Database connection timeout").
		WithDetail("Retry attempt 3 of 3 failed").
		WithCause(errors.New("network unreachable")).
		Build()

	fmt.Printf("Complex error:\n%v\n", complexErr)

	fmt.Println("\n=== CUSTOM ERROR BEST PRACTICES ===")
	fmt.Println("✓ Implement Error() string method")
	fmt.Println("✓ Include relevant context (field names, IDs, timestamps)")
	fmt.Println("✓ Implement Unwrap() for wrapped errors")
	fmt.Println("✓ Use errors.As() for type-safe error checking")
	fmt.Println("✓ Add custom methods for error-specific logic")
	fmt.Println("✓ Make error types comparable when possible")
	fmt.Println("✓ Export error types for external use")
	fmt.Println("✓ Consider error codes for APIs")
	fmt.Println("✗ Don't include sensitive data in error messages")
	fmt.Println("✗ Don't make errors too generic")
}
