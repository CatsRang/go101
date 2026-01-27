package main

import (
	"fmt"
	"io"
	"strings"
)

// ============================================================================
// STANDARD LIBRARY INTERFACES
// ============================================================================
// Go's standard library defines powerful, composable interfaces that are
// used throughout the ecosystem.

// ============================================================================
// 1. io.Reader - One of the Most Important Interfaces
// ============================================================================
// type Reader interface {
//     Read(p []byte) (n int, err error)
// }

// StringReader implements io.Reader for a string
type StringReader struct {
	data string
	pos  int
}

func NewStringReader(s string) *StringReader {
	return &StringReader{data: s, pos: 0}
}

// Read implements io.Reader
func (sr *StringReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}

	n = copy(p, sr.data[sr.pos:])
	sr.pos += n
	return n, nil
}

// CountingReader wraps an io.Reader and counts bytes read
type CountingReader struct {
	reader    io.Reader
	bytesRead int64
}

func NewCountingReader(r io.Reader) *CountingReader {
	return &CountingReader{reader: r}
}

func (cr *CountingReader) Read(p []byte) (n int, err error) {
	n, err = cr.reader.Read(p)
	cr.bytesRead += int64(n)
	return n, err
}

func (cr *CountingReader) BytesRead() int64 {
	return cr.bytesRead
}

// ============================================================================
// 2. io.Writer - Writing Data
// ============================================================================
// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

// MemoryWriter implements io.Writer to write to memory
type MemoryWriter struct {
	data []byte
}

func NewMemoryWriter() *MemoryWriter {
	return &MemoryWriter{data: make([]byte, 0)}
}

// Write implements io.Writer
func (mw *MemoryWriter) Write(p []byte) (n int, err error) {
	mw.data = append(mw.data, p...)
	return len(p), nil
}

func (mw *MemoryWriter) String() string {
	return string(mw.data)
}

// PrefixWriter adds a prefix to each write
type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{writer: w, prefix: prefix}
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	prefixed := append([]byte(pw.prefix), p...)
	return pw.writer.Write(prefixed)
}

// ============================================================================
// 3. error Interface - Error Handling
// ============================================================================
// type error interface {
//     Error() string
// }

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", ve.Field, ve.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", nfe.Resource, nfe.ID)
}

// MultiError represents multiple errors
type MultiError struct {
	Errors []error
}

func (me MultiError) Error() string {
	if len(me.Errors) == 0 {
		return "no errors"
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	msg := fmt.Sprintf("%d errors occurred:", len(me.Errors))
	for i, err := range me.Errors {
		msg += fmt.Sprintf("\n  %d. %s", i+1, err.Error())
	}
	return msg
}

// ============================================================================
// 4. fmt.Stringer - String Representation
// ============================================================================
// type Stringer interface {
//     String() string
// }

// User implements fmt.Stringer
type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
}

// String implements fmt.Stringer
func (u User) String() string {
	status := "inactive"
	if u.Active {
		status = "active"
	}
	return fmt.Sprintf("User{ID: %d, Username: %s, Email: %s, Status: %s}",
		u.ID, u.Username, u.Email, status)
}

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// String implements fmt.Stringer
func (p Point) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", p.X, p.Y)
}

// ============================================================================
// 5. io.Closer - Resource Cleanup
// ============================================================================
// type Closer interface {
//     Close() error
// }

// Connection represents a network connection
type Connection struct {
	host   string
	port   int
	closed bool
}

func NewConnection(host string, port int) *Connection {
	return &Connection{host: host, port: port, closed: false}
}

// Close implements io.Closer
func (c *Connection) Close() error {
	if c.closed {
		return fmt.Errorf("connection already closed")
	}
	c.closed = true
	fmt.Printf("Connection to %s:%d closed\n", c.host, c.port)
	return nil
}

func (c *Connection) String() string {
	status := "open"
	if c.closed {
		status = "closed"
	}
	return fmt.Sprintf("Connection{%s:%d, status: %s}", c.host, c.port, status)
}

// ============================================================================
// 6. Combining Standard Interfaces
// ============================================================================

// ReadWriteCloser implements multiple standard interfaces
type ReadWriteCloser struct {
	data   []byte
	pos    int
	closed bool
}

func NewReadWriteCloser() *ReadWriteCloser {
	return &ReadWriteCloser{
		data: make([]byte, 0),
	}
}

// Read implements io.Reader
func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	if rwc.closed {
		return 0, fmt.Errorf("reader is closed")
	}
	if rwc.pos >= len(rwc.data) {
		return 0, io.EOF
	}
	n = copy(p, rwc.data[rwc.pos:])
	rwc.pos += n
	return n, nil
}

// Write implements io.Writer
func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	if rwc.closed {
		return 0, fmt.Errorf("writer is closed")
	}
	rwc.data = append(rwc.data, p...)
	return len(p), nil
}

// Close implements io.Closer
func (rwc *ReadWriteCloser) Close() error {
	if rwc.closed {
		return fmt.Errorf("already closed")
	}
	rwc.closed = true
	return nil
}

// String implements fmt.Stringer
func (rwc *ReadWriteCloser) String() string {
	status := "open"
	if rwc.closed {
		status = "closed"
	}
	return fmt.Sprintf("ReadWriteCloser{size: %d bytes, status: %s}", len(rwc.data), status)
}

// ============================================================================
// 7. Using Standard Library Functions
// ============================================================================

// CopyData demonstrates using io.Copy with io.Reader and io.Writer
func CopyData(src io.Reader, dst io.Writer) (int64, error) {
	return io.Copy(dst, src)
}

