package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// Using SQLite for easy demonstration (no external database needed)
	// Install: go get github.com/mattn/go-sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// User model
type DBUser struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt sql.NullTime // Nullable field
}

func main() {
	// === Open Database Connection ===
	// For SQLite, the file will be created if it doesn't exist
	db, err := sql.Open("sqlite3", ":memory:") // Use in-memory DB for demo
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// === Connection Pool Settings ===
	db.SetMaxOpenConns(25)                 // Maximum open connections
	db.SetMaxIdleConns(5)                  // Maximum idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection max lifetime

	// === Verify Connection ===
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Connected to database successfully!")

	// === Create Table ===
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME
	);`

	if _, err := db.ExecContext(ctx, createTableSQL); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table created successfully!")

	// === INSERT - Create Users ===
	fmt.Println("\n=== Creating Users ===")

	// Using Exec for simple inserts
	result, err := db.ExecContext(ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Alice", "alice@example.com")
	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}

	lastID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Inserted user with ID: %d, rows affected: %d\n", lastID, rowsAffected)

	// Insert more users
	users := []struct {
		name, email string
	}{
		{"Bob", "bob@example.com"},
		{"Charlie", "charlie@example.com"},
	}

	for _, u := range users {
		_, err := db.ExecContext(ctx,
			"INSERT INTO users (name, email) VALUES (?, ?)",
			u.name, u.email)
		if err != nil {
			log.Printf("Failed to insert %s: %v", u.name, err)
		}
	}

	// === SELECT - Query Single Row ===
	fmt.Println("\n=== Query Single User ===")

	var user DBUser
	err = db.QueryRowContext(ctx,
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?",
		1).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		fmt.Println("User not found")
	} else if err != nil {
		log.Fatal("Failed to query user:", err)
	} else {
		fmt.Printf("Found user: %+v\n", user)
	}

	// === SELECT - Query Multiple Rows ===
	fmt.Println("\n=== Query All Users ===")

	rows, err := db.QueryContext(ctx,
		"SELECT id, name, email, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows.Close()

	var allUsers []DBUser
	for rows.Next() {
		var u DBUser
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		allUsers = append(allUsers, u)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatal("Error during rows iteration:", err)
	}

	for _, u := range allUsers {
		fmt.Printf("  User %d: %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// === UPDATE ===
	fmt.Println("\n=== Update User ===")

	result, err = db.ExecContext(ctx,
		"UPDATE users SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		"Alice Updated", 1)
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}

	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Updated %d row(s)\n", rowsAffected)

	// === Prepared Statements (Better for repeated queries) ===
	fmt.Println("\n=== Prepared Statements ===")

	stmt, err := db.PrepareContext(ctx,
		"SELECT name, email FROM users WHERE id = ?")
	if err != nil {
		log.Fatal("Failed to prepare statement:", err)
	}
	defer stmt.Close()

	for id := 1; id <= 3; id++ {
		var name, email string
		err := stmt.QueryRowContext(ctx, id).Scan(&name, &email)
		if err != nil {
			log.Printf("Failed to query user %d: %v", id, err)
			continue
		}
		fmt.Printf("  ID %d: %s (%s)\n", id, name, email)
	}

	// === Transactions ===
	fmt.Println("\n=== Transactions ===")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	// Perform operations within transaction
	_, err = tx.ExecContext(ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Dave", "dave@example.com")
	if err != nil {
		tx.Rollback()
		log.Fatal("Transaction failed:", err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Eve", "eve@example.com")
	if err != nil {
		tx.Rollback()
		log.Fatal("Transaction failed:", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
	fmt.Println("Transaction committed successfully!")

	// === DELETE ===
	fmt.Println("\n=== Delete User ===")

	result, err = db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", 5)
	if err != nil {
		log.Fatal("Failed to delete user:", err)
	}

	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Deleted %d row(s)\n", rowsAffected)

	// === Final Count ===
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal("Failed to count users:", err)
	}
	fmt.Printf("\nTotal users in database: %d\n", count)

	fmt.Println("\nDatabase examples completed!")
}

// Note: To run this example, install the SQLite driver:
// go get github.com/mattn/go-sqlite3
//
// For PostgreSQL, use:
// import _ "github.com/lib/pq"
// db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
//
// For MySQL, use:
// import _ "github.com/go-sql-driver/mysql"
// db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/dbname")
