package main

import (
	"fmt"
	"time"
)

// ============================================================================
// INTERFACE COMPOSITION: Building Larger Interfaces from Smaller Ones
// ============================================================================

// Reader interface (simplified io.Reader)
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer interface (simplified io.Writer)
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Closer interface (simplified io.Closer)
type Closer interface {
	Close() error
}

// ReadWriter combines Reader and Writer interfaces
type ReadWriter interface {
	Reader
	Writer
}

// ReadWriteCloser combines all three
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// File implements ReadWriteCloser
type File struct {
	name   string
	data   []byte
	pos    int
	closed bool
}

func NewFile(name string) *File {
	return &File{
		name: name,
		data: make([]byte, 0),
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.closed {
		return 0, fmt.Errorf("file is closed")
	}
	if f.pos >= len(f.data) {
		return 0, fmt.Errorf("EOF")
	}

	n = copy(p, f.data[f.pos:])
	f.pos += n
	return n, nil
}

func (f *File) Write(p []byte) (n int, err error) {
	if f.closed {
		return 0, fmt.Errorf("file is closed")
	}
	f.data = append(f.data, p...)
	return len(p), nil
}

func (f *File) Close() error {
	if f.closed {
		return fmt.Errorf("file already closed")
	}
	f.closed = true
	fmt.Printf("File '%s' closed\n", f.name)
	return nil
}

// ============================================================================
// BUILDING COMPLEX INTERFACES
// ============================================================================

// Logger interface
type Logger interface {
	Log(message string)
}

// MetricsCollector interface
type MetricsCollector interface {
	RecordMetric(name string, value float64)
}

// Notifier interface
type Notifier interface {
	Notify(event string)
}

// ObservableService combines multiple interfaces
type ObservableService interface {
	Logger
	MetricsCollector
	Notifier
}

// Service implements ObservableService
type Service struct {
	name string
}

func (s Service) Log(message string) {
	fmt.Printf("[LOG][%s] %s\n", s.name, message)
}

func (s Service) RecordMetric(name string, value float64) {
	fmt.Printf("[METRIC][%s] %s: %.2f\n", s.name, name, value)
}

func (s Service) Notify(event string) {
	fmt.Printf("[NOTIFY][%s] Event: %s\n", s.name, event)
}

// ============================================================================
// INTERFACE SEGREGATION
// ============================================================================

// Database operations split into focused interfaces

// Readable allows reading operations
type Readable interface {
	Get(key string) (string, error)
	List() []string
}

// Writable allows writing operations
type Writable interface {
	Set(key, value string) error
	Delete(key string) error
}

// Transactional allows transaction operations
type Transactional interface {
	Begin() error
	Commit() error
	Rollback() error
}

// FullDatabase combines all database operations
type FullDatabase interface {
	Readable
	Writable
	Transactional
}

// ReadOnlyDB implements only Readable
type ReadOnlyDB struct {
	data map[string]string
}

func NewReadOnlyDB() *ReadOnlyDB {
	return &ReadOnlyDB{
		data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
}

func (db *ReadOnlyDB) Get(key string) (string, error) {
	if val, exists := db.data[key]; exists {
		return val, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (db *ReadOnlyDB) List() []string {
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys
}

// FullDB implements FullDatabase
type FullDB struct {
	data        map[string]string
	inTransaction bool
	backup      map[string]string
}

func NewFullDB() *FullDB {
	return &FullDB{
		data: make(map[string]string),
	}
}

func (db *FullDB) Get(key string) (string, error) {
	if val, exists := db.data[key]; exists {
		return val, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (db *FullDB) List() []string {
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys
}

func (db *FullDB) Set(key, value string) error {
	db.data[key] = value
	fmt.Printf("Set %s = %s\n", key, value)
	return nil
}

func (db *FullDB) Delete(key string) error {
	delete(db.data, key)
	fmt.Printf("Deleted %s\n", key)
	return nil
}

func (db *FullDB) Begin() error {
	if db.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}
	db.backup = make(map[string]string)
	for k, v := range db.data {
		db.backup[k] = v
	}
	db.inTransaction = true
	fmt.Println("Transaction started")
	return nil
}

func (db *FullDB) Commit() error {
	if !db.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}
	db.backup = nil
	db.inTransaction = false
	fmt.Println("Transaction committed")
	return nil
}

func (db *FullDB) Rollback() error {
	if !db.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}
	db.data = db.backup
	db.backup = nil
	db.inTransaction = false
	fmt.Println("Transaction rolled back")
	return nil
}

// ============================================================================
// REAL-WORLD EXAMPLE: HTTP Handler
// ============================================================================

// Authenticator handles authentication
type Authenticator interface {
	Authenticate(token string) (bool, error)
}

// Authorizer handles authorization
type Authorizer interface {
	Authorize(user, resource string) bool
}

// RateLimiter handles rate limiting
type RateLimiter interface {
	Allow(identifier string) bool
}

// Middleware combines security-related interfaces
type Middleware interface {
	Authenticator
	Authorizer
	RateLimiter
}

// SecurityMiddleware implements Middleware
type SecurityMiddleware struct {
	validTokens  map[string]bool
	permissions  map[string][]string
	requestCount map[string]int
	maxRequests  int
}

func NewSecurityMiddleware(maxRequests int) *SecurityMiddleware {
	return &SecurityMiddleware{
		validTokens: map[string]bool{
			"token123": true,
			"token456": true,
		},
		permissions: map[string][]string{
			"admin": {"read", "write", "delete"},
			"user":  {"read"},
		},
		requestCount: make(map[string]int),
		maxRequests:  maxRequests,
	}
}

func (sm *SecurityMiddleware) Authenticate(token string) (bool, error) {
	if sm.validTokens[token] {
		return true, nil
	}
	return false, fmt.Errorf("invalid token")
}

func (sm *SecurityMiddleware) Authorize(user, resource string) bool {
	permissions, exists := sm.permissions[user]
	if !exists {
		return false
	}
	for _, perm := range permissions {
		if perm == resource {
			return true
		}
	}
	return false
}

func (sm *SecurityMiddleware) Allow(identifier string) bool {
	sm.requestCount[identifier]++
	return sm.requestCount[identifier] <= sm.maxRequests
}

// ============================================================================
// EMPTY INTERFACE IN COMPOSITION
// ============================================================================

// Serializer handles serialization
type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
}

// Deserializer handles deserialization
type Deserializer interface {
	Deserialize(data []byte) (interface{}, error)
}

// Codec combines serialization and deserialization
type Codec interface {
	Serializer
	Deserializer
}

// JSONCodec implements Codec
type JSONCodec struct{}

func (jc JSONCodec) Serialize(data interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("JSON: %v", data)), nil
}

func (jc JSONCodec) Deserialize(data []byte) (interface{}, error) {
	return string(data), nil
}

func main() {
	fmt.Println("=== 1. Basic Interface Composition ===")

	file := NewFile("test.txt")

	// File implements ReadWriteCloser
	var rwc ReadWriteCloser = file

	// Write some data
	rwc.Write([]byte("Hello, World!"))
	rwc.Write([]byte(" Go interfaces are awesome!"))

	// Read it back
	buf := make([]byte, 50)
	n, _ := rwc.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))

	// Close the file
	rwc.Close()

	fmt.Println("\n=== 2. Observable Service ===")

	service := Service{name: "UserService"}

	// Use as ObservableService
	var obs ObservableService = service

	obs.Log("Service initialized")
	obs.RecordMetric("users_count", 1250.0)
	obs.Notify("new_user_registered")

	// Can also use individual interfaces
	var logger Logger = service
	logger.Log("Using logger interface")

	var metrics MetricsCollector = service
	metrics.RecordMetric("response_time_ms", 45.2)

	fmt.Println("\n=== 3. Interface Segregation ===")

	// Read-only database - only implements Readable
	roDB := NewReadOnlyDB()
	fmt.Println("Read-only database:")
	val, _ := roDB.Get("key1")
	fmt.Printf("  key1 = %s\n", val)
	fmt.Printf("  All keys: %v\n", roDB.List())

	// Full database - implements all interfaces
	fullDB := NewFullDB()
	fmt.Println("\nFull database:")

	// Use as Writable
	var writable Writable = fullDB
	writable.Set("name", "Alice")
	writable.Set("age", "30")

	// Use as Readable
	var readable Readable = fullDB
	val, _ = readable.Get("name")
	fmt.Printf("  name = %s\n", val)

	// Use as Transactional
	var txn Transactional = fullDB
	txn.Begin()
	writable.Set("temp", "temporary")
	txn.Rollback()

	fmt.Println("\n=== 4. Security Middleware Composition ===")

	middleware := NewSecurityMiddleware(3)

	// Use as Authenticator
	var auth Authenticator = middleware
	isValid, err := auth.Authenticate("token123")
	fmt.Printf("Token authentication: %v (err: %v)\n", isValid, err)

	isValid, err = auth.Authenticate("invalid")
	fmt.Printf("Invalid token: %v (err: %v)\n", isValid, err)

	// Use as Authorizer
	var authz Authorizer = middleware
	canRead := authz.Authorize("admin", "read")
	canDelete := authz.Authorize("user", "delete")
	fmt.Printf("Admin can read: %v\n", canRead)
	fmt.Printf("User can delete: %v\n", canDelete)

	// Use as RateLimiter
	var limiter RateLimiter = middleware
	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow("user123")
		status := "allowed"
		if !allowed {
			status = "denied"
		}
		fmt.Printf("Request %d: %s\n", i, status)
	}

	fmt.Println("\n=== 5. Codec Composition ===")

	codec := JSONCodec{}

	// Use as Codec (combined interface)
	var c Codec = codec

	// Serialize
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	serialized, _ := c.Serialize(data)
	fmt.Printf("Serialized: %s\n", string(serialized))

	// Deserialize
	deserialized, _ := c.Deserialize(serialized)
	fmt.Printf("Deserialized: %v\n", deserialized)

	fmt.Println("\n=== 6. Function Accepting Composed Interfaces ===")

	processDatabase := func(db FullDatabase) {
		fmt.Println("\nProcessing with full database access:")
		db.Begin()
		db.Set("transaction_key", "transaction_value")
		keys := db.List()
		fmt.Printf("Keys during transaction: %v\n", keys)
		db.Commit()
	}

	processDatabase(fullDB)

	// This wouldn't work with ReadOnlyDB:
	// processDatabase(roDB) // ✗ Compile error!

	// But this would work (accepts only Readable):
	processReadOnly := func(db Readable) {
		fmt.Println("\nProcessing with read-only access:")
		keys := db.List()
		fmt.Printf("Available keys: %v\n", keys)
	}

	processReadOnly(roDB)     // ✓ Works
	processReadOnly(fullDB)   // ✓ Also works (FullDB implements Readable)

	fmt.Println("\n=== Benefits of Interface Composition ===")
	fmt.Println("✓ Build complex interfaces from simple ones")
	fmt.Println("✓ Promotes interface segregation (smaller, focused interfaces)")
	fmt.Println("✓ Enables flexible API design")
	fmt.Println("✓ Types can implement subsets of functionality")
	fmt.Println("✓ Consumers can request only what they need")
	fmt.Println("✓ Makes testing easier (mock smaller interfaces)")
	fmt.Println("✓ Follows SOLID principles (especially ISP)")

	time.Sleep(10 * time.Millisecond) // Small delay for output
}
