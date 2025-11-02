package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Contact struct holds the information for a single contact.
type Contact struct {
	Name        string
	PhoneNumber string
	Email       string
}

// ContactBook manages a collection of contacts.
// It uses a map for efficient lookups by name.
type ContactBook struct {
	contacts map[string]Contact
}

// NewContactBook creates and returns a new, initialized ContactBook.
func NewContactBook() *ContactBook {
	return &ContactBook{
		contacts: make(map[string]Contact),
	}
}

// Add adds or updates a contact in the book.
func (cb *ContactBook) Add(contact Contact) {
	fmt.Printf("Adding/Updating contact for '%s'\n", contact.Name)
	cb.contacts[contact.Name] = contact
}

// Get retrieves a contact by name.
// It returns the contact and a boolean indicating if the contact was found.
func (cb *ContactBook) Get(name string) (Contact, bool) {
	contact, found := cb.contacts[name]
	return contact, found
}

// Delete removes a contact by name.
func (cb *ContactBook) Delete(name string) bool {
	if _, found := cb.contacts[name]; found {
		delete(cb.contacts, name)
		return true
	}
	return false
}

// List displays all contacts in alphabetical order by name.
func (cb *ContactBook) List() {
	if len(cb.contacts) == 0 {
		fmt.Println("Contact book is empty.")
		return
	}

	// To print in sorted order, we extract the keys (names) into a slice, sort the slice,
	// and then iterate over the sorted slice to get the contact data from the map.
	var names []string
	for name := range cb.contacts {
		names = append(names, name)
	}
	sort.Strings(names) // Sort the names alphabetically

	fmt.Println("\n--- All Contacts ---")
	for _, name := range names {
		contact := cb.contacts[name]
		fmt.Printf("  Name: %-15s | Phone: %-15s | Email: %s\n", contact.Name, contact.PhoneNumber, contact.Email)
	}
	fmt.Println("--------------------")
}

// main function runs the interactive command-line interface for the contact book.
func main() {
	book := NewContactBook()
	// Pre-populate with some data
	book.Add(Contact{Name: "Alice", PhoneNumber: "123-456-7890", Email: "alice@example.com"})
	book.Add(Contact{Name: "Bob", PhoneNumber: "987-654-3210", Email: "bob@example.com"})

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Contact Management System!")

	for {
		fmt.Println("\nAvailable commands: add, get, delete, list, exit")
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input) // Split input by spaces
		if len(parts) == 0 {
			continue
		}
		command := parts[0]

		switch command {
		case "add":
			if len(parts) != 4 {
				fmt.Println("Usage: add <name> <phone> <email>")
				continue
			}
			book.Add(Contact{Name: parts[1], PhoneNumber: parts[2], Email: parts[3]})
		case "get":
			if len(parts) != 2 {
				fmt.Println("Usage: get <name>")
				continue
			}
			if contact, found := book.Get(parts[1]); found {
				fmt.Printf("Found contact -> Name: %s, Phone: %s, Email: %s\n", contact.Name, contact.PhoneNumber, contact.Email)
			} else {
				fmt.Printf("Contact '%s' not found.\n", parts[1])
			}
		case "delete":
			if len(parts) != 2 {
				fmt.Println("Usage: delete <name>")
				continue
			}
			if book.Delete(parts[1]) {
				fmt.Printf("Deleted contact '%s'.\n", parts[1])
			} else {
				fmt.Printf("Contact '%s' not found.\n", parts[1])
			}
		case "list":
			book.List()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
