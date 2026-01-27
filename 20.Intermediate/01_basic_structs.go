package main

import "fmt"

// Person struct defines a basic structure
type Person struct {
	FirstName string
	LastName  string
	Age       int
	Email     string
}

// Product struct demonstrates different field types
type Product struct {
	ID          int
	Name        string
	Price       float64
	InStock     bool
	Categories  []string
	Attributes  map[string]string
}

// Address struct shows nested structures
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// Employee struct with embedded Address
type Employee struct {
	ID      int
	Name    string
	Salary  float64
	Address Address // Nested struct
}

func main() {
	fmt.Println("=== Basic Struct Usage ===")

	// Creating a struct with field names
	person1 := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Email:     "john.doe@example.com",
	}
	fmt.Printf("Person 1: %+v\n", person1)

	// Creating a struct with positional values (not recommended)
	person2 := Person{"Jane", "Smith", 25, "jane.smith@example.com"}
	fmt.Printf("Person 2: %+v\n", person2)

	// Creating a struct with partial initialization
	person3 := Person{
		FirstName: "Bob",
		LastName:  "Johnson",
		// Age and Email will have zero values (0 and "")
	}
	fmt.Printf("Person 3 (partial): %+v\n", person3)

	// Zero value struct (all fields have their zero values)
	var person4 Person
	fmt.Printf("Person 4 (zero value): %+v\n", person4)

	fmt.Println("\n=== Accessing and Modifying Struct Fields ===")

	// Accessing fields
	fmt.Printf("First Name: %s\n", person1.FirstName)
	fmt.Printf("Full Name: %s %s\n", person1.FirstName, person1.LastName)

	// Modifying fields
	person1.Age = 31
	person1.Email = "john.newemail@example.com"
	fmt.Printf("Updated Person 1: %+v\n", person1)

	fmt.Println("\n=== Struct with Different Field Types ===")

	product := Product{
		ID:      101,
		Name:    "Laptop",
		Price:   999.99,
		InStock: true,
		Categories: []string{"Electronics", "Computers"},
		Attributes: map[string]string{
			"Brand":  "TechCorp",
			"Model":  "Pro 2024",
			"Color":  "Silver",
		},
	}
	fmt.Printf("Product: %+v\n", product)
	fmt.Printf("Categories: %v\n", product.Categories)
	fmt.Printf("Brand: %s\n", product.Attributes["Brand"])

	fmt.Println("\n=== Nested Structs ===")

	emp := Employee{
		ID:     1001,
		Name:   "Alice Johnson",
		Salary: 75000.00,
		Address: Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94102",
		},
	}
	fmt.Printf("Employee: %+v\n", emp)
	fmt.Printf("Employee works in: %s, %s\n", emp.Address.City, emp.Address.State)

	fmt.Println("\n=== Pointer to Struct ===")

	// Creating a pointer to a struct
	personPtr := &Person{
		FirstName: "Charlie",
		LastName:  "Brown",
		Age:       28,
		Email:     "charlie@example.com",
	}

	// Go automatically dereferences pointers to structs
	fmt.Printf("Person pointer: %+v\n", personPtr)
	fmt.Printf("First Name (via pointer): %s\n", personPtr.FirstName)

	// Modifying through pointer
	personPtr.Age = 29
	fmt.Printf("Updated age: %d\n", personPtr.Age)

	fmt.Println("\n=== Comparing Structs ===")

	p1 := Person{FirstName: "Test", LastName: "User", Age: 20}
	p2 := Person{FirstName: "Test", LastName: "User", Age: 20}
	p3 := Person{FirstName: "Test", LastName: "User", Age: 21}

	fmt.Printf("p1 == p2: %v\n", p1 == p2) // true (all fields are equal)
	fmt.Printf("p1 == p3: %v\n", p1 == p3) // false (Age differs)

	// Note: Structs with slices or maps cannot be compared with ==
	// product1 == product2 would cause a compile error
}
