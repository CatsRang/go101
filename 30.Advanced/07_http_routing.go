package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Go 1.22+ introduced enhanced routing patterns in http.ServeMux
// Features: Method-specific routing, path parameters, wildcards

// User represents a user entity
type UserEntity struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// In-memory user store
var users = map[int]*UserEntity{
	1: {ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
	2: {ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
}

var nextID = 3

func main() {
	mux := http.NewServeMux()

	// === Basic Route ===
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the API! Available endpoints:\n")
		fmt.Fprintf(w, "GET  /users      - List all users\n")
		fmt.Fprintf(w, "GET  /users/{id} - Get user by ID\n")
		fmt.Fprintf(w, "POST /users      - Create new user\n")
		fmt.Fprintf(w, "PUT  /users/{id} - Update user\n")
		fmt.Fprintf(w, "DELETE /users/{id} - Delete user\n")
	})

	// === Go 1.22+ Method-Specific Routing ===

	// GET /users - List all users
	mux.HandleFunc("GET /users", listUsers)

	// POST /users - Create a new user
	mux.HandleFunc("POST /users", createUser)

	// === Go 1.22+ Path Parameters with {param} ===

	// GET /users/{id} - Get user by ID
	mux.HandleFunc("GET /users/{id}", getUser)

	// PUT /users/{id} - Update user
	mux.HandleFunc("PUT /users/{id}", updateUser)

	// DELETE /users/{id} - Delete user
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	// === Go 1.22+ Wildcard Pattern with {path...} ===

	// GET /files/{path...} - Serve files with wildcard path
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path") // Gets everything after /files/
		fmt.Fprintf(w, "Requested file path: %s\n", path)
	})

	// === Health Check Endpoint ===
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	addr := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Println("\nTest with:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl http://localhost:8080/users")
	fmt.Println("  curl http://localhost:8080/users/1")
	fmt.Println("  curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"Charlie\",\"email\":\"charlie@example.com\"}' http://localhost:8080/users")
	fmt.Println("  curl -X DELETE http://localhost:8080/users/1")
	fmt.Println("  curl http://localhost:8080/files/dir/subdir/file.txt")

	log.Fatal(http.ListenAndServe(addr, mux))
}

// listUsers returns all users
func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userList := make([]*UserEntity, 0, len(users))
	for _, u := range users {
		userList = append(userList, u)
	}

	json.NewEncoder(w).Encode(userList)
}

// getUser returns a single user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Go 1.22+: Access path parameters with r.PathValue()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, exists := users[id]
	if !exists {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// createUser creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Email == "" {
		http.Error(w, `{"error":"name and email are required"}`, http.StatusBadRequest)
		return
	}

	user := &UserEntity{
		ID:        nextID,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}
	users[nextID] = user
	nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// updateUser updates an existing user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, exists := users[id]
	if !exists {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	json.NewEncoder(w).Encode(user)
}

// deleteUser deletes a user by ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	if _, exists := users[id]; !exists {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}
