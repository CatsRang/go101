package main

import (
	"testing"
)

// Run tests: go test -v ./30.Advanced/
// Run with coverage: go test -cover ./30.Advanced/
// Run benchmarks: go test -bench=. ./30.Advanced/

// === Basic Unit Test ===
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

// === Table-Driven Tests ===
func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		want      float64
		wantError bool
	}{
		{"positive numbers", 10, 2, 5, false},
		{"negative dividend", -10, 2, -5, false},
		{"decimal result", 7, 2, 3.5, false},
		{"divide by zero", 10, 0, 0, true},
		{"zero dividend", 0, 5, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)

			if tt.wantError {
				if err == nil {
					t.Errorf("Divide(%v, %v) expected error, got nil", tt.a, tt.b)
				}
				return
			}

			if err != nil {
				t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
				return
			}

			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// === Subtests ===
func TestIsPalindrome(t *testing.T) {
	t.Run("valid palindromes", func(t *testing.T) {
		palindromes := []string{"radar", "level", "A", "aa", "Radar"}
		for _, p := range palindromes {
			if !IsPalindrome(p) {
				t.Errorf("IsPalindrome(%q) = false; want true", p)
			}
		}
	})

	t.Run("non-palindromes", func(t *testing.T) {
		nonPalindromes := []string{"hello", "world", "ab"}
		for _, np := range nonPalindromes {
			if IsPalindrome(np) {
				t.Errorf("IsPalindrome(%q) = true; want false", np)
			}
		}
	})

	t.Run("empty string", func(t *testing.T) {
		if !IsPalindrome("") {
			t.Error("IsPalindrome(\"\") = false; want true")
		}
	})
}

// === Testing Structs and Methods ===
func TestUserValidate(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		wantError string
	}{
		{
			name:      "valid user",
			user:      User{ID: 1, Name: "John", Email: "john@example.com"},
			wantError: "",
		},
		{
			name:      "missing name",
			user:      User{ID: 1, Name: "", Email: "john@example.com"},
			wantError: "name is required",
		},
		{
			name:      "missing email",
			user:      User{ID: 1, Name: "John", Email: ""},
			wantError: "email is required",
		},
		{
			name:      "invalid email format",
			user:      User{ID: 1, Name: "John", Email: "invalid-email"},
			wantError: "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantError == "" {
				if err != nil {
					t.Errorf("User.Validate() unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("User.Validate() expected error %q, got nil", tt.wantError)
				return
			}

			if err.Error() != tt.wantError {
				t.Errorf("User.Validate() error = %q; want %q", err.Error(), tt.wantError)
			}
		})
	}
}

// === Benchmark Tests ===
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}

// === Benchmark with Different Input Sizes ===
func BenchmarkFibonacciIterativeSizes(b *testing.B) {
	sizes := []int{10, 20, 30, 40}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10))+string(rune('0'+size%10)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FibonacciIterative(size)
			}
		})
	}
}

// === Test Helper Functions ===
func TestStringSliceContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	tests := []struct {
		target string
		want   bool
	}{
		{"apple", true},
		{"banana", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			got := StringSliceContains(slice, tt.target)
			if got != tt.want {
				t.Errorf("StringSliceContains(..., %q) = %v; want %v", tt.target, got, tt.want)
			}
		})
	}
}

// === Go 1.24+ T.Context() Example ===
// func TestWithContext(t *testing.T) {
//     ctx := t.Context() // Automatically cancelled when test ends
//
//     // Use ctx for operations that need cancellation
//     select {
//     case <-ctx.Done():
//         t.Log("Context was cancelled")
//     default:
//         t.Log("Context is still active")
//     }
// }

// === Fuzz Testing (Go 1.18+) ===
func FuzzIsPalindrome(f *testing.F) {
	// Seed corpus
	f.Add("radar")
	f.Add("hello")
	f.Add("")
	f.Add("a")

	f.Fuzz(func(t *testing.T, s string) {
		// Just check it doesn't panic
		_ = IsPalindrome(s)
	})
}

// === Test with Parallel Execution ===
func TestParallel(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"test1", 1, 2, 3},
		{"test2", 5, 5, 10},
		{"test3", -1, 1, 0},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable (not needed in Go 1.22+)
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run subtests in parallel

			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
