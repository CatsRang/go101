package main

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// PRACTICAL EXERCISE: PAYMENT PROCESSING SYSTEM
// ============================================================================
// This system demonstrates:
// - Interface-based design
// - Multiple payment methods (credit card, PayPal, cryptocurrency)
// - Polymorphism
// - Composition
// - Error handling with custom error types

// ============================================================================
// 1. CORE INTERFACES
// ============================================================================

// PaymentMethod defines the interface for all payment methods
type PaymentMethod interface {
	// ProcessPayment processes a payment and returns transaction ID and error
	ProcessPayment(amount float64) (string, error)

	// GetName returns the payment method name
	GetName() string

	// Validate validates the payment method
	Validate() error
}

// Refundable interface for payment methods that support refunds
type Refundable interface {
	ProcessRefund(transactionID string, amount float64) error
}

// RecurringPayment interface for subscription-based payments
type RecurringPayment interface {
	SetupRecurring(amount float64, interval time.Duration) error
	CancelRecurring() error
}

// ============================================================================
// 2. CUSTOM ERROR TYPES
// ============================================================================

// PaymentError represents a payment processing error
type PaymentError struct {
	Method  string
	Amount  float64
	Reason  string
	Code    int
}

func (pe PaymentError) Error() string {
	return fmt.Sprintf("[%s] Payment of $%.2f failed (code %d): %s",
		pe.Method, pe.Amount, pe.Code, pe.Reason)
}

// InsufficientFundsError represents insufficient funds
type InsufficientFundsError struct {
	Available float64
	Required  float64
}

func (ife InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: have $%.2f, need $%.2f",
		ife.Available, ife.Required)
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", ve.Field, ve.Message)
}

// ============================================================================
// 3. CREDIT CARD PAYMENT
// ============================================================================

// CreditCard represents a credit card payment method
type CreditCard struct {
	CardNumber     string
	CardHolder     string
	ExpiryMonth    int
	ExpiryYear     int
	CVV            string
	BillingAddress string
}

// NewCreditCard creates a new credit card payment method
func NewCreditCard(number, holder, cvv string, month, year int) *CreditCard {
	return &CreditCard{
		CardNumber:  number,
		CardHolder:  holder,
		ExpiryMonth: month,
		ExpiryYear:  year,
		CVV:         cvv,
	}
}

// ProcessPayment implements PaymentMethod
func (cc *CreditCard) ProcessPayment(amount float64) (string, error) {
	if err := cc.Validate(); err != nil {
		return "", err
	}

	if amount <= 0 {
		return "", PaymentError{
			Method: cc.GetName(),
			Amount: amount,
			Reason: "amount must be positive",
			Code:   400,
		}
	}

	// Simulate payment processing
	transactionID := fmt.Sprintf("CC-%s-%d", cc.maskCardNumber(), time.Now().Unix())
	fmt.Printf("Processing credit card payment: $%.2f\n", amount)
	fmt.Printf("Card: %s\n", cc.maskCardNumber())
	fmt.Printf("Transaction ID: %s\n", transactionID)

	return transactionID, nil
}

// GetName implements PaymentMethod
func (cc *CreditCard) GetName() string {
	return "Credit Card"
}

// Validate implements PaymentMethod
func (cc *CreditCard) Validate() error {
	if len(cc.CardNumber) < 13 || len(cc.CardNumber) > 19 {
		return ValidationError{Field: "CardNumber", Message: "invalid length"}
	}
	if cc.CardHolder == "" {
		return ValidationError{Field: "CardHolder", Message: "required"}
	}
	if cc.ExpiryMonth < 1 || cc.ExpiryMonth > 12 {
		return ValidationError{Field: "ExpiryMonth", Message: "must be 1-12"}
	}
	if len(cc.CVV) != 3 && len(cc.CVV) != 4 {
		return ValidationError{Field: "CVV", Message: "must be 3 or 4 digits"}
	}

	// Check expiry
	now := time.Now()
	if cc.ExpiryYear < now.Year() || (cc.ExpiryYear == now.Year() && cc.ExpiryMonth < int(now.Month())) {
		return ValidationError{Field: "Expiry", Message: "card has expired"}
	}

	return nil
}

// ProcessRefund implements Refundable
func (cc *CreditCard) ProcessRefund(transactionID string, amount float64) error {
	fmt.Printf("Processing refund: $%.2f to card %s (transaction: %s)\n",
		amount, cc.maskCardNumber(), transactionID)
	return nil
}

// maskCardNumber masks the card number for security
func (cc *CreditCard) maskCardNumber() string {
	if len(cc.CardNumber) < 4 {
		return "****"
	}
	return "****-****-****-" + cc.CardNumber[len(cc.CardNumber)-4:]
}

// ============================================================================
// 4. PAYPAL PAYMENT
// ============================================================================

// PayPal represents PayPal payment method
type PayPal struct {
	Email    string
	Password string
	Balance  float64
}

