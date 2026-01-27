package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// MULTIPLE ERROR HANDLING (Go 1.20+)
// ============================================================================
// errors.Join() allows combining multiple errors into one
// fmt.Errorf() with multiple %w also wraps multiple errors

// ============================================================================
// 1. BASIC errors.Join()
// ============================================================================

func demonstrateBasicJoin() {
	err1 := errors.New("first error")
	err2 := errors.New("second error")
	err3 := errors.New("third error")

	// Join multiple errors into one
	multiErr := errors.Join(err1, err2, err3)

	fmt.Printf("Combined error: %v\n", multiErr)
	fmt.Println("\nEach error on separate line:")
	fmt.Println(multiErr)
}

// ============================================================================
// 2. VALIDATION WITH MULTIPLE ERRORS
// ============================================================================

type User struct {
	Name  string
	Email string
	Age   int
}

func validateUser(u User) error {
	var errs []error

	if u.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if len(u.Name) < 3 {
		errs = append(errs, errors.New("name must be at least 3 characters"))
	}
	if u.Email == "" {
		errs = append(errs, errors.New("email is required"))
	}
	if !strings.Contains(u.Email, "@") {
		errs = append(errs, errors.New("email must contain @"))
	}
	if u.Age < 18 {
		errs = append(errs, errors.New("age must be 18 or older"))
	}
	if u.Age > 150 {
		errs = append(errs, errors.New("age is unrealistic"))
	}

	// Join all validation errors
	return errors.Join(errs...)
}

// ============================================================================
// 3. CHECKING INDIVIDUAL ERRORS WITH errors.Is()
// ============================================================================

var (
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
	ErrAgeInvalid    = errors.New("age is invalid")
)

func validateUserWithSentinels(name, email string, age int) error {
	var errs []error

	if name == "" {
		errs = append(errs, ErrNameRequired)
	}
	if email == "" {
		errs = append(errs, ErrEmailRequired)
	}
	if age < 0 || age > 150 {
		errs = append(errs, ErrAgeInvalid)
	}

	return errors.Join(errs...)
}

func demonstrateErrorsIs() {
	err := validateUserWithSentinels("", "", -1)

	// errors.Is() works with joined errors
	if errors.Is(err, ErrNameRequired) {
		fmt.Println("✓ Name is required")
	}
	if errors.Is(err, ErrEmailRequired) {
		fmt.Println("✓ Email is required")
	}
	if errors.Is(err, ErrAgeInvalid) {
		fmt.Println("✓ Age is invalid")
	}
}

// ============================================================================
// 4. MULTIPLE %w IN fmt.Errorf() (Go 1.20+)
// ============================================================================

func processData() error {
	dbErr := errors.New("database connection failed")
	cacheErr := errors.New("cache update failed")

	// Go 1.20+: Multiple %w in single fmt.Errorf
	return fmt.Errorf("data processing failed: %w, %w", dbErr, cacheErr)
}

func demonstrateMultipleWrap() {
	err := processData()
	fmt.Printf("Error: %v\n", err)

	// Check for both wrapped errors
	if errors.Is(err, errors.New("database connection failed")) {
		fmt.Println("Contains database error")
	}
	if errors.Is(err, errors.New("cache update failed")) {
		fmt.Println("Contains cache error")
	}
}

// ============================================================================
// 5. FILE PROCESSING WITH MULTIPLE ERRORS
// ============================================================================

type FileError struct {
	Filename string
	Err      error
}

func (fe *FileError) Error() string {
	return fmt.Sprintf("file %q: %v", fe.Filename, fe.Err)
}

func (fe *FileError) Unwrap() error {
	return fe.Err
}

