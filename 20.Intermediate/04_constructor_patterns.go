package main

import (
	"fmt"
	"time"
)

// User demonstrates basic constructor pattern
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
}

// NewUser is a basic constructor function
func NewUser(id int, username, email string) *User {
	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
}

// Database demonstrates constructor with validation
type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

// NewDatabase creates a database with validation
func NewDatabase(host string, port int, username, password, dbName string) (*Database, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	return &Database{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Name:     dbName,
	}, nil
}

// Config demonstrates functional options pattern
type Config struct {
	Host           string
	Port           int
	Timeout        time.Duration
	MaxConnections int
	EnableLogging  bool
	RetryAttempts  int
}

// Option is a function that configures Config
type Option func(*Config)

// NewConfig creates a Config with default values and applies options
func NewConfig(options ...Option) *Config {
	// Set default values
	cfg := &Config{
		Host:           "localhost",
		Port:           8080,
		Timeout:        30 * time.Second,
		MaxConnections: 100,
		EnableLogging:  true,
		RetryAttempts:  3,
	}

	// Apply all options
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

// WithHost sets the host
func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

// WithPort sets the port
func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithTimeout sets the timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithMaxConnections sets max connections
func WithMaxConnections(max int) Option {
	return func(c *Config) {
		c.MaxConnections = max
	}
}

// WithLogging enables or disables logging
func WithLogging(enable bool) Option {
	return func(c *Config) {
		c.EnableLogging = enable
	}
}

// WithRetryAttempts sets retry attempts
func WithRetryAttempts(attempts int) Option {
	return func(c *Config) {
		c.RetryAttempts = attempts
	}
}

// Server demonstrates builder pattern
type Server struct {
	host      string
	port      int
	timeout   time.Duration
	tlsConfig *TLSConfig
	logger    *Logger
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// Logger represents logging configuration
type Logger struct {
	Enabled bool
	Level   string
}

// ServerBuilder builds a Server
type ServerBuilder struct {
	server *Server
}

// NewServerBuilder creates a new ServerBuilder
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		server: &Server{
			host:    "localhost",
			port:    8080,
			timeout: 30 * time.Second,
		},
	}
}

// Host sets the host
func (b *ServerBuilder) Host(host string) *ServerBuilder {
	b.server.host = host
	return b
}

// Port sets the port
func (b *ServerBuilder) Port(port int) *ServerBuilder {
	b.server.port = port
	return b
}

// Timeout sets the timeout
func (b *ServerBuilder) Timeout(timeout time.Duration) *ServerBuilder {
	b.server.timeout = timeout
	return b
}

