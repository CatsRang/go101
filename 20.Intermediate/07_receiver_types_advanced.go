package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 1. INTERFACE COMPLIANCE: Pointer vs Value Receivers
// ============================================================================

// Counter demonstrates interface compliance
type Counter struct {
	count int
}

// IncrementPointer uses pointer receiver - can modify state
func (c *Counter) IncrementPointer() {
	c.count++
}

// GetCount returns current count
func (c *Counter) GetCount() int {
	return c.count
}

// IncrementValue uses value receiver - modifies a copy
func (c Counter) IncrementValue() {
	c.count++ // This only modifies the copy!
}

// Incrementer interface for demonstrating interface compliance
type Incrementer interface {
	IncrementPointer()
	GetCount() int
}

// ============================================================================
// 2. MONITORED READER: The Classic Use Case
// ============================================================================

// MonitoredReader wraps an io.Reader to track metrics
// This MUST use pointer receiver to satisfy io.Reader interface
type MonitoredReader struct {
	reader    io.Reader
	bytesRead int64
	readCalls int
	mu        sync.Mutex // MUST NOT be copied - requires pointer receiver
}

// NewMonitoredReader creates a new monitored reader
func NewMonitoredReader(r io.Reader) *MonitoredReader {
	return &MonitoredReader{
		reader: r,
	}
}

// Read implements io.Reader interface with POINTER RECEIVER
// This is required because:
// 1. io.Reader interface expects pointer receiver
// 2. We need to track state (bytesRead, readCalls)
// 3. sync.Mutex cannot be copied safely
func (r *MonitoredReader) Read(buf []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, err = r.reader.Read(buf)
	r.bytesRead += int64(n)
	r.readCalls++

	return n, err
}

// GetStats returns reading statistics
func (r *MonitoredReader) GetStats() (bytesRead int64, calls int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.bytesRead, r.readCalls
}

// BrokenMonitoredReader demonstrates what happens with value receiver
type BrokenMonitoredReader struct {
	reader    io.Reader
	bytesRead int64
	readCalls int
}

// Read with VALUE RECEIVER - This is broken!
// Problem: Changes to bytesRead and readCalls don't persist
func (r BrokenMonitoredReader) Read(buf []byte) (n int, err error) {
	n, err = r.reader.Read(buf)
	r.bytesRead += int64(n) // Only modifies the copy!
	r.readCalls++           // Only modifies the copy!
	return n, err
}

// ============================================================================
// 3. MEMORY EFFICIENCY: Large Structs
// ============================================================================

// SmallStruct is tiny (16 bytes) - value receiver is fine
type SmallStruct struct {
	ID   int64
	Flag bool
}

// ProcessValue with value receiver - copies 16 bytes
func (s SmallStruct) ProcessValue() string {
	return fmt.Sprintf("SmallStruct ID: %d", s.ID)
}

// ProcessPointer with pointer receiver - copies 8 bytes (pointer)
func (s *SmallStruct) ProcessPointer() string {
	return fmt.Sprintf("SmallStruct ID: %d", s.ID)
}

// LargeStruct is big (>1KB) - pointer receiver is essential
type LargeStruct struct {
	Data     [1024]byte // 1KB array
	Metadata map[string]string
	History  []string
	Config   struct {
		Settings [100]int64
		Flags    [100]bool
	}
}

// ProcessValue with value receiver - copies >1KB on every call!
func (l LargeStruct) ProcessValue() int {
	return len(l.Data)
}

// ProcessPointer with pointer receiver - copies only 8 bytes
func (l *LargeStruct) ProcessPointer() int {
	return len(l.Data)
}

// ============================================================================
// 4. MUTABILITY AND STATE TRACKING
// ============================================================================

// RateLimiter tracks request rates with state
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	mu          sync.Mutex
}

// NewRateLimiter creates a rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0),
	}
}

// Allow checks if request is allowed - MUST use pointer receiver
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Remove old requests outside the window
	cutoff := now.Add(-rl.window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range rl.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	rl.requests = validRequests

	// Check if we can allow this request
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}

// GetRequestCount returns current request count
func (rl *RateLimiter) GetRequestCount() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return len(rl.requests)
}

// ============================================================================
// 5. SYNC.MUTEX AND UNSAFE COPYING
// ============================================================================

