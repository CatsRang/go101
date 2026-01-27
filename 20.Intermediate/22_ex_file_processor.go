package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// PRACTICAL EXERCISE: FILE PROCESSOR WITH COMPREHENSIVE ERROR HANDLING
// ============================================================================
// This system demonstrates all Week 6 error handling concepts:
// - Custom error types
// - Error wrapping
// - Multiple error handling (errors.Join)
// - Panic/recover
// - Defer for cleanup

// ============================================================================
// 1. CUSTOM ERROR TYPES
// ============================================================================

// FileError represents file-related errors
type FileError struct {
	Filename  string
	Operation string
	Err       error
}

func (fe *FileError) Error() string {
	return fmt.Sprintf("file error during %s on '%s': %v",
		fe.Operation, fe.Filename, fe.Err)
}

func (fe *FileError) Unwrap() error {
	return fe.Err
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s=%v - %s",
		ve.Field, ve.Value, ve.Message)
}

// ProcessingError represents data processing errors
type ProcessingError struct {
	Line    int
	Content string
	Err     error
}

func (pe *ProcessingError) Error() string {
	return fmt.Sprintf("processing error at line %d (%q): %v",
		pe.Line, pe.Content, pe.Err)
}

func (pe *ProcessingError) Unwrap() error {
	return pe.Err
}

// Sentinel errors
var (
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrEmptyFile       = errors.New("file is empty")
	ErrPermissionError = errors.New("permission denied")
	ErrCorruptedData   = errors.New("corrupted data")
)

// ============================================================================
// 2. FILE STRUCTURE
// ============================================================================

// File simulates a file with content
type File struct {
	Name    string
	Content []string
	IsOpen  bool
	LineNum int
}

// NewFile creates a new file with content
func NewFile(name string, lines []string) *File {
	return &File{
		Name:    name,
		Content: lines,
		IsOpen:  false,
		LineNum: 0,
	}
}

// Open opens the file
func (f *File) Open() error {
	if f.Name == "" {
		return &ValidationError{
			Field:   "filename",
			Value:   f.Name,
			Message: "cannot be empty",
		}
	}

	if strings.HasSuffix(f.Name, ".missing") {
		return &FileError{
			Filename:  f.Name,
			Operation: "open",
			Err:       ErrFileNotFound,
		}
	}

	if strings.HasSuffix(f.Name, ".denied") {
		return &FileError{
			Filename:  f.Name,
			Operation: "open",
			Err:       ErrPermissionError,
		}
	}

	f.IsOpen = true
	fmt.Printf("✓ Opened file: %s\n", f.Name)
	return nil
}

// ReadLine reads the next line
func (f *File) ReadLine() (string, error) {
	if !f.IsOpen {
		return "", &FileError{
			Filename:  f.Name,
			Operation: "read",
			Err:       errors.New("file not open"),
		}
	}

	if f.LineNum >= len(f.Content) {
		return "", fmt.Errorf("EOF")
	}

	line := f.Content[f.LineNum]
	f.LineNum++
	return line, nil
}

// Close closes the file
func (f *File) Close() error {
	if !f.IsOpen {
		return nil // Already closed
	}
	f.IsOpen = false
	fmt.Printf("✓ Closed file: %s\n", f.Name)
	return nil
}

// ============================================================================
// 3. DATA PROCESSOR
// ============================================================================

// Record represents a processed record
type Record struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// Processor handles file processing
type Processor struct {
	Stats ProcessingStats
}

// ProcessingStats tracks processing statistics
type ProcessingStats struct {
	TotalFiles      int
	SuccessfulFiles int
	FailedFiles     int
	TotalRecords    int
	ValidRecords    int
	InvalidRecords  int
	Errors          []error
}

// NewProcessor creates a new processor
func NewProcessor() *Processor {
	return &Processor{
		Stats: ProcessingStats{
			Errors: make([]error, 0),
		},
	}
}