func processFiles(filenames []string) error {
	var errs []error

	for _, filename := range filenames {
		// Simulate file processing
		if strings.HasSuffix(filename, ".bad") {
			errs = append(errs, &FileError{
				Filename: filename,
				Err:      errors.New("corrupted file"),
			})
		}
		if strings.HasSuffix(filename, ".missing") {
			errs = append(errs, &FileError{
				Filename: filename,
				Err:      errors.New("file not found"),
			})
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to process %d files: %w", len(errs), errors.Join(errs...))
	}
	return nil
}

// ============================================================================
// 6. PARALLEL OPERATIONS WITH MULTIPLE ERRORS
// ============================================================================

func fetchFromAPI(endpoint string) error {
	if strings.Contains(endpoint, "users") {
		return errors.New("users API timeout")
	}
	if strings.Contains(endpoint, "products") {
		return errors.New("products API unavailable")
	}
	return nil
}

func fetchAllData() error {
	endpoints := []string{
		"/api/users",
		"/api/products",
		"/api/orders",
	}

	var errs []error
	for _, endpoint := range endpoints {
		if err := fetchFromAPI(endpoint); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", endpoint, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// ============================================================================
// 7. CUSTOM MULTI-ERROR TYPE
// ============================================================================

type MultiError struct {
	Errors []error
}

func (me *MultiError) Error() string {
	if len(me.Errors) == 0 {
		return "no errors"
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d errors occurred:\n", len(me.Errors)))
	for i, err := range me.Errors {
		builder.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
	}
	return builder.String()
}

// Unwrap returns all errors (Go 1.20+ supports this)
func (me *MultiError) Unwrap() []error {
	return me.Errors
}

func processWithCustomMultiError() error {
	return &MultiError{
		Errors: []error{
			errors.New("validation failed"),
			errors.New("authorization failed"),
			errors.New("rate limit exceeded"),
		},
	}
}

// ============================================================================
// 8. HANDLING JOINED ERRORS
// ============================================================================

func handleJoinedErrors(err error) {
	if err == nil {
		fmt.Println("No errors")
		return
	}

	fmt.Println("Handling joined errors:")

	// Method 1: Print all errors at once
	fmt.Printf("All errors: %v\n", err)

	// Method 2: Check for specific errors
	if errors.Is(err, ErrNameRequired) {
		fmt.Println("  → Please provide a name")
	}
	if errors.Is(err, ErrEmailRequired) {
		fmt.Println("  → Please provide an email")
	}
	if errors.Is(err, ErrAgeInvalid) {
		fmt.Println("  → Please provide a valid age")
	}
}

// ============================================================================
// 9. BATCH OPERATIONS WITH ERROR COLLECTION
// ============================================================================

type Task struct {
	ID   int
	Name string
}

func executeTask(task Task) error {
	if task.ID%2 == 0 {
		return fmt.Errorf("task %d (%s) failed", task.ID, task.Name)
	}
	return nil
}

func executeBatch(tasks []Task) error {
	var errs []error
	successCount := 0

	for _, task := range tasks {
		if err := executeTask(task); err != nil {
			errs = append(errs, err)
		} else {
			successCount++
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d tasks succeeded, %d failed: %w",
			successCount, len(errs), errors.Join(errs...))
	}
	return nil
}

// ============================================================================
// 10. CLEANUP WITH MULTIPLE ERRORS
// ============================================================================

type Resource struct {
	Name   string
	Closed bool
}

func (r *Resource) Close() error {
	if r.Closed {
		return fmt.Errorf("resource %s already closed", r.Name)
	}
	if r.Name == "bad" {
		return fmt.Errorf("failed to close %s", r.Name)
	}
	r.Closed = true
	return nil
}

func closeAllResources(resources []*Resource) error {
	var errs []error

	for _, res := range resources {
		if err := res.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	// Return all close errors together
	return errors.Join(errs...)
}

func main() {
	fmt.Println("=== 1. Basic errors.Join() ===")
	demonstrateBasicJoin()

	fmt.Println("\n=== 2. Validation with Multiple Errors ===")
	user := User{Name: "AB", Email: "invalid", Age: 16}
	if err := validateUser(user); err != nil {
		fmt.Println("Validation errors:")
		fmt.Println(err)
	}

	// Valid user
	validUser := User{Name: "Alice Johnson", Email: "alice@example.com", Age: 25}
	if err := validateUser(validUser); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("User is valid!")
	}

	fmt.Println("\n=== 3. errors.Is() with Joined Errors ===")
	demonstrateErrorsIs()

	fmt.Println("\n=== 4. Multiple %w in fmt.Errorf() ===")
	demonstrateMultipleWrap()

	fmt.Println("\n=== 5. File Processing ===")
	files := []string{
		"data1.txt",
		"data2.bad",
		"data3.txt",
		"data4.missing",
		"data5.bad",
	}

	if err := processFiles(files); err != nil {
		fmt.Printf("File processing error:\n%v\n", err)
	}

	fmt.Println("\n=== 6. Parallel API Fetching ===")
	if err := fetchAllData(); err != nil {
		fmt.Println("API fetch errors:")
		fmt.Println(err)
	}

	fmt.Println("\n=== 7. Custom Multi-Error Type ===")
	if err := processWithCustomMultiError(); err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println("\n=== 8. Handling Joined Errors ===")
	err := validateUserWithSentinels("", "", 200)
	handleJoinedErrors(err)

	fmt.Println("\n=== 9. Batch Task Execution ===")
	tasks := []Task{
		{ID: 1, Name: "Task One"},
		{ID: 2, Name: "Task Two"},
		{ID: 3, Name: "Task Three"},
		{ID: 4, Name: "Task Four"},
		{ID: 5, Name: "Task Five"},
	}

	if err := executeBatch(tasks); err != nil {
		fmt.Printf("Batch execution:\n%v\n", err)
	}

	fmt.Println("\n=== 10. Resource Cleanup ===")
	resources := []*Resource{
		{Name: "db1", Closed: false},
		{Name: "db2", Closed: false},
		{Name: "bad", Closed: false},
		{Name: "cache", Closed: false},
	}

	if err := closeAllResources(resources); err != nil {
		fmt.Println("Close errors:")
		fmt.Println(err)
	}

	fmt.Println("\n=== MULTIPLE ERROR HANDLING BEST PRACTICES ===")
	fmt.Println("✓ Use errors.Join() to combine multiple errors (Go 1.20+)")
	fmt.Println("✓ Use multiple %w in fmt.Errorf() (Go 1.20+)")
	fmt.Println("✓ Collect all validation errors before returning")
	fmt.Println("✓ Continue processing to find all errors, not just the first")
	fmt.Println("✓ errors.Is() works with joined errors")
	fmt.Println("✓ Implement Unwrap() []error for custom multi-error types")
	fmt.Println("✓ Include count of successful vs failed operations")
	fmt.Println("✓ Use for cleanup operations (close all, even if some fail)")
	fmt.Println("✗ Don't join errors unnecessarily (single error is simpler)")
	fmt.Println("✗ Don't lose information by merging too aggressively")

	fmt.Println("\n=== WHEN TO USE MULTIPLE ERRORS ===")
	fmt.Println("Use when:")
	fmt.Println("  • Validating multiple fields")
	fmt.Println("  • Processing batch operations")
	fmt.Println("  • Cleaning up multiple resources")
	fmt.Println("  • Parallel operations that can all fail")
	fmt.Println("  • You want to report all problems, not just the first")
	fmt.Println("\nAvoid when:")
	fmt.Println("  • Errors are sequential (use wrapping instead)")
	fmt.Println("  • Only one error can occur at a time")
	fmt.Println("  • First error should stop processing")
}