// SafeCounter uses pointer receiver with mutex
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// Increment safely increments with pointer receiver
func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.value++
}

// GetValue safely gets value
func (sc *SafeCounter) GetValue() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.value
}

// UnsafeCounter demonstrates the problem with value receiver
type UnsafeCounter struct {
	mu    sync.Mutex
	value int
}

// BrokenIncrement with value receiver - DANGEROUS!
// This would copy the mutex, which is unsafe
// Go vet will warn about this
func (uc UnsafeCounter) BrokenIncrement() {
	uc.mu.Lock()   // Locking a COPY of the mutex!
	defer uc.mu.Unlock()
	uc.value++ // Modifying a COPY of the value!
}

// ============================================================================
// 6. INTERFACE SATISFACTION EXAMPLES
// ============================================================================

// Writer interface that expects pointer receiver
type Writer interface {
	Write(data []byte) (int, error)
}

// BufferWriter implements Writer
type BufferWriter struct {
	buffer bytes.Buffer
	writes int
}

// Write implements Writer interface - MUST use pointer receiver
func (bw *BufferWriter) Write(data []byte) (int, error) {
	bw.writes++
	return bw.buffer.Write(data)
}

// GetWrites returns write count
func (bw *BufferWriter) GetWrites() int {
	return bw.writes
}

// GetContent returns buffer content
func (bw *BufferWriter) GetContent() string {
	return bw.buffer.String()
}

// ============================================================================
// 7. PRACTICAL DECISION GUIDE
// ============================================================================

// Point represents a small, immutable point - value receiver is fine
type Point struct {
	X, Y float64
}

// Distance calculates distance from origin - value receiver
func (p Point) Distance() float64 {
	return p.X*p.X + p.Y*p.Y
}

// Add returns new point - value receiver for immutable semantics
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Circle represents a mutable circle - pointer receiver
type Circle struct {
	Center Point
	Radius float64
}

// Grow increases radius - pointer receiver to modify
func (c *Circle) Grow(delta float64) {
	c.Radius += delta
}