// ProcessFile processes a single file with full error handling
func (p *Processor) ProcessFile(file *File) (err error) {
	p.Stats.TotalFiles++

	// Defer for cleanup and error tracking
	defer func() {
		// Panic recovery
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during file processing: %v", r)
			p.Stats.FailedFiles++
			fmt.Printf("✗ Recovered from panic: %v\n", err)
		}

		// Close file (even if panic occurred)
		if file.IsOpen {
			file.Close()
		}

		// Track error
		if err != nil {
			p.Stats.Errors = append(p.Stats.Errors, err)
			p.Stats.FailedFiles++
		} else {
			p.Stats.SuccessfulFiles++
		}
	}()

	// Open file
	if err := file.Open(); err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	// Collect errors during processing
	var lineErrors []error

	// Process each line
	for {
		line, err := file.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		p.Stats.TotalRecords++

		// Process line with error handling
		if err := p.processLine(file.LineNum-1, line); err != nil {
			lineErrors = append(lineErrors, err)
			p.Stats.InvalidRecords++
		} else {
			p.Stats.ValidRecords++
		}
	}

	// Join all line errors if any occurred
	if len(lineErrors) > 0 {
		return errors.Join(lineErrors...)
	}

	return nil
}

// processLine processes a single line
func (p *Processor) processLine(lineNum int, content string) error {
	// Simulate corrupted data detection
	if strings.Contains(content, "CORRUPTED") {
		return &ProcessingError{
			Line:    lineNum,
			Content: content,
			Err:     ErrCorruptedData,
		}
	}

	// Simulate invalid format
	if strings.Contains(content, "INVALID") {
		return &ProcessingError{
			Line:    lineNum,
			Content: content,
			Err:     ErrInvalidFormat,
		}
	}

	// Simulate panic on specific content
	if strings.Contains(content, "PANIC") {
		panic("unexpected data encountered")
	}

	// Valid line
	fmt.Printf("  ✓ Processed line %d: %s\n", lineNum, content)
	return nil
}

// ProcessBatch processes multiple files
func (p *Processor) ProcessBatch(files []*File) error {
	fmt.Printf("\n=== Processing Batch of %d Files ===\n", len(files))

	for i, file := range files {
		fmt.Printf("\nFile %d/%d: %s\n", i+1, len(files), file.Name)

		err := p.ProcessFile(file)
		if err != nil {
			// Log error but continue processing other files
			fmt.Printf("Error processing %s:\n%v\n", file.Name, err)
		}
	}

	return nil
}

