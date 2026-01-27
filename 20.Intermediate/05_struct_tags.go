package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// Product demonstrates basic JSON struct tags
type Product struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	InStock     bool     `json:"in_stock"`
	Tags        []string `json:"tags"`
}

// User demonstrates various JSON tag options
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never serialized
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Age       int       `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Book demonstrates XML struct tags
type Book struct {
	XMLName xml.Name `xml:"book"`
	ID      int      `xml:"id,attr"`
	Title   string   `xml:"title"`
	Author  string   `xml:"author"`
	ISBN    string   `xml:"isbn,attr"`
	Price   float64  `xml:"price"`
	Pages   int      `xml:"pages"`
}

// Address demonstrates nested structs with tags
type Address struct {
	Street  string `json:"street" xml:"street"`
	City    string `json:"city" xml:"city"`
	State   string `json:"state" xml:"state"`
	ZipCode string `json:"zip_code" xml:"zipCode"`
	Country string `json:"country,omitempty" xml:"country,omitempty"`
}

// Customer demonstrates embedded structs with tags
type Customer struct {
	ID      int     `json:"id" xml:"id,attr"`
	Name    string  `json:"name" xml:"name"`
	Email   string  `json:"email" xml:"email"`
	Address Address `json:"address" xml:"address"`
	Phone   string  `json:"phone,omitempty" xml:"phone,omitempty"`
}

// APIResponse demonstrates complex nested structures
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Config demonstrates multiple tag types
type Config struct {
	Host     string        `json:"host" yaml:"host" toml:"host"`
	Port     int           `json:"port" yaml:"port" toml:"port"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout" toml:"timeout"`
	Database DBConfig      `json:"database" yaml:"database" toml:"database"`
}

// DBConfig represents database configuration
type DBConfig struct {
	Driver   string `json:"driver" yaml:"driver" toml:"driver"`
	Host     string `json:"host" yaml:"host" toml:"host"`
	Port     int    `json:"port" yaml:"port" toml:"port"`
	Username string `json:"username" yaml:"username" toml:"username"`
	Password string `json:"-" yaml:"-" toml:"-"` // Never serialized
	DBName   string `json:"db_name" yaml:"db_name" toml:"db_name"`
}

// Employee demonstrates string/number type conversion tags
type Employee struct {
	ID         int     `json:"id,string"` // Convert int to string in JSON
	Name       string  `json:"name"`
	Salary     float64 `json:"salary,string"` // Convert float to string
	Department string  `json:"department"`
	IsActive   bool    `json:"is_active"`
}

// BlogPost demonstrates time formatting
type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Published time.Time `json:"published"`
	Tags      []string  `json:"tags,omitempty"`
	Views     int       `json:"views,omitempty"`
}