// Area calculates area - pointer receiver for consistency
func (c *Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// ============================================================================
// MAIN: DEMONSTRATIONS
// ============================================================================

func main() {
	fmt.Println("=== 1. Interface Compliance ===")

	counter := &Counter{count: 0}

	// This works because Counter methods use pointer receiver
	var inc Incrementer = counter
	inc.IncrementPointer()
	inc.IncrementPointer()
	fmt.Printf("Counter value: %d\n", inc.GetCount())

	// Value receiver doesn't modify original
	counter.IncrementValue()
	fmt.Printf("After IncrementValue: %d (no change!)\n", counter.GetCount())

	fmt.Println("\n=== 2. MonitoredReader: Correct Implementation ===")

	data := "Hello, World! This is a test of the monitored reader."
	reader := NewMonitoredReader(strings.NewReader(data))

	// Satisfies io.Reader interface
	var _ io.Reader = reader

	// Read some data
	buf := make([]byte, 10)
	n, _ := reader.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))

	// Read more data
	n, _ = reader.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))

	// Check stats - state was tracked correctly
	bytesRead, calls := reader.GetStats()
	fmt.Printf("Total bytes read: %d, Total calls: %d\n", bytesRead, calls)

	fmt.Println("\n=== 3. BrokenMonitoredReader: Value Receiver Problem ===")

	brokenReader := BrokenMonitoredReader{
		reader: strings.NewReader(data),
	}

	// Read data
	buf2 := make([]byte, 10)
	brokenReader.Read(buf2)
	brokenReader.Read(buf2)

	// Stats are ALWAYS ZERO because Read() operates on a copy!
	fmt.Printf("Broken reader - bytes tracked: %d (should be >0!)\n", brokenReader.bytesRead)
	fmt.Printf("Broken reader - calls tracked: %d (should be 2!)\n", brokenReader.readCalls)

	fmt.Println("\n=== 4. Memory Efficiency Comparison ===")

	small := SmallStruct{ID: 42, Flag: true}
	fmt.Printf("SmallStruct size: ~16 bytes\n")
	fmt.Printf("  Value receiver: %s\n", small.ProcessValue())
	fmt.Printf("  Pointer receiver: %s\n", small.ProcessPointer())
	fmt.Println("  For small structs, either approach is fine")

	large := &LargeStruct{}
	large.Data[0] = 'A'
	fmt.Printf("\nLargeStruct size: >1KB\n")
	fmt.Printf("  Pointer receiver result: %d\n", large.ProcessPointer())
	fmt.Println("  ⚠️  Value receiver would copy >1KB on every call!")
	fmt.Println("  Always use pointer receiver for large structs")

	fmt.Println("\n=== 5. State Tracking with RateLimiter ===")

	limiter := NewRateLimiter(3, 1*time.Second)

	fmt.Println("Attempting 5 requests (limit: 3 per second):")
	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow()
		status := "✓ ALLOWED"
		if !allowed {
			status = "✗ DENIED"
		}
		fmt.Printf("  Request %d: %s (total: %d)\n", i, status, limiter.GetRequestCount())
	}

	fmt.Println("\n=== 6. Mutex Safety ===")

	safeCounter := &SafeCounter{}

	// Safe concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			safeCounter.Increment()
		}()
	}
	wg.Wait()

	fmt.Printf("Safe counter (pointer receiver): %d (correct!)\n", safeCounter.GetValue())

	fmt.Println("⚠️  UnsafeCounter.BrokenIncrement would copy the mutex!")
	fmt.Println("    Go vet would warn: 'passes lock by value'")

	fmt.Println("\n=== 7. Interface Satisfaction ===")

	writer := &BufferWriter{}

	// Satisfies Writer interface
	var w Writer = writer

	w.Write([]byte("Hello "))
	w.Write([]byte("World!"))

	fmt.Printf("Writes: %d\n", writer.GetWrites())
	fmt.Printf("Content: %s\n", writer.GetContent())

	fmt.Println("\n=== 8. Immutable vs Mutable Semantics ===")

	// Point: small, immutable - value receiver
	p1 := Point{X: 3, Y: 4}
	p2 := Point{X: 1, Y: 2}
	p3 := p1.Add(p2)
	fmt.Printf("Point p1: %+v (unchanged)\n", p1)
	fmt.Printf("Point p3: %+v (new point)\n", p3)

	// Circle: mutable - pointer receiver
	circle := &Circle{Center: Point{X: 0, Y: 0}, Radius: 5}
	fmt.Printf("Circle before: radius=%.2f, area=%.2f\n", circle.Radius, circle.Area())
	circle.Grow(3)
	fmt.Printf("Circle after:  radius=%.2f, area=%.2f\n", circle.Radius, circle.Area())

	fmt.Println("\n=== DECISION GUIDE ===")
	fmt.Println("┌─────────────────────────────────────────────┬──────────┬───────┐")
	fmt.Println("│ Scenario                                    │ Pointer  │ Value │")
	fmt.Println("├─────────────────────────────────────────────┼──────────┼───────┤")
	fmt.Println("│ Implements interface (io.Reader, etc.)     │    ✓     │   ✗   │")
	fmt.Println("│ Method modifies receiver fields            │    ✓     │   ✗   │")
	fmt.Println("│ Contains sync.Mutex or channels            │    ✓     │   ✗   │")
	fmt.Println("│ Large struct (>128 bytes)                  │    ✓     │   ✗   │")
	fmt.Println("│ Small struct + immutable                   │  Either  │   ✓   │")
	fmt.Println("│ Reader/Writer wrapper pattern              │    ✓     │   ✗   │")
	fmt.Println("│ Needs to track state/metrics               │    ✓     │   ✗   │")
	fmt.Println("│ Performance-critical path                  │    ✓     │   -   │")
	fmt.Println("└─────────────────────────────────────────────┴──────────┴───────┘")

	fmt.Println("\n=== KEY TAKEAWAYS ===")
	fmt.Println("1. Use pointer receiver when method needs to modify the receiver")
	fmt.Println("2. Use pointer receiver for structs with sync.Mutex, channels, or other non-copyable fields")
	fmt.Println("3. Use pointer receiver for large structs to avoid expensive copies")
	fmt.Println("4. Use pointer receiver when implementing interfaces (usually)")
	fmt.Println("5. Use value receiver for small, immutable types (like Point)")
	fmt.Println("6. For consistency: if any method needs pointer receiver, use it for all methods")
	fmt.Println("7. MonitoredReader pattern ALWAYS requires pointer receiver")
}
