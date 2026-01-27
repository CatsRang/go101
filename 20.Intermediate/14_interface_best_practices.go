package main

import (
	"fmt"
	"time"
)

// ============================================================================
// INTERFACE BEST PRACTICES IN GO
// ============================================================================

// ============================================================================
// 1. ACCEPT INTERFACES, RETURN CONCRETE TYPES
// ============================================================================

// ❌ BAD: Returning interface
type BadDataStore interface {
	Save(key string, value interface{}) error
	Load(key string) (interface{}, error)
}

func NewBadDataStore() BadDataStore {
	return &MemoryStore{data: make(map[string]interface{})}
}

// ✓ GOOD: Return concrete type
type MemoryStore struct {
	data map[string]interface{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]interface{})}
}

func (ms *MemoryStore) Save(key string, value interface{}) error {
	ms.data[key] = value
	return nil
}

func (ms *MemoryStore) Load(key string) (interface{}, error) {
	if val, ok := ms.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}

// ============================================================================
// 2. KEEP INTERFACES SMALL
// ============================================================================

// ❌ BAD: Large interface with many methods (God interface)
type BadUserService interface {
	CreateUser(name, email string) error
	UpdateUser(id int, name, email string) error
	DeleteUser(id int) error
	GetUser(id int) (User, error)
	ListUsers() ([]User, error)
	AuthenticateUser(email, password string) (bool, error)
	ChangePassword(id int, oldPass, newPass string) error
	SendEmail(id int, subject, body string) error
	LogActivity(id int, activity string) error
}

// ✓ GOOD: Small, focused interfaces
type UserCreator interface {
	CreateUser(name, email string) error
}

type UserReader interface {
	GetUser(id int) (User, error)
	ListUsers() ([]User, error)
}

type UserUpdater interface {
	UpdateUser(id int, name, email string) error
}

type UserDeleter interface {
	DeleteUser(id int) error
}

type UserAuthenticator interface {
	AuthenticateUser(email, password string) (bool, error)
}

// Compose when needed
type UserManager interface {
	UserCreator
	UserReader
	UserUpdater
	UserDeleter
}

type User struct {
	ID    int
	Name  string
	Email string
}

// ============================================================================
// 3. DEFINE INTERFACES AT THE POINT OF USE
// ============================================================================

// ❌ BAD: Defining interface with implementation
// file: database.go
type Database struct {
	conn string
}

type DatabaseOperations interface { // Don't define this here!
	Query(sql string) ([]Row, error)
	Execute(sql string) error
}

func (db *Database) Query(sql string) ([]Row, error) {
	return nil, nil
}

func (db *Database) Execute(sql string) error {
	return nil
}

type Row struct {
	Data map[string]interface{}
}

// ✓ GOOD: Define interface where it's used
// file: userservice.go (consumer)

// UserService only needs Query method, so it defines what it needs
type Querier interface {
	Query(sql string) ([]Row, error)
}

type UserService struct {
	db Querier // Depends on interface, not concrete Database
}

func NewUserService(db Querier) *UserService {
	return &UserService{db: db}
}

func (us *UserService) GetUsers() ([]User, error) {
	us.db.Query("SELECT * FROM users")
	// ... implementation
	return nil, nil
}

// ============================================================================
// 4. DESIGN FOR TESTABILITY
// ============================================================================

// ✓ Interface makes testing easy
type EmailSender interface {
	Send(to, subject, body string) error
}

// Production implementation
type SMTPEmailSender struct {
	host string
	port int
}

func (ses *SMTPEmailSender) Send(to, subject, body string) error {
	fmt.Printf("Sending email to %s via SMTP\n", to)
	return nil
}

// Mock implementation for testing
type MockEmailSender struct {
	SentEmails []EmailRecord
}

type EmailRecord struct {
	To      string
	Subject string
	Body    string
}