func main() {
	fmt.Println("=== Basic JSON Serialization ===")

	product := Product{
		ID:          101,
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       1299.99,
		InStock:     true,
		Tags:        []string{"electronics", "computers", "sale"},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(product)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}
	fmt.Printf("JSON: %s\n", string(jsonData))

	// Pretty print JSON
	jsonPretty, _ := json.MarshalIndent(product, "", "  ")
	fmt.Printf("\nPretty JSON:\n%s\n", string(jsonPretty))

	fmt.Println("\n=== JSON Unmarshaling ===")

	jsonInput := `{
		"id": 102,
		"name": "Smartphone",
		"description": "Latest model",
		"price": 899.99,
		"in_stock": false,
		"tags": ["electronics", "mobile"]
	}`

	var product2 Product
	err = json.Unmarshal([]byte(jsonInput), &product2)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}
	fmt.Printf("Unmarshaled product: %+v\n", product2)

	fmt.Println("\n=== JSON Tag Options ===")

	user := User{
		ID:        1,
		Username:  "johndoe",
		Email:     "john@example.com",
		Password:  "secret123", // Will not be serialized
		FirstName: "John",
		LastName:  "Doe",
		Age:       0, // Will be omitted due to omitempty
		CreatedAt: time.Now(),
		// UpdatedAt is zero value, will be omitted
	}

	userJSON, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("User JSON:\n%s\n", string(userJSON))
	fmt.Println("Note: Password is not included, Age and UpdatedAt are omitted")

	fmt.Println("\n=== XML Serialization ===")

	book := Book{
		ID:     1,
		Title:  "Learning Go",
		Author: "John Doe",
		ISBN:   "978-1234567890",
		Price:  49.99,
		Pages:  450,
	}

	xmlData, err := xml.MarshalIndent(book, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling XML: %v\n", err)
		return
	}
	fmt.Printf("XML:\n%s\n", string(xmlData))

	fmt.Println("\n=== Nested Structs with Tags ===")

	customer := Customer{
		ID:    1001,
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94102",
			Country: "USA",
		},
		Phone: "555-0123",
	}

	// JSON
	customerJSON, _ := json.MarshalIndent(customer, "", "  ")
	fmt.Printf("Customer JSON:\n%s\n", string(customerJSON))

	// XML
	customerXML, _ := xml.MarshalIndent(customer, "", "  ")
	fmt.Printf("\nCustomer XML:\n%s\n", string(customerXML))

	fmt.Println("\n=== API Response Pattern ===")

	// Success response
	successResp := APIResponse{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data: map[string]interface{}{
			"users": []string{"user1", "user2", "user3"},
			"count": 3,
		},
	}

	successJSON, _ := json.MarshalIndent(successResp, "", "  ")
	fmt.Printf("Success Response:\n%s\n", string(successJSON))

	// Error response
	errorResp := APIResponse{
		Status: "error",
		Error: &ErrorInfo{
			Code:    404,
			Message: "Resource not found",
			Details: "The requested user does not exist",
		},
	}

	errorJSON, _ := json.MarshalIndent(errorResp, "", "  ")
	fmt.Printf("\nError Response:\n%s\n", string(errorJSON))

	fmt.Println("\n=== String Type Conversion ===")

	employee := Employee{
		ID:         12345,
		Name:       "Bob Smith",
		Salary:     75000.50,
		Department: "Engineering",
		IsActive:   true,
	}

	empJSON, _ := json.MarshalIndent(employee, "", "  ")
	fmt.Printf("Employee (ID and Salary as strings):\n%s\n", string(empJSON))

	fmt.Println("\n=== Complex Configuration ===")

	config := Config{
		Host:    "api.example.com",
		Port:    8080,
		Timeout: 30 * time.Second,
		Database: DBConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Username: "admin",
			Password: "secret", // Will not be serialized
			DBName:   "myapp",
		},
	}

	configJSON, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config JSON:\n%s\n", string(configJSON))
	fmt.Println("Note: Database password is not included")

	fmt.Println("\n=== Blog Post Example ===")

	post := BlogPost{
		ID:        1,
		Title:     "Introduction to Go Struct Tags",
		Content:   "Struct tags are a powerful feature...",
		Author:    "Jane Doe",
		Published: time.Now(),
		Tags:      []string{"go", "programming", "tutorial"},
		Views:     0, // Will be omitted due to omitempty
	}

	postJSON, _ := json.MarshalIndent(post, "", "  ")
	fmt.Printf("Blog Post:\n%s\n", string(postJSON))

	fmt.Println("\n=== Struct Tag Guidelines ===")
	fmt.Println("Common JSON tag options:")
	fmt.Println("  json:\"field_name\"        - Rename field in JSON")
	fmt.Println("  json:\"-\"                 - Exclude field from JSON")
	fmt.Println("  json:\"field,omitempty\"   - Omit if zero value")
	fmt.Println("  json:\"field,string\"      - Convert to/from string")
	fmt.Println("")
	fmt.Println("Common XML tag options:")
	fmt.Println("  xml:\"element\"            - Element name")
	fmt.Println("  xml:\"field,attr\"         - XML attribute")
	fmt.Println("  xml:\"-\"                  - Exclude from XML")
	fmt.Println("  xml:\"field,omitempty\"    - Omit if zero value")
}
