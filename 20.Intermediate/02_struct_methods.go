package main

import (
	"fmt"
	"math"
)

// Rectangle demonstrates methods with different receiver types
type Rectangle struct {
	Width  float64
	Height float64
}

// Area calculates the area using a value receiver
// Value receiver gets a copy of the struct
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter using a value receiver
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Scale modifies the rectangle using a pointer receiver
// Pointer receiver can modify the original struct
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// SetDimensions sets new dimensions using a pointer receiver
func (r *Rectangle) SetDimensions(width, height float64) {
	r.Width = width
	r.Height = height
}

// IsSquare checks if rectangle is a square (value receiver is fine here)
func (r Rectangle) IsSquare() bool {
	return r.Width == r.Height
}

// Circle demonstrates methods for a different type
type Circle struct {
	Radius float64
}

// Area calculates circle area
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circumference calculates circle circumference
func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// Grow increases the circle radius by a pointer receiver
func (c *Circle) Grow(increment float64) {
	c.Radius += increment
}

// BankAccount demonstrates practical use of pointer receivers
type BankAccount struct {
	AccountNumber string
	Owner         string
	Balance       float64
}

// Deposit adds money to the account (pointer receiver needed)
func (b *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		b.Balance += amount
		fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, b.Balance)
	}
}

// Withdraw removes money from the account (pointer receiver needed)
func (b *BankAccount) Withdraw(amount float64) bool {
	if amount > 0 && amount <= b.Balance {
		b.Balance -= amount
		fmt.Printf("Withdrew $%.2f. New balance: $%.2f\n", amount, b.Balance)
		return true
	}
	fmt.Printf("Withdrawal failed. Insufficient funds.\n")
	return false
}

// GetBalance returns the current balance (value receiver is fine)
func (b BankAccount) GetBalance() float64 {
	return b.Balance
}

// DisplayInfo shows account information (value receiver is fine)
func (b BankAccount) DisplayInfo() {
	fmt.Printf("Account: %s, Owner: %s, Balance: $%.2f\n",
		b.AccountNumber, b.Owner, b.Balance)
}

// Counter demonstrates why pointer receiver matters
type Counter struct {
	Value int
}

// IncrementValue tries to increment using value receiver (won't work)
func (c Counter) IncrementValue() {
	c.Value++ // This modifies a copy, not the original
}

// IncrementPointer increments using pointer receiver (works correctly)
func (c *Counter) IncrementPointer() {
	c.Value++ // This modifies the original
}

// Reset sets counter to zero
func (c *Counter) Reset() {
	c.Value = 0
}

// GetValue returns the current value
func (c Counter) GetValue() int {
	return c.Value
}

func main() {
	fmt.Println("=== Value Receiver vs Pointer Receiver ===")

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Initial rectangle: %+v\n", rect)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
	fmt.Printf("Is square: %v\n", rect.IsSquare())

	fmt.Println("\n=== Modifying with Pointer Receiver ===")

	// Scale modifies the original struct
	rect.Scale(2.0)
	fmt.Printf("After scaling by 2: %+v\n", rect)
	fmt.Printf("New area: %.2f\n", rect.Area())

	// Go automatically takes the address when calling pointer receiver methods
	rect.SetDimensions(15, 15)
	fmt.Printf("After setting dimensions: %+v\n", rect)
	fmt.Printf("Is square now: %v\n", rect.IsSquare())

	fmt.Println("\n=== Circle Methods ===")

	circle := Circle{Radius: 5.0}
	fmt.Printf("Circle: %+v\n", circle)
	fmt.Printf("Area: %.2f\n", circle.Area())
	fmt.Printf("Circumference: %.2f\n", circle.Circumference())

	circle.Grow(3.0)
	fmt.Printf("After growing by 3: %+v\n", circle)
	fmt.Printf("New area: %.2f\n", circle.Area())

	fmt.Println("\n=== Bank Account Example ===")

	account := BankAccount{
		AccountNumber: "ACC-001",
		Owner:         "John Doe",
		Balance:       1000.00,
	}

	account.DisplayInfo()
	account.Deposit(500.00)
	account.Withdraw(200.00)
	account.Withdraw(2000.00) // Will fail
	fmt.Printf("Final balance: $%.2f\n", account.GetBalance())

	fmt.Println("\n=== Pointer Receiver Importance ===")

	counter := Counter{Value: 0}
	fmt.Printf("Initial counter: %+v\n", counter)

	// This won't work - value receiver modifies a copy
	counter.IncrementValue()
	fmt.Printf("After IncrementValue: %+v (no change!)\n", counter)

	// This works - pointer receiver modifies the original
	counter.IncrementPointer()
	fmt.Printf("After IncrementPointer: %+v (changed!)\n", counter)

	counter.IncrementPointer()
	counter.IncrementPointer()
	fmt.Printf("After 2 more increments: %+v\n", counter)

	counter.Reset()
	fmt.Printf("After reset: %+v\n", counter)

	fmt.Println("\n=== Method Calls on Pointers ===")

	// You can call methods on pointers too
	rectPtr := &Rectangle{Width: 8, Height: 4}
	fmt.Printf("Rectangle pointer: %+v\n", rectPtr)

	// Go automatically dereferences for value receiver methods
	fmt.Printf("Area (via pointer): %.2f\n", rectPtr.Area())

	// Pointer receiver methods work naturally
	rectPtr.Scale(0.5)
	fmt.Printf("After scaling: %+v\n", rectPtr)

	fmt.Println("\n=== Guidelines for Choosing Receiver Type ===")
	fmt.Println("Use POINTER receiver when:")
	fmt.Println("  - Method needs to modify the receiver")
	fmt.Println("  - Struct is large (avoid copying)")
	fmt.Println("  - Consistency: if one method uses pointer, others should too")
	fmt.Println("\nUse VALUE receiver when:")
	fmt.Println("  - Method doesn't modify the receiver")
	fmt.Println("  - Struct is small and copyable")
	fmt.Println("  - Working with immutable types")
}