// ReadAll reads all data from a reader
func ReadAll(r io.Reader) ([]byte, error) {
	buf := make([]byte, 0, 512)
	for {
		if len(buf) == cap(buf) {
			// Grow buffer
			newBuf := make([]byte, len(buf), 2*cap(buf)+1)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := r.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]

		if err != nil {
			if err == io.EOF {
				return buf, nil
			}
			return buf, err
		}
	}
}

func main() {
	fmt.Println("=== 1. io.Reader Interface ===")

	// Using custom StringReader
	reader := NewStringReader("Hello, Go interfaces!")
	buf := make([]byte, 10)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
		}
		if err == io.EOF {
			fmt.Println("Reached end of data")
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
	}

	fmt.Println("\n=== 2. Counting Reader Wrapper ===")

	// Wrap strings.Reader with CountingReader
	originalReader := strings.NewReader("Count these bytes!")
	countingReader := NewCountingReader(originalReader)

	// Read all data
	data, _ := ReadAll(countingReader)
	fmt.Printf("Read: %s\n", string(data))
	fmt.Printf("Total bytes read: %d\n", countingReader.BytesRead())

	fmt.Println("\n=== 3. io.Writer Interface ===")

	writer := NewMemoryWriter()

	// Write data
	writer.Write([]byte("Hello, "))
	writer.Write([]byte("Writer "))
	writer.Write([]byte("Interface!"))

	fmt.Printf("Written data: %s\n", writer.String())

	fmt.Println("\n=== 4. Prefix Writer Wrapper ===")

	memWriter := NewMemoryWriter()
	prefixWriter := NewPrefixWriter(memWriter, "[LOG] ")

	prefixWriter.Write([]byte("First log message\n"))
	prefixWriter.Write([]byte("Second log message\n"))

	fmt.Printf("Output:\n%s", memWriter.String())

	fmt.Println("\n=== 5. error Interface ===")

	// Different error types
	errors := []error{
		ValidationError{Field: "email", Message: "invalid format"},
		NotFoundError{Resource: "User", ID: 123},
		fmt.Errorf("generic error"),
		MultiError{Errors: []error{
			fmt.Errorf("first error"),
			fmt.Errorf("second error"),
			fmt.Errorf("third error"),
		}},
	}

	for i, err := range errors {
		fmt.Printf("%d. %s\n", i+1, err.Error())
	}

	fmt.Println("\n=== 6. fmt.Stringer Interface ===")

	user := User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Active:   true,
	}

	point := Point{X: 3.14, Y: 2.71}

	// fmt.Println uses the String() method if available
	fmt.Println("User:", user)
	fmt.Println("Point:", point)

	// Direct call
	fmt.Println("User.String():", user.String())
	fmt.Println("Point.String():", point.String())

	fmt.Println("\n=== 7. io.Closer Interface ===")

	conn := NewConnection("localhost", 8080)
	fmt.Println(conn)

	// Close the connection
	err := conn.Close()
	if err != nil {
		fmt.Printf("Error closing: %v\n", err)
	}
	fmt.Println(conn)

	// Try to close again
	err = conn.Close()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== 8. Combining Multiple Interfaces ===")

	rwc := NewReadWriteCloser()
	fmt.Println("Initial:", rwc)

	// Write data
	rwc.Write([]byte("Hello, "))
	rwc.Write([]byte("combined "))
	rwc.Write([]byte("interfaces!"))
	fmt.Println("After writing:", rwc)

	// Read data back
	readBuf := make([]byte, 50)
	n, _ := rwc.Read(readBuf)
	fmt.Printf("Read: %s\n", string(readBuf[:n]))

	// Close
	rwc.Close()
	fmt.Println("After closing:", rwc)

	fmt.Println("\n=== 9. Using with Standard Library ===")

	// io.Copy works with any io.Reader and io.Writer
	source := strings.NewReader("Data to copy via io.Copy")
	destination := NewMemoryWriter()

	bytesCopied, _ := CopyData(source, destination)
	fmt.Printf("Copied %d bytes: %s\n", bytesCopied, destination.String())

	// Using strings.Reader (standard library io.Reader)
	stdReader := strings.NewReader("Standard library reader")
	customWriter := NewMemoryWriter()

	io.Copy(customWriter, stdReader)
	fmt.Printf("Result: %s\n", customWriter.String())

	fmt.Println("\n=== Standard Interfaces Summary ===")
	fmt.Println("\nio.Reader:")
	fmt.Println("  Read(p []byte) (n int, err error)")
	fmt.Println("  → Used for reading data from various sources")
	fmt.Println("  → Examples: files, network, strings, bytes")
	fmt.Println("\nio.Writer:")
	fmt.Println("  Write(p []byte) (n int, err error)")
	fmt.Println("  → Used for writing data to various destinations")
	fmt.Println("  → Examples: files, network, buffers")
	fmt.Println("\nerror:")
	fmt.Println("  Error() string")
	fmt.Println("  → The only required method for error types")
	fmt.Println("  → Enable custom error types with additional context")
	fmt.Println("\nfmt.Stringer:")
	fmt.Println("  String() string")
	fmt.Println("  → Controls how your type is printed")
	fmt.Println("  → Used by fmt.Print, fmt.Println, etc.")
	fmt.Println("\nio.Closer:")
	fmt.Println("  Close() error")
	fmt.Println("  → Used for resource cleanup")
	fmt.Println("  → Often combined with defer")

	fmt.Println("\n=== Why These Interfaces Are Powerful ===")
	fmt.Println("✓ Small and focused (1 method each)")
	fmt.Println("✓ Composable (io.ReadWriter, io.ReadCloser, etc.)")
	fmt.Println("✓ Widely adopted across the ecosystem")
	fmt.Println("✓ Enable polymorphism without inheritance")
	fmt.Println("✓ Work with standard library functions")
	fmt.Println("✓ Make testing easier (easy to mock)")
	fmt.Println("✓ Encourage good API design")
}