// WithTLS configures TLS
func (b *ServerBuilder) WithTLS(certFile, keyFile string) *ServerBuilder {
	b.server.tlsConfig = &TLSConfig{
		Enabled:  true,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	return b
}

// WithLogger configures logging
func (b *ServerBuilder) WithLogger(level string) *ServerBuilder {
	b.server.logger = &Logger{
		Enabled: true,
		Level:   level,
	}
	return b
}

// Build creates the Server
func (b *ServerBuilder) Build() *Server {
	return b.server
}

// ConnectionPool demonstrates factory pattern with pooling
type ConnectionPool struct {
	maxConnections int
	connections    []string
}

// NewConnectionPool creates a connection pool
func NewConnectionPool(max int) *ConnectionPool {
	pool := &ConnectionPool{
		maxConnections: max,
		connections:    make([]string, 0, max),
	}

	// Initialize connections
	for i := 0; i < max; i++ {
		pool.connections = append(pool.connections,
			fmt.Sprintf("conn-%d", i))
	}

	return pool
}

// Account demonstrates zero value usability pattern
type Account struct {
	balance float64
	locked  bool
}

// NewAccount creates a new account with initial balance
func NewAccount(initialBalance float64) *Account {
	return &Account{
		balance: initialBalance,
		locked:  false,
	}
}

// Deposit adds money (works even with zero value)
func (a *Account) Deposit(amount float64) {
	if !a.locked && amount > 0 {
		a.balance += amount
	}
}

// GetBalance returns balance (works even with zero value)
func (a *Account) GetBalance() float64 {
	return a.balance
}

func main() {
	fmt.Println("=== Basic Constructor Pattern ===")

	user1 := NewUser(1, "johndoe", "john@example.com")
	fmt.Printf("User: %+v\n", user1)

	user2 := NewUser(2, "janedoe", "jane@example.com")
	fmt.Printf("User: %+v\n", user2)

	fmt.Println("\n=== Constructor with Validation ===")

	db1, err := NewDatabase("localhost", 5432, "admin", "secret", "mydb")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Database: %+v\n", db1)
	}

	// Invalid database configuration
	db2, err := NewDatabase("", 5432, "admin", "secret", "mydb")
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
	}

	db3, err := NewDatabase("localhost", -1, "admin", "secret", "mydb")
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
	}

	fmt.Println("\n=== Functional Options Pattern ===")

	// Use defaults
	config1 := NewConfig()
	fmt.Printf("Default config: %+v\n", config1)

	// Override some options
	config2 := NewConfig(
		WithHost("api.example.com"),
		WithPort(443),
		WithTimeout(60*time.Second),
	)
	fmt.Printf("Custom config: %+v\n", config2)

	// Override all options
	config3 := NewConfig(
		WithHost("db.example.com"),
		WithPort(5432),
		WithTimeout(120*time.Second),
		WithMaxConnections(200),
		WithLogging(false),
		WithRetryAttempts(5),
	)
	fmt.Printf("Fully custom config: %+v\n", config3)

	fmt.Println("\n=== Builder Pattern ===")

	// Simple server
	server1 := NewServerBuilder().
		Host("example.com").
		Port(8080).
		Build()
	fmt.Printf("Simple server: %+v\n", server1)

	// Server with TLS and logging
	server2 := NewServerBuilder().
		Host("secure.example.com").
		Port(443).
		Timeout(60 * time.Second).
		WithTLS("/path/to/cert.pem", "/path/to/key.pem").
		WithLogger("DEBUG").
		Build()
	fmt.Printf("Secure server: %+v\n", server2)

	fmt.Println("\n=== Factory Pattern ===")

	pool := NewConnectionPool(5)
	fmt.Printf("Connection pool: %+v\n", pool)
	fmt.Printf("Pool size: %d\n", len(pool.connections))

	fmt.Println("\n=== Zero Value Usability ===")

	// Constructor with initial value
	account1 := NewAccount(1000.00)
	fmt.Printf("Account 1 balance: $%.2f\n", account1.GetBalance())
	account1.Deposit(500.00)
	fmt.Printf("Account 1 after deposit: $%.2f\n", account1.GetBalance())

	// Zero value is also usable
	var account2 Account // Zero value: balance=0, locked=false
	fmt.Printf("Account 2 (zero value) balance: $%.2f\n", account2.GetBalance())
	account2.Deposit(250.00)
	fmt.Printf("Account 2 after deposit: $%.2f\n", account2.GetBalance())

	fmt.Println("\n=== Constructor Pattern Guidelines ===")
	fmt.Println("1. Basic Constructor (NewX):")
	fmt.Println("   - Simple initialization")
	fmt.Println("   - Returns pointer (*X) for mutable types")
	fmt.Println("")
	fmt.Println("2. Constructor with Validation:")
	fmt.Println("   - Returns (*X, error)")
	fmt.Println("   - Validates inputs before creation")
	fmt.Println("")
	fmt.Println("3. Functional Options Pattern:")
	fmt.Println("   - Many optional parameters")
	fmt.Println("   - Good defaults, easy to extend")
	fmt.Println("   - Clean API for users")
	fmt.Println("")
	fmt.Println("4. Builder Pattern:")
	fmt.Println("   - Complex object construction")
	fmt.Println("   - Step-by-step configuration")
	fmt.Println("   - Fluent interface (method chaining)")
	fmt.Println("")
	fmt.Println("5. Zero Value Usability:")
	fmt.Println("   - Design structs to work with zero values")
	fmt.Println("   - Makes constructors optional")
	fmt.Println("   - Simplifies usage")
}
