package main

import (
	"fmt"
	"time"
)

// Person is a base struct
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// FullName returns the person's full name
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// Introduce prints an introduction
func (p Person) Introduce() {
	fmt.Printf("Hi, I'm %s and I'm %d years old.\n", p.FullName(), p.Age)
}

// Employee embeds Person (composition, not inheritance)
type Employee struct {
	Person          // Embedded struct (anonymous field)
	EmployeeID      string
	Department      string
	Salary          float64
	HireDate        time.Time
}

// GetDetails returns employee details
func (e Employee) GetDetails() string {
	return fmt.Sprintf("Employee ID: %s, Name: %s, Dept: %s",
		e.EmployeeID, e.FullName(), e.Department)
}

// Introduce overrides the embedded Person's Introduce method
func (e Employee) Introduce() {
	fmt.Printf("Hello, I'm %s, working in %s department.\n",
		e.FullName(), e.Department)
}

// Student embeds Person
type Student struct {
	Person
	StudentID string
	Major     string
	GPA       float64
}

// GetInfo returns student information
func (s Student) GetInfo() string {
	return fmt.Sprintf("Student: %s, Major: %s, GPA: %.2f",
		s.FullName(), s.Major, s.GPA)
}

// Address struct for composition
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// FullAddress returns formatted address
func (a Address) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s", a.Street, a.City, a.State, a.ZipCode)
}

// Contact struct with multiple embedded types
type Contact struct {
	Person
	Address
	Phone string
	Email string
}

// Engine component
type Engine struct {
	Type       string
	Horsepower int
	FuelType   string
}

// Start starts the engine
func (e Engine) Start() {
	fmt.Printf("%s engine starting... (HP: %d)\n", e.Type, e.Horsepower)
}

// Wheels component
type Wheels struct {
	Count int
	Size  int
}

// Car composes Engine and Wheels
type Car struct {
	Brand  string
	Model  string
	Year   int
	Engine // Embedded
	Wheels // Embedded
}

// Drive demonstrates using embedded components
func (c Car) Drive() {
	fmt.Printf("Driving %d %s %s with %d wheels\n",
		c.Year, c.Brand, c.Model, c.Wheels.Count)
}

// Logger provides logging functionality
type Logger struct {
	Prefix string
}

// Log prints a log message
func (l Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Prefix, message)
}

// Database provides database functionality
type Database struct {
	ConnectionString string
}

// Connect simulates database connection
func (d Database) Connect() {
	fmt.Printf("Connecting to database: %s\n", d.ConnectionString)
}

// Query simulates database query
func (d Database) Query(sql string) {
	fmt.Printf("Executing query: %s\n", sql)
}

// Service composes Logger and Database
type Service struct {
	Name string
	Logger
	Database
}

// Initialize initializes the service
func (s Service) Initialize() {
	s.Log("Service initializing: " + s.Name)
	s.Connect()
	s.Log("Service initialized successfully")
}

// Example with name collision
type A struct {
	Name string
}

func (a A) GetName() string {
	return "A: " + a.Name
}

type B struct {
	Name string
}

func (b B) GetName() string {
	return "B: " + b.Name
}

// C embeds both A and B (name collision)
type C struct {
	A
	B
	Value int
}

func main() {
	fmt.Println("=== Basic Struct Embedding ===")

	emp := Employee{
		Person: Person{
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
		},
		EmployeeID: "EMP001",
		Department: "Engineering",
		Salary:     75000.00,
		HireDate:   time.Now(),
	}

	// Accessing embedded struct fields directly
	fmt.Printf("Name: %s\n", emp.FullName())
	fmt.Printf("Age: %d\n", emp.Age) // Direct access to embedded field
	fmt.Printf("Department: %s\n", emp.Department)

	// Calling methods from embedded struct
	emp.Introduce() // Uses Employee's overridden Introduce

	// Can still call the embedded Person's method explicitly
	emp.Person.Introduce()

	fmt.Println("\n=== Student Example ===")

	student := Student{
		Person: Person{
			FirstName: "Jane",
			LastName:  "Smith",
			Age:       20,
		},
		StudentID: "STU12345",
		Major:     "Computer Science",
		GPA:       3.8,
	}

	fmt.Println(student.GetInfo())
	student.Introduce() // Uses embedded Person's Introduce

	fmt.Println("\n=== Multiple Embedded Structs ===")

	contact := Contact{
		Person: Person{
			FirstName: "Alice",
			LastName:  "Johnson",
			Age:       28,
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94102",
		},
		Phone: "555-0123",
		Email: "alice@example.com",
	}

	// Access fields and methods from both embedded structs
	fmt.Printf("Contact: %s\n", contact.FullName())
	fmt.Printf("Address: %s\n", contact.FullAddress())
	fmt.Printf("Phone: %s, Email: %s\n", contact.Phone, contact.Email)

	// Can also access embedded struct fields directly
	fmt.Printf("City: %s\n", contact.City) // From embedded Address

	fmt.Println("\n=== Composition for Functionality ===")

	car := Car{
		Brand: "Tesla",
		Model: "Model 3",
		Year:  2024,
		Engine: Engine{
			Type:       "Electric",
			Horsepower: 283,
			FuelType:   "Electric",
		},
		Wheels: Wheels{
			Count: 4,
			Size:  19,
		},
	}

	fmt.Printf("Car: %d %s %s\n", car.Year, car.Brand, car.Model)
	car.Engine.Start() // Explicit access to embedded method
	car.Start()        // Direct access to embedded method
	car.Drive()

	fmt.Println("\n=== Service with Multiple Components ===")

	service := Service{
		Name: "UserService",
		Logger: Logger{
			Prefix: "USER_SVC",
		},
		Database: Database{
			ConnectionString: "postgres://localhost:5432/mydb",
		},
	}

	service.Initialize()
	service.Query("SELECT * FROM users")
	service.Log("Query completed")

	fmt.Println("\n=== Handling Name Collisions ===")

	c := C{
		A:     A{Name: "ValueA"},
		B:     B{Name: "ValueB"},
		Value: 42,
	}

	// When there's a name collision, must use explicit path
	fmt.Printf("A.Name: %s\n", c.A.Name)
	fmt.Printf("B.Name: %s\n", c.B.Name)

	// Methods also need explicit path when there's collision
	fmt.Println(c.A.GetName())
	fmt.Println(c.B.GetName())

	// This would cause compile error (ambiguous):
	// fmt.Println(c.Name)      // Error: ambiguous selector c.Name
	// fmt.Println(c.GetName()) // Error: ambiguous selector c.GetName

	fmt.Println("\n=== Composition vs Inheritance ===")
	fmt.Println("Go Philosophy: 'Composition over inheritance'")
	fmt.Println("Benefits:")
	fmt.Println("  - More flexible and explicit")
	fmt.Println("  - Avoids complex inheritance hierarchies")
	fmt.Println("  - Easier to understand and maintain")
	fmt.Println("  - No fragile base class problem")
	fmt.Println("  - Can compose multiple behaviors easily")
}
