package main

import (
	"fmt"
	"math"
)

// ============================================================================
// INTERFACES: The Go Way of Polymorphism
// ============================================================================

// Shape interface defines behavior for geometric shapes
// Interfaces specify WHAT, not HOW
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle implements Shape interface (implicitly)
type Rectangle struct {
	Width  float64
	Height float64
}

// Area calculates rectangle area
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates rectangle perimeter
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle implements Shape interface (implicitly)
type Circle struct {
	Radius float64
}

// Area calculates circle area
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates circle circumference
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Triangle implements Shape interface
type Triangle struct {
	A, B, C float64 // Side lengths
}

// Area calculates triangle area using Heron's formula
func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2 // Semi-perimeter
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter calculates triangle perimeter
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// PrintShapeInfo accepts any type that implements Shape
func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// CalculateTotalArea calculates total area of multiple shapes
func CalculateTotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// ============================================================================
// MULTIPLE INTERFACES
// ============================================================================

// Animal interface defines animal behavior
type Animal interface {
	Speak() string
	Move() string
}

// Pet interface defines pet-specific behavior
type Pet interface {
	Name() string
	Owner() string
}

// Dog implements both Animal and Pet
type Dog struct {
	name  string
	owner string
}

func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) Move() string {
	return "Running on four legs"
}

func (d Dog) Name() string {
	return d.name
}

func (d Dog) Owner() string {
	return d.owner
}