func (mes *MockEmailSender) Send(to, subject, body string) error {
	mes.SentEmails = append(mes.SentEmails, EmailRecord{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	fmt.Printf("Mock: Recording email to %s\n", to)
	return nil
}

// Service that uses EmailSender
type NotificationService struct {
	emailSender EmailSender
}

func NewNotificationService(sender EmailSender) *NotificationService {
	return &NotificationService{emailSender: sender}
}

func (ns *NotificationService) NotifyUser(email, message string) error {
	return ns.emailSender.Send(email, "Notification", message)
}

// ============================================================================
// 5. INTERFACE SEGREGATION PRINCIPLE (ISP)
// ============================================================================

// ✓ GOOD: Clients should not depend on methods they don't use

type Reader interface {
	Read() (string, error)
}

type Writer interface {
	Write(data string) error
}

// ReadOnlyService only needs Reader
type ReadOnlyService struct {
	reader Reader
}

func NewReadOnlyService(r Reader) *ReadOnlyService {
	return &ReadOnlyService{reader: r}
}

func (ros *ReadOnlyService) Process() {
	data, _ := ros.reader.Read()
	fmt.Printf("Processing: %s\n", data)
}

// WriteOnlyService only needs Writer
type WriteOnlyService struct {
	writer Writer
}

func NewWriteOnlyService(w Writer) *WriteOnlyService {
	return &WriteOnlyService{writer: w}
}

func (wos *WriteOnlyService) Save(data string) {
	wos.writer.Write(data)
}

// ============================================================================
// 6. AVOID EMPTY INTERFACE WHEN POSSIBLE
// ============================================================================

// ❌ BAD: Overuse of interface{}
func BadProcess(data interface{}) interface{} {
	// No type safety!
	return data
}

// ✓ GOOD: Use concrete types or generics (Go 1.18+)
func GoodProcess(data string) string {
	return data
}

// Or use generics (Go 1.18+)
// func GenericProcess[T any](data T) T {
//     return data
// }

// ============================================================================
// 7. NIL INTERFACE VALUES
// ============================================================================

// Demonstrating nil interface gotcha
type Logger interface {
	Log(message string)
}

type FileLogger struct {
	filename string
}

func (fl *FileLogger) Log(message string) {
	if fl == nil {
		fmt.Println("Nil logger, can't log:", message)
		return
	}
	fmt.Printf("[%s] %s\n", fl.filename, message)
}

func GetLogger(useFile bool) Logger {
	if useFile {
		return &FileLogger{filename: "app.log"}
	}
	// ⚠️ DANGER: Returning typed nil
	var logger *FileLogger = nil
	return logger // This is NOT a nil interface!
}

// ============================================================================
// 8. COMPOSITION OVER LARGE INTERFACES
// ============================================================================

// Instead of one large interface, compose smaller ones

type Validator interface {
	Validate() error
}

type Saver interface {
	Save() error
}

type Loader interface {
	Load() error
}

// Compose as needed
type Persistable interface {
	Validator
	Saver
	Loader
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return nil
}

func (p *Product) Save() error {
	fmt.Printf("Saving product: %s\n", p.Name)
	return nil
}

func (p *Product) Load() error {
	fmt.Printf("Loading product: %d\n", p.ID)
	return nil
}

// ============================================================================
// 9. INTERFACE NAMING CONVENTIONS
// ============================================================================

// ✓ Single-method interfaces often end in "-er"
type Runner interface {
	Run() error
}

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Closer interface {
	Close() error
}

// ✓ Describes behavior, not data
type Processor interface {
	Process(data []byte) error
}

// ❌ BAD: Naming based on data/implementation
// type UserData interface { ... }
// type DatabaseInterface interface { ... }

func main() {
	fmt.Println("=== 1. Accept Interfaces, Return Concrete Types ===")

	// Return concrete type
	store := NewMemoryStore()
	store.Save("key1", "value1")

	// Accept interface
	var dataStore interface {
		Save(key string, value interface{}) error
		Load(key string) (interface{}, error)
	} = store

	dataStore.Save("key2", "value2")
	val, _ := dataStore.Load("key2")
	fmt.Printf("Loaded: %v\n", val)

	fmt.Println("\n=== 2. Small vs Large Interfaces ===")

	fmt.Println("✓ Small interfaces are more flexible")
	fmt.Println("✓ Easier to implement and test")
	fmt.Println("✓ Can be composed when needed")
	fmt.Println("✗ Large interfaces are rigid and hard to implement")

	fmt.Println("\n=== 3. Testability with Interfaces ===")

	// Production code
	productionSender := &SMTPEmailSender{host: "smtp.example.com", port: 587}
	productionService := NewNotificationService(productionSender)
	productionService.NotifyUser("user@example.com", "Hello from production")

	// Test code
	mockSender := &MockEmailSender{SentEmails: make([]EmailRecord, 0)}
	testService := NewNotificationService(mockSender)
	testService.NotifyUser("test@example.com", "Hello from test")

	fmt.Printf("Mock recorded %d emails\n", len(mockSender.SentEmails))
	if len(mockSender.SentEmails) > 0 {
		fmt.Printf("First email to: %s\n", mockSender.SentEmails[0].To)
	}

	fmt.Println("\n=== 4. Interface Segregation ===")

	// FileSystem implements both Reader and Writer
	type FileSystem struct {
		data string
	}

	fs := &FileSystem{data: "file contents"}

	fs.Read = func() (string, error) { return fs.data, nil }
	fs.Write = func(data string) error { fs.data = data; return nil }

	// But services only depend on what they need
	readService := NewReadOnlyService(fs)
	readService.Process()

	writeService := NewWriteOnlyService(fs)
	writeService.Save("new data")

	fmt.Println("\n=== 5. Nil Interface Gotcha ===")

	logger1 := GetLogger(true)
	if logger1 != nil {
		logger1.Log("This works")
	}

	logger2 := GetLogger(false)
	// ⚠️ logger2 != nil (it's a typed nil), but the underlying value is nil!
	if logger2 != nil {
		fmt.Println("Logger2 is not nil (but underlying value is!)")
		logger2.Log("This might not work as expected")
	}

	fmt.Println("\n=== 6. Composition Example ===")

	product := &Product{
		ID:    1,
		Name:  "Laptop",
		Price: 999.99,
	}

	// Use as Validator
	var validator Validator = product
	if err := validator.Validate(); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Validation passed")
	}

	// Use as Persistable
	var persistable Persistable = product
	persistable.Save()

	fmt.Println("\n=== INTERFACE BEST PRACTICES SUMMARY ===")
	fmt.Println("\n1. Accept interfaces, return concrete types")
	fmt.Println("   • Makes your functions flexible")
	fmt.Println("   • Callers can pass any implementation")
	fmt.Println("")
	fmt.Println("2. Keep interfaces small (1-3 methods)")
	fmt.Println("   • Easier to implement")
	fmt.Println("   • More composable")
	fmt.Println("   • Better testability")
	fmt.Println("")
	fmt.Println("3. Define interfaces where they're used")
	fmt.Println("   • Consumer defines what it needs")
	fmt.Println("   • Avoids unnecessary dependencies")
	fmt.Println("   • Follows Dependency Inversion Principle")
	fmt.Println("")
	fmt.Println("4. Use interfaces for testability")
	fmt.Println("   • Easy to create mocks")
	fmt.Println("   • Test without external dependencies")
	fmt.Println("   • Faster and more reliable tests")
	fmt.Println("")
	fmt.Println("5. Follow Interface Segregation Principle")
	fmt.Println("   • Clients depend only on methods they use")
	fmt.Println("   • Split large interfaces into smaller ones")
	fmt.Println("")
	fmt.Println("6. Avoid interface{}/any when possible")
	fmt.Println("   • Prefer concrete types")
	fmt.Println("   • Use generics (Go 1.18+) when appropriate")
	fmt.Println("   • Keep type safety")
	fmt.Println("")
	fmt.Println("7. Be aware of nil interface values")
	fmt.Println("   • An interface with a nil underlying value != nil")
	fmt.Println("   • Always check before dereferencing")
	fmt.Println("")
	fmt.Println("8. Name interfaces descriptively")
	fmt.Println("   • Single-method: use -er suffix (Reader, Writer)")
	fmt.Println("   • Describe behavior, not implementation")

	time.Sleep(10 * time.Millisecond)
}