// NewPayPal creates a new PayPal payment method
func NewPayPal(email, password string, balance float64) *PayPal {
	return &PayPal{
		Email:    email,
		Password: password,
		Balance:  balance,
	}
}

// ProcessPayment implements PaymentMethod
func (pp *PayPal) ProcessPayment(amount float64) (string, error) {
	if err := pp.Validate(); err != nil {
		return "", err
	}

	if amount > pp.Balance {
		return "", InsufficientFundsError{
			Available: pp.Balance,
			Required:  amount,
		}
	}

	// Process payment
	pp.Balance -= amount
	transactionID := fmt.Sprintf("PP-%s-%d", strings.Split(pp.Email, "@")[0], time.Now().Unix())

	fmt.Printf("Processing PayPal payment: $%.2f\n", amount)
	fmt.Printf("Account: %s\n", pp.Email)
	fmt.Printf("Remaining balance: $%.2f\n", pp.Balance)
	fmt.Printf("Transaction ID: %s\n", transactionID)

	return transactionID, nil
}

// GetName implements PaymentMethod
func (pp *PayPal) GetName() string {
	return "PayPal"
}

// Validate implements PaymentMethod
func (pp *PayPal) Validate() error {
	if !strings.Contains(pp.Email, "@") {
		return ValidationError{Field: "Email", Message: "invalid email format"}
	}
	if pp.Password == "" {
		return ValidationError{Field: "Password", Message: "required"}
	}
	return nil
}

// ProcessRefund implements Refundable
func (pp *PayPal) ProcessRefund(transactionID string, amount float64) error {
	pp.Balance += amount
	fmt.Printf("Processing PayPal refund: $%.2f to %s (transaction: %s)\n",
		amount, pp.Email, transactionID)
	fmt.Printf("New balance: $%.2f\n", pp.Balance)
	return nil
}

// ============================================================================
// 5. CRYPTOCURRENCY PAYMENT
// ============================================================================

// Cryptocurrency represents crypto payment method
type Cryptocurrency struct {
	WalletAddress string
	CryptoType    string // BTC, ETH, etc.
	Balance       float64
}

// NewCryptocurrency creates a new cryptocurrency payment method
func NewCryptocurrency(address, cryptoType string, balance float64) *Cryptocurrency {
	return &Cryptocurrency{
		WalletAddress: address,
		CryptoType:    cryptoType,
		Balance:       balance,
	}
}

// ProcessPayment implements PaymentMethod
func (cr *Cryptocurrency) ProcessPayment(amount float64) (string, error) {
	if err := cr.Validate(); err != nil {
		return "", err
	}

	if amount > cr.Balance {
		return "", InsufficientFundsError{
			Available: cr.Balance,
			Required:  amount,
		}
	}

	// Process payment
	cr.Balance -= amount
	transactionID := fmt.Sprintf("%s-%s-%d",
		cr.CryptoType, cr.WalletAddress[:8], time.Now().Unix())

	fmt.Printf("Processing %s payment: $%.2f\n", cr.CryptoType, amount)
	fmt.Printf("Wallet: %s...%s\n", cr.WalletAddress[:6], cr.WalletAddress[len(cr.WalletAddress)-4:])
	fmt.Printf("Remaining balance: $%.2f\n", cr.Balance)
	fmt.Printf("Transaction ID: %s\n", transactionID)

	return transactionID, nil
}

// GetName implements PaymentMethod
func (cr *Cryptocurrency) GetName() string {
	return fmt.Sprintf("Cryptocurrency (%s)", cr.CryptoType)
}

// Validate implements PaymentMethod
func (cr *Cryptocurrency) Validate() error {
	if len(cr.WalletAddress) < 26 {
		return ValidationError{Field: "WalletAddress", Message: "invalid address"}
	}
	if cr.CryptoType == "" {
		return ValidationError{Field: "CryptoType", Message: "required"}
	}
	return nil
}

// ============================================================================
// 6. PAYMENT PROCESSOR
// ============================================================================

// PaymentProcessor handles payment processing
type PaymentProcessor struct {
	transactions []Transaction
}

// Transaction represents a payment transaction
type Transaction struct {
	ID        string
	Method    string
	Amount    float64
	Status    string
	Timestamp time.Time
}

// NewPaymentProcessor creates a new payment processor
func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{
		transactions: make([]Transaction, 0),
	}
}

// Process processes a payment using any PaymentMethod
func (pp *PaymentProcessor) Process(method PaymentMethod, amount float64) error {
	fmt.Printf("\n=== Processing Payment with %s ===\n", method.GetName())

	transactionID, err := method.ProcessPayment(amount)
	if err != nil {
		pp.recordTransaction(Transaction{
			ID:        "FAILED",
			Method:    method.GetName(),
			Amount:    amount,
			Status:    "FAILED",
			Timestamp: time.Now(),
		})
		return err
	}

	pp.recordTransaction(Transaction{
		ID:        transactionID,
		Method:    method.GetName(),
		Amount:    amount,
		Status:    "SUCCESS",
		Timestamp: time.Now(),
	})

	fmt.Printf("✓ Payment successful!\n")
	return nil
}

