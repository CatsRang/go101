package main

import (
	"fmt"
	"time"
)

// Category represents a product category
type Category struct {
	ID          int
	Name        string
	Description string
}

// NewCategory creates a new category
func NewCategory(id int, name, description string) *Category {
	return &Category{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

// DisplayInfo shows category information
func (c Category) DisplayInfo() {
	fmt.Printf("Category [%d]: %s - %s\n", c.ID, c.Name, c.Description)
}

// Supplier represents a product supplier
type Supplier struct {
	ID      int
	Name    string
	Contact string
	Email   string
	Phone   string
	Address string
}

// NewSupplier creates a new supplier
func NewSupplier(id int, name, contact, email, phone, address string) *Supplier {
	return &Supplier{
		ID:      id,
		Name:    name,
		Contact: contact,
		Email:   email,
		Phone:   phone,
		Address: address,
	}
}

// DisplayInfo shows supplier information
func (s Supplier) DisplayInfo() {
	fmt.Printf("Supplier [%d]: %s\n", s.ID, s.Name)
	fmt.Printf("  Contact: %s\n", s.Contact)
	fmt.Printf("  Email: %s, Phone: %s\n", s.Email, s.Phone)
	fmt.Printf("  Address: %s\n", s.Address)
}

// Product represents an inventory item
type Product struct {
	ID          int
	SKU         string
	Name        string
	Description string
	Price       float64
	Cost        float64
	Quantity    int
	MinStock    int
	Category    *Category
	Supplier    *Supplier
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new product with validation
func NewProduct(id int, sku, name, description string, price, cost float64,
	quantity, minStock int, category *Category, supplier *Supplier) (*Product, error) {

	if sku == "" {
		return nil, fmt.Errorf("SKU cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("product name cannot be empty")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if cost < 0 {
		return nil, fmt.Errorf("cost cannot be negative")
	}
	if quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	return &Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		Cost:        cost,
		Quantity:    quantity,
		MinStock:    minStock,
		Category:    category,
		Supplier:    supplier,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// AddStock increases product quantity
func (p *Product) AddStock(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	p.Quantity += quantity
	p.UpdatedAt = time.Now()
	fmt.Printf("Added %d units to %s. New quantity: %d\n", quantity, p.Name, p.Quantity)
	return nil
}

// RemoveStock decreases product quantity
func (p *Product) RemoveStock(quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if quantity > p.Quantity {
		return fmt.Errorf("insufficient stock: have %d, need %d", p.Quantity, quantity)
	}
	p.Quantity -= quantity
	p.UpdatedAt = time.Now()
	fmt.Printf("Removed %d units from %s. Remaining: %d\n", quantity, p.Name, p.Quantity)
	return nil
}

// UpdatePrice changes product price
func (p *Product) UpdatePrice(newPrice float64) error {
	if newPrice < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	oldPrice := p.Price
	p.Price = newPrice
	p.UpdatedAt = time.Now()
	fmt.Printf("Updated price for %s: $%.2f -> $%.2f\n", p.Name, oldPrice, newPrice)
	return nil
}

// IsLowStock checks if product is below minimum stock level
func (p Product) IsLowStock() bool {
	return p.Quantity <= p.MinStock
}

// GetProfit calculates profit per unit
func (p Product) GetProfit() float64 {
	return p.Price - p.Cost
}

// GetProfitMargin calculates profit margin percentage
func (p Product) GetProfitMargin() float64 {
	if p.Price == 0 {
		return 0
	}
	return (p.GetProfit() / p.Price) * 100
}

// GetTotalValue calculates total inventory value
func (p Product) GetTotalValue() float64 {
	return float64(p.Quantity) * p.Cost
}

// DisplayInfo shows detailed product information
func (p Product) DisplayInfo() {
	fmt.Printf("\n=== Product Information ===\n")
	fmt.Printf("ID: %d, SKU: %s\n", p.ID, p.SKU)
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Price: $%.2f, Cost: $%.2f\n", p.Price, p.Cost)
	fmt.Printf("Profit per unit: $%.2f (%.1f%% margin)\n", p.GetProfit(), p.GetProfitMargin())
	fmt.Printf("Quantity: %d (Min: %d)\n", p.Quantity, p.MinStock)
	if p.IsLowStock() {
		fmt.Printf("⚠️  LOW STOCK WARNING!\n")
	}
	fmt.Printf("Total Value: $%.2f\n", p.GetTotalValue())
	if p.Category != nil {
		fmt.Printf("Category: %s\n", p.Category.Name)
	}
	if p.Supplier != nil {
		fmt.Printf("Supplier: %s\n", p.Supplier.Name)
	}
	fmt.Printf("Created: %s\n", p.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", p.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// InventoryManager manages the inventory system
type InventoryManager struct {
	products   map[int]*Product
	categories map[int]*Category
	suppliers  map[int]*Supplier
	nextID     int
}

// NewInventoryManager creates a new inventory manager
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		products:   make(map[int]*Product),
		categories: make(map[int]*Category),
		suppliers:  make(map[int]*Supplier),
		nextID:     1,
	}
}

// AddCategory adds a category to the system
func (im *InventoryManager) AddCategory(category *Category) {
	im.categories[category.ID] = category
	fmt.Printf("Added category: %s\n", category.Name)
}

// AddSupplier adds a supplier to the system
func (im *InventoryManager) AddSupplier(supplier *Supplier) {
	im.suppliers[supplier.ID] = supplier
	fmt.Printf("Added supplier: %s\n", supplier.Name)
}

// AddProduct adds a product to inventory
func (im *InventoryManager) AddProduct(product *Product) {
	im.products[product.ID] = product
	fmt.Printf("Added product: %s (SKU: %s)\n", product.Name, product.SKU)
}

// GetProduct retrieves a product by ID
func (im *InventoryManager) GetProduct(id int) (*Product, error) {
	product, exists := im.products[id]
	if !exists {
		return nil, fmt.Errorf("product with ID %d not found", id)
	}
	return product, nil
}

// GetProductsByCategoryID returns all products in a category
func (im *InventoryManager) GetProductsByCategoryID(categoryID int) []*Product {
	var products []*Product
	for _, p := range im.products {
		if p.Category != nil && p.Category.ID == categoryID {
			products = append(products, p)
		}
	}
	return products
}

// GetLowStockProducts returns products below minimum stock
func (im *InventoryManager) GetLowStockProducts() []*Product {
	var lowStock []*Product
	for _, p := range im.products {
		if p.IsLowStock() {
			lowStock = append(lowStock, p)
		}
	}
	return lowStock
}

// GetTotalInventoryValue calculates total value of all inventory
func (im *InventoryManager) GetTotalInventoryValue() float64 {
	total := 0.0
	for _, p := range im.products {
		total += p.GetTotalValue()
	}
	return total
}

// DisplayInventorySummary shows inventory statistics
func (im *InventoryManager) DisplayInventorySummary() {
	fmt.Printf("\n=== Inventory Summary ===\n")
	fmt.Printf("Total Products: %d\n", len(im.products))
	fmt.Printf("Total Categories: %d\n", len(im.categories))
	fmt.Printf("Total Suppliers: %d\n", len(im.suppliers))
	fmt.Printf("Total Inventory Value: $%.2f\n", im.GetTotalInventoryValue())

	lowStock := im.GetLowStockProducts()
	if len(lowStock) > 0 {
		fmt.Printf("\n⚠️  Low Stock Alert: %d products\n", len(lowStock))
		for _, p := range lowStock {
			fmt.Printf("  - %s: %d units (min: %d)\n", p.Name, p.Quantity, p.MinStock)
		}
	}
}

func main() {
	fmt.Println("=== Inventory Management System ===\n")

	// Initialize inventory manager
	manager := NewInventoryManager()

	// Create categories
	electronics := NewCategory(1, "Electronics", "Electronic devices and accessories")
	furniture := NewCategory(2, "Furniture", "Office and home furniture")
	supplies := NewCategory(3, "Office Supplies", "General office supplies")

	manager.AddCategory(electronics)
	manager.AddCategory(furniture)
	manager.AddCategory(supplies)

	fmt.Println()

	// Create suppliers
	techCorp := NewSupplier(1, "TechCorp Inc.", "John Smith",
		"john@techcorp.com", "555-0101", "123 Tech Street, Silicon Valley, CA")
	furnitureCo := NewSupplier(2, "Furniture Co.", "Jane Doe",
		"jane@furnitureco.com", "555-0102", "456 Furniture Ave, Portland, OR")
	officeMax := NewSupplier(3, "OfficeMax", "Bob Johnson",
		"bob@officemax.com", "555-0103", "789 Supply Blvd, New York, NY")

	manager.AddSupplier(techCorp)
	manager.AddSupplier(furnitureCo)
	manager.AddSupplier(officeMax)

	fmt.Println()

	// Create products
	laptop, _ := NewProduct(1, "LAPTOP-001", "Dell XPS 15", "High-performance laptop",
		1299.99, 900.00, 15, 5, electronics, techCorp)

	mouse, _ := NewProduct(2, "MOUSE-001", "Logitech MX Master", "Wireless mouse",
		99.99, 60.00, 3, 10, electronics, techCorp)

	desk, _ := NewProduct(3, "DESK-001", "Standing Desk", "Adjustable height desk",
		499.99, 300.00, 8, 3, furniture, furnitureCo)

	chair, _ := NewProduct(4, "CHAIR-001", "Ergonomic Chair", "Comfortable office chair",
		299.99, 180.00, 12, 5, furniture, furnitureCo)

	paper, _ := NewProduct(5, "PAPER-001", "Printer Paper", "500 sheets white paper",
		9.99, 5.00, 2, 20, supplies, officeMax)

	manager.AddProduct(laptop)
	manager.AddProduct(mouse)
	manager.AddProduct(desk)
	manager.AddProduct(chair)
	manager.AddProduct(paper)

	fmt.Println()

	// Display category information
	fmt.Println("=== Categories ===")
	electronics.DisplayInfo()
	furniture.DisplayInfo()
	supplies.DisplayInfo()

	fmt.Println()

	// Display supplier information
	fmt.Println("=== Suppliers ===")
	techCorp.DisplayInfo()
	fmt.Println()

	// Display product details
	laptop.DisplayInfo()
	mouse.DisplayInfo()

	fmt.Println("\n=== Stock Operations ===")

	// Add stock
	laptop.AddStock(10)
	mouse.AddStock(15)

	// Remove stock
	laptop.RemoveStock(3)
	desk.RemoveStock(2)

	// Update price
	mouse.UpdatePrice(89.99)

	// Display inventory summary
	manager.DisplayInventorySummary()

	fmt.Println("\n=== Products by Category ===")
	electronicsProducts := manager.GetProductsByCategoryID(electronics.ID)
	fmt.Printf("\nElectronics Category (%d products):\n", len(electronicsProducts))
	for _, p := range electronicsProducts {
		fmt.Printf("  - %s: %d units @ $%.2f\n", p.Name, p.Quantity, p.Price)
	}

	furnitureProducts := manager.GetProductsByCategoryID(furniture.ID)
	fmt.Printf("\nFurniture Category (%d products):\n", len(furnitureProducts))
	for _, p := range furnitureProducts {
		fmt.Printf("  - %s: %d units @ $%.2f\n", p.Name, p.Quantity, p.Price)
	}

	fmt.Println("\n=== Error Handling Examples ===")

	// Try to remove too much stock
	err := laptop.RemoveStock(1000)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Try to set negative price
	err = laptop.UpdatePrice(-100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Try to create invalid product
	_, err = NewProduct(999, "", "Invalid Product", "No SKU", 10.00, 5.00, 1, 0, nil, nil)
	if err != nil {
		fmt.Printf("Error creating product: %v\n", err)
	}

	// Final summary
	manager.DisplayInventorySummary()

	fmt.Println("\n=== System Concepts Demonstrated ===")
	fmt.Println("✓ Struct definition and composition")
	fmt.Println("✓ Constructor patterns with validation")
	fmt.Println("✓ Methods with value and pointer receivers")
	fmt.Println("✓ Embedded structs (Category, Supplier in Product)")
	fmt.Println("✓ Error handling")
	fmt.Println("✓ Encapsulation and data management")
	fmt.Println("✓ Business logic implementation")
}