// PrintStats prints processing statistics
func (p *Processor) PrintStats() {
	fmt.Println("\n=== Processing Statistics ===")
	fmt.Printf("Total Files:      %d\n", p.Stats.TotalFiles)
	fmt.Printf("  Successful:     %d\n", p.Stats.SuccessfulFiles)
	fmt.Printf("  Failed:         %d\n", p.Stats.FailedFiles)
	fmt.Printf("\nTotal Records:    %d\n", p.Stats.TotalRecords)
	fmt.Printf("  Valid:          %d\n", p.Stats.ValidRecords)
	fmt.Printf("  Invalid:        %d\n", p.Stats.InvalidRecords)

	if len(p.Stats.Errors) > 0 {
		fmt.Printf("\nErrors Encountered: %d\n", len(p.Stats.Errors))
		for i, err := range p.Stats.Errors {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
	}
}

// ============================================================================
// 4. ERROR ANALYSIS
// ============================================================================

// AnalyzeErrors analyzes collected errors
func (p *Processor) AnalyzeErrors() {
	if len(p.Stats.Errors) == 0 {
		fmt.Println("\n=== No Errors to Analyze ===")
		return
	}

	fmt.Println("\n=== Error Analysis ===")

	var (
		fileErrors       int
		validationErrors int
		processingErrors int
		otherErrors      int
	)

	for _, err := range p.Stats.Errors {
		// Type assertion to categorize errors
		var fileErr *FileError
		var valErr *ValidationError
		var procErr *ProcessingError

		switch {
		case errors.As(err, &fileErr):
			fileErrors++
		case errors.As(err, &valErr):
			validationErrors++
		case errors.As(err, &procErr):
			processingErrors++
		default:
			otherErrors++
		}

		// Check for specific sentinel errors
		if errors.Is(err, ErrFileNotFound) {
			fmt.Println("  → File not found error detected")
		}
		if errors.Is(err, ErrPermissionError) {
			fmt.Println("  → Permission error detected")
		}
		if errors.Is(err, ErrCorruptedData) {
			fmt.Println("  → Corrupted data detected")
		}
	}

	fmt.Printf("\nError Categories:\n")
	fmt.Printf("  File Errors:       %d\n", fileErrors)
	fmt.Printf("  Validation Errors: %d\n", validationErrors)
	fmt.Printf("  Processing Errors: %d\n", processingErrors)
	fmt.Printf("  Other Errors:      %d\n", otherErrors)
}

// ============================================================================
// MAIN: DEMONSTRATION
// ============================================================================

func main() {
	fmt.Println("=== FILE PROCESSOR WITH COMPREHENSIVE ERROR HANDLING ===")
	fmt.Println("Demonstrating:")
	fmt.Println("  • Custom error types")
	fmt.Println("  • Error wrapping with %w")
	fmt.Println("  • Multiple error handling (errors.Join)")
	fmt.Println("  • Panic recovery")
	fmt.Println("  • Defer for cleanup")
	fmt.Println("  • Sentinel errors")

	// Create processor
	processor := NewProcessor()

	// Create test files with various scenarios
	files := []*File{
		// Successful file
		NewFile("users.txt", []string{
			"User 1: John Doe, john@example.com",
			"User 2: Jane Smith, jane@example.com",
			"User 3: Bob Johnson, bob@example.com",
		}),

		// File with invalid data
		NewFile("data.txt", []string{
			"Valid data line 1",
			"INVALID data line 2",
			"Valid data line 3",
			"CORRUPTED data line 4",
		}),

		// File that will panic
		NewFile("dangerous.txt", []string{
			"Normal line 1",
			"PANIC at line 2",
			"This won't be reached",
		}),

		// File not found
		NewFile("missing.missing", []string{}),

		// Permission denied
		NewFile("secret.denied", []string{}),

		// Empty but valid file
		NewFile("empty.txt", []string{}),
	}

	// Process batch
	processor.ProcessBatch(files)

	// Print statistics
	processor.PrintStats()

	// Analyze errors
	processor.AnalyzeErrors()

	fmt.Println("\n=== CONCEPTS DEMONSTRATED ===")
	fmt.Println("✓ Custom error types (FileError, ValidationError, ProcessingError)")
	fmt.Println("✓ Error wrapping with %w for context")
	fmt.Println("✓ Multiple error collection with errors.Join()")
	fmt.Println("✓ Sentinel errors for common cases")
	fmt.Println("✓ errors.Is() for checking wrapped errors")
	fmt.Println("✓ errors.As() for type-safe error extraction")
	fmt.Println("✓ Panic recovery in defer")
	fmt.Println("✓ Defer for resource cleanup")
	fmt.Println("✓ Named return values for error modification")
	fmt.Println("✓ Error statistics and analysis")
	fmt.Println("✓ Graceful degradation (continue on error)")

	fmt.Println("\n=== ERROR HANDLING PATTERNS USED ===")
	fmt.Println("1. Open → Defer Close pattern")
	fmt.Println("2. Defer with panic recovery")
	fmt.Println("3. Named return for error modification")
	fmt.Println("4. Error collection during iteration")
	fmt.Println("5. errors.Join() for multiple line errors")
	fmt.Println("6. Error wrapping at each layer")
	fmt.Println("7. Error type checking with errors.As()")
	fmt.Println("8. Sentinel error checking with errors.Is()")
	fmt.Println("9. Statistics tracking for observability")
	fmt.Println("10. Graceful error handling (continue processing)")

	time.Sleep(10 * time.Millisecond) // Small delay for output
}