// Refund processes a refund if the payment method supports it
func (pp *PaymentProcessor) Refund(method PaymentMethod, transactionID string, amount float64) error {
	fmt.Printf("\n=== Processing Refund ===\n")

	// Check if method supports refunds using type assertion
	if refundable, ok := method.(Refundable); ok {
		return refundable.ProcessRefund(transactionID, amount)
	}

	return fmt.Errorf("payment method %s does not support refunds", method.GetName())
}

// recordTransaction records a transaction
func (pp *PaymentProcessor) recordTransaction(tx Transaction) {
	pp.transactions = append(pp.transactions, tx)
}

// GetTransactionHistory returns all transactions
func (pp *PaymentProcessor) GetTransactionHistory() []Transaction {
	return pp.transactions
}

// PrintHistory prints transaction history
func (pp *PaymentProcessor) PrintHistory() {
	fmt.Println("\n=== Transaction History ===")
	for i, tx := range pp.transactions {
		fmt.Printf("%d. [%s] %s - %s: $%.2f (%s)\n",
			i+1, tx.Status, tx.Method, tx.ID, tx.Amount,
			tx.Timestamp.Format("15:04:05"))
	}
}

// ============================================================================
// MAIN: DEMONSTRATION
// ============================================================================

func main() {
	fmt.Println("=== PAYMENT PROCESSING SYSTEM ===\n")

	// Create payment processor
	processor := NewPaymentProcessor()

	// Create different payment methods
	creditCard := NewCreditCard("4532123456789012", "John Doe", "123", 12, 2025)
	paypal := NewPayPal("alice@example.com", "secure123", 500.00)
	bitcoin := NewCryptocurrency("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "BTC", 1000.00)

	// Store payment methods in a slice (polymorphism)
	paymentMethods := []PaymentMethod{creditCard, paypal, bitcoin}

	fmt.Println("Available payment methods:")
	for i, method := range paymentMethods {
		fmt.Printf("%d. %s\n", i+1, method.GetName())
	}

	// Process payments with different methods
	fmt.Println("\n=== Processing Payments ===")

	// Credit card payment
	err := processor.Process(creditCard, 99.99)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// PayPal payment
	err = processor.Process(paypal, 150.00)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Bitcoin payment
	err = processor.Process(bitcoin, 250.00)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Try payment with insufficient funds
	err = processor.Process(paypal, 1000.00)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Process refunds (demonstrating Refundable interface)
	fmt.Println("\n=== Processing Refunds ===")

	// Refund to credit card
	txID := "CC-****-****-****-9012-" + fmt.Sprint(time.Now().Unix())
	processor.Refund(creditCard, txID, 50.00)

	// Refund to PayPal
	processor.Refund(paypal, "PP-alice-123", 30.00)

	// Try to refund to Bitcoin (doesn't implement Refundable)
	err = processor.Refund(bitcoin, "BTC-1A1zP1eP-123", 25.00)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Show transaction history
	processor.PrintHistory()

	// Demonstrate validation
	fmt.Println("\n=== Validation Examples ===")

	invalidCard := NewCreditCard("123", "Jane Doe", "12", 13, 2020)
	err = processor.Process(invalidCard, 50.00)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

	invalidPaypal := NewPayPal("invalid-email", "pass", 100.00)
	err = processor.Process(invalidPaypal, 25.00)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

	// Type assertions to check capabilities
	fmt.Println("\n=== Checking Payment Method Capabilities ===")

	for _, method := range paymentMethods {
		fmt.Printf("\n%s:\n", method.GetName())

		// Check if refundable
		if _, ok := method.(Refundable); ok {
			fmt.Println("  ✓ Supports refunds")
		} else {
			fmt.Println("  ✗ Does not support refunds")
		}

		// Check if supports recurring payments
		if _, ok := method.(RecurringPayment); ok {
			fmt.Println("  ✓ Supports recurring payments")
		} else {
			fmt.Println("  ✗ Does not support recurring payments")
		}
	}

	fmt.Println("\n=== System Concepts Demonstrated ===")
	fmt.Println("✓ Interface-based polymorphism")
	fmt.Println("✓ Multiple implementations of PaymentMethod interface")
	fmt.Println("✓ Interface composition (Refundable, RecurringPayment)")
	fmt.Println("✓ Type assertions to check capabilities")
	fmt.Println("✓ Custom error types for specific error handling")
	fmt.Println("✓ Validation pattern")
	fmt.Println("✓ Accept interfaces (PaymentMethod), return concrete types")
	fmt.Println("✓ Real-world application of Go interfaces")
}