// Cat implements both Animal and Pet
type Cat struct {
	name  string
	owner string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Move() string {
	return "Walking gracefully"
}

func (c Cat) Name() string {
	return c.name
}

func (c Cat) Owner() string {
	return c.owner
}

// Bird implements only Animal
type Bird struct {
	species string
}

func (b Bird) Speak() string {
	return "Chirp!"
}

func (b Bird) Move() string {
	return "Flying through the air"
}

// MakeAnimalSpeak accepts anything that implements Animal
func MakeAnimalSpeak(a Animal) {
	fmt.Printf("The animal says: %s\n", a.Speak())
	fmt.Printf("It moves by: %s\n", a.Move())
}

// IntroducePet accepts anything that implements Pet
func IntroducePet(p Pet) {
	fmt.Printf("This is %s, owned by %s\n", p.Name(), p.Owner())
}

// ============================================================================
// IMPLICIT IMPLEMENTATION
// ============================================================================

// Writer interface (simplified version of io.Writer)
type Writer interface {
	Write(data string) error
}

// ConsoleWriter writes to console
type ConsoleWriter struct{}

// Write implements Writer interface
func (cw ConsoleWriter) Write(data string) error {
	fmt.Println(data)
	return nil
}

// FileWriter simulates file writing
type FileWriter struct {
	filename string
}

// Write implements Writer interface
func (fw FileWriter) Write(data string) error {
	fmt.Printf("[Writing to %s]: %s\n", fw.filename, data)
	return nil
}

// WriteMessage uses any Writer implementation
func WriteMessage(w Writer, message string) {
	w.Write(message)
}

// ============================================================================
// INTERFACE VALUES
// ============================================================================

// Greeter interface
type Greeter interface {
	Greet() string
}

// EnglishGreeter greets in English
type EnglishGreeter struct {
	name string
}

func (eg EnglishGreeter) Greet() string {
	return "Hello, " + eg.name
}

// SpanishGreeter greets in Spanish
type SpanishGreeter struct {
	name string
}

func (sg SpanishGreeter) Greet() string {
	return "Hola, " + sg.name
}

// FrenchGreeter greets in French
type FrenchGreeter struct {
	name string
}

func (fg FrenchGreeter) Greet() string {
	return "Bonjour, " + fg.name
}

// ============================================================================
// NIL INTERFACES
// ============================================================================

// Validator interface for validation
type Validator interface {
	Validate() bool
}

// Email type with validation
type Email struct {
	address string
}

func (e Email) Validate() bool {
	return len(e.address) > 0 && contains(e.address, "@")
}

// Simple contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== 1. Basic Interface Implementation ===")

	// Create different shapes
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	triangle := Triangle{A: 3, B: 4, C: 5}

	// All implement Shape interface
	fmt.Println("Rectangle:")
	PrintShapeInfo(rect)

	fmt.Println("\nCircle:")
	PrintShapeInfo(circle)

	fmt.Println("\nTriangle:")
	PrintShapeInfo(triangle)

	fmt.Println("\n=== 2. Collection of Interfaces ===")

	// Store different shapes in a slice of Shape interface
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 7},
		Triangle{A: 3, B: 4, C: 5},
		Rectangle{Width: 8, Height: 3},
	}

	totalArea := CalculateTotalArea(shapes)
	fmt.Printf("Total area of all shapes: %.2f\n", totalArea)

	fmt.Println("\n=== 3. Multiple Interfaces ===")

	dog := Dog{name: "Buddy", owner: "Alice"}
	cat := Cat{name: "Whiskers", owner: "Bob"}
	bird := Bird{species: "Sparrow"}

	// Dog and Cat implement both Animal and Pet
	fmt.Println("Dog as Animal:")
	MakeAnimalSpeak(dog)

	fmt.Println("\nCat as Animal:")
	MakeAnimalSpeak(cat)

	fmt.Println("\nBird as Animal:")
	MakeAnimalSpeak(bird)

	fmt.Println("\nDog as Pet:")
	IntroducePet(dog)

	fmt.Println("\nCat as Pet:")
	IntroducePet(cat)

	// Bird doesn't implement Pet, so this would be a compile error:
	// IntroducePet(bird) // ✗ Cannot do this!

	fmt.Println("\n=== 4. Implicit Implementation ===")

	console := ConsoleWriter{}
	fileWriter := FileWriter{filename: "output.txt"}

	fmt.Println("Using ConsoleWriter:")
	WriteMessage(console, "Hello from console!")

	fmt.Println("\nUsing FileWriter:")
	WriteMessage(fileWriter, "Hello from file!")

	// No explicit "implements" keyword needed!
	// Types automatically satisfy interfaces by implementing methods

	fmt.Println("\n=== 5. Interface Values ===")

	// Store different greeters
	greeters := []Greeter{
		EnglishGreeter{name: "John"},
		SpanishGreeter{name: "Carlos"},
		FrenchGreeter{name: "Pierre"},
		EnglishGreeter{name: "Jane"},
	}

	fmt.Println("Greetings:")
	for i, greeter := range greeters {
		fmt.Printf("%d: %s\n", i+1, greeter.Greet())
	}

	fmt.Println("\n=== 6. Nil Interface Check ===")

	var validator Validator

	// validator is nil (no concrete type assigned)
	if validator == nil {
		fmt.Println("Validator is nil")
	}

	// Assign a concrete type
	validator = Email{address: "user@example.com"}
	if validator != nil {
		fmt.Printf("Validator is not nil: %v\n", validator.Validate())
	}

	// Empty email (invalid)
	validator = Email{address: ""}
	fmt.Printf("Empty email validation: %v\n", validator.Validate())

	// Valid email
	validator = Email{address: "test@domain.com"}
	fmt.Printf("Valid email validation: %v\n", validator.Validate())

	fmt.Println("\n=== 7. Polymorphism in Action ===")

	// Create a slice of Animals (polymorphic collection)
	animals := []Animal{
		Dog{name: "Max", owner: "Charlie"},
		Cat{name: "Luna", owner: "Diana"},
		Bird{species: "Parrot"},
		Dog{name: "Rocky", owner: "Eve"},
	}

	fmt.Println("All animals:")
	for i, animal := range animals {
		fmt.Printf("\n%d. %s - %s\n", i+1, animal.Speak(), animal.Move())
	}

	fmt.Println("\n=== Key Concepts ===")
	fmt.Println("✓ Interfaces define behavior (methods), not data")
	fmt.Println("✓ Types implicitly implement interfaces (no 'implements' keyword)")
	fmt.Println("✓ A type can implement multiple interfaces")
	fmt.Println("✓ Interface values can hold any concrete type that implements it")
	fmt.Println("✓ Interfaces enable polymorphism in Go")
	fmt.Println("✓ Keep interfaces small and focused (often 1-3 methods)")
	fmt.Println("✓ Accept interfaces, return concrete types")
}
