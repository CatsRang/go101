package main

import (
	"arena"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// ============================================================================
// Go Memory Arenas (Experimental Feature)
// ============================================================================
//
// To run this code, you must set the GOEXPERIMENT environment variable:
//   GOEXPERIMENT=arenas go run 99.etc/02_arenas.go
//
// To run with address sanitizer (detects use-after-free):
//   GOEXPERIMENT=arenas go run -asan 99.etc/02_arenas.go
//
// ============================================================================
// IMPORTANT: EXPERIMENTAL STATUS
// ============================================================================
// As of Go 1.20+, arenas are an experimental feature that may be changed
// incompatibly or removed at any time. The proposal is on hold indefinitely
// due to API safety concerns (use-after-free risks).
//
// For production code, consider alternatives like sync.Pool.
// ============================================================================

// Example structs for demonstration
type Person struct {
	ID   int
	Name string
	Age  int
}

type ComplexData struct {
	Items   []int
	Mapping map[string]int // Note: maps cannot be allocated directly in arenas
	Nested  *Person
}

type Node struct {
	Value int
	Next  *Node
}

func main() {
	fmt.Println("=" + repeatStr("=", 60))
	fmt.Println("Go Memory Arenas - Comprehensive Examples")
	fmt.Println("=" + repeatStr("=", 60))

	example1BasicUsage()
	example2SliceAllocation()
	example3Clone()
	example4ReflectArenaNew()
	example5LinkedListInArena()
	example6AppendBehavior()
	example7StringWorkaround()
	example8BatchProcessing()
	example9MemoryComparison()

	fmt.Println("\n" + repeatStr("=", 61))
	fmt.Println("All examples completed successfully!")
	fmt.Println(repeatStr("=", 61))
}

// ============================================================================
// Example 1: Basic Arena Usage
// ============================================================================
func example1BasicUsage() {
	printHeader("Example 1: Basic Arena Usage")

	// Create a new arena - all memory allocated from this arena will be freed together
	mem := arena.NewArena()
	defer mem.Free() // Ensure we free the memory when done

	fmt.Println("Arena created.")

	// Allocate a single object in the arena
	// arena.New[T] returns a pointer to a new value of type T
	p1 := arena.New[Person](mem)
	p1.ID = 1
	p1.Name = "Alice"
	p1.Age = 30

	fmt.Printf("Allocated Person: %+v at %p\n", *p1, p1)

	// Allocate another object
	p2 := arena.New[Person](mem)
	p2.ID = 2
	p2.Name = "Bob"
	p2.Age = 25

	fmt.Printf("Allocated Person: %+v at %p\n", *p2, p2)

	// Both p1 and p2 will be freed when mem.Free() is called
}

// ============================================================================
// Example 2: Slice Allocation in Arena
// ============================================================================
func example2SliceAllocation() {
	printHeader("Example 2: Slice Allocation")

	mem := arena.NewArena()
	defer mem.Free()

	// arena.MakeSlice[T] creates a slice backed by arena memory
	// Syntax: arena.MakeSlice[T](arena, length, capacity)
	people := arena.MakeSlice[Person](mem, 3, 5)

	// Note: The slice header is on the stack, but the underlying array is in the arena
	people[0] = Person{ID: 10, Name: "Charlie", Age: 35}
	people[1] = Person{ID: 11, Name: "Diana", Age: 28}
	people[2] = Person{ID: 12, Name: "Eve", Age: 42}

	fmt.Printf("Slice len=%d, cap=%d\n", len(people), cap(people))
	fmt.Println("Allocated Slice in Arena:")
	for i, p := range people {
		fmt.Printf("  [%d] %+v (addr: %p)\n", i, p, &people[i])
	}

	// Allocating large slices - arenas are most efficient for MiB-scale allocations
	largeSlice := arena.MakeSlice[int](mem, 1000, 1000)
	fmt.Printf("Large slice allocated: len=%d, cap=%d\n", len(largeSlice), cap(largeSlice))
}

// ============================================================================
// Example 3: arena.Clone - Escaping Arena Lifetime
// ============================================================================
func example3Clone() {
	printHeader("Example 3: arena.Clone")

	mem := arena.NewArena()

	// Allocate a Person in the arena
	arenaP := arena.New[Person](mem)
	arenaP.ID = 100
	arenaP.Name = "Frank"
	arenaP.Age = 50

	// Clone copies the value to the heap, escaping the arena lifetime
	// Clone works with pointers, slices, and strings
	heapP := arena.Clone(arenaP)

	fmt.Printf("Arena allocated: %+v at %p\n", *arenaP, arenaP)
	fmt.Printf("Heap cloned:     %+v at %p\n", *heapP, heapP)
	fmt.Printf("Same pointer? %v\n", arenaP == heapP) // false

	// Clone a slice
	arenaSlice := arena.MakeSlice[int](mem, 5, 5)
	for i := range arenaSlice {
		arenaSlice[i] = i * 10
	}
	heapSlice := arena.Clone(arenaSlice)

	fmt.Printf("Arena slice: %v\n", arenaSlice)
	fmt.Printf("Heap slice:  %v\n", heapSlice)

	// Free the arena - arenaP and arenaSlice become invalid
	mem.Free()

	// heapP and heapSlice are still valid!
	fmt.Printf("After Free - Heap pointer still valid: %+v\n", *heapP)
	fmt.Printf("After Free - Heap slice still valid: %v\n", heapSlice)

	// WARNING: Accessing arenaP or arenaSlice here would be undefined behavior!
}

// ============================================================================
// Example 4: Using Reflect with Arenas
// ============================================================================
func example4ReflectArenaNew() {
	printHeader("Example 4: Reflect ArenaNew")

	mem := arena.NewArena()
	defer mem.Free()

	// reflect.ArenaNew allows arena allocation when the type is known at runtime
	typ := reflect.TypeOf((*Person)(nil)).Elem()
	value := reflect.ArenaNew(mem, typ)

	// Set fields using reflection
	personPtr := value.Interface().(*Person)
	personPtr.ID = 200
	personPtr.Name = "Grace"
	personPtr.Age = 33

	fmt.Printf("Reflect-allocated Person: %+v\n", *personPtr)

	// Useful for generic serialization/deserialization into arenas
	intType := reflect.TypeOf((*int)(nil)).Elem()
	intVal := reflect.ArenaNew(mem, intType)
	intPtr := intVal.Interface().(*int)
	*intPtr = 42
	fmt.Printf("Reflect-allocated int: %d\n", *intPtr)
}

// ============================================================================
// Example 5: Linked List in Arena
// ============================================================================
func example5LinkedListInArena() {
	printHeader("Example 5: Linked List in Arena")

	mem := arena.NewArena()
	defer mem.Free()

	// Create a linked list entirely within the arena
	var head *Node
	var current *Node

	for i := 1; i <= 5; i++ {
		node := arena.New[Node](mem)
		node.Value = i * 100

		if head == nil {
			head = node
			current = head
		} else {
			current.Next = node
			current = node
		}
	}

	// Traverse and print
	fmt.Print("Linked list: ")
	for n := head; n != nil; n = n.Next {
		fmt.Printf("%d", n.Value)
		if n.Next != nil {
			fmt.Print(" -> ")
		}
	}
	fmt.Println()

	// All nodes freed together when arena is freed - no need to traverse and free
}

// ============================================================================
// Example 6: Append Behavior (Important Limitation!)
// ============================================================================
func example6AppendBehavior() {
	printHeader("Example 6: Append Behavior (IMPORTANT!)")

	mem := arena.NewArena()
	defer mem.Free()

	// Create a slice with specific capacity in arena
	slice := arena.MakeSlice[int](mem, 0, 3)
	fmt.Printf("Initial: len=%d, cap=%d\n", len(slice), cap(slice))

	// Appending within capacity stays in arena
	slice = append(slice, 1, 2, 3)
	fmt.Printf("After 3 appends (within cap): len=%d, cap=%d\n", len(slice), cap(slice))

	// WARNING: Appending beyond capacity moves the slice to the HEAP!
	sliceBeforeGrow := slice
	slice = append(slice, 4) // This exceeds capacity
	fmt.Printf("After 4th append (exceeds cap): len=%d, cap=%d\n", len(slice), cap(slice))

	// The slices now point to different backing arrays
	fmt.Printf("Same backing array? %v\n", &sliceBeforeGrow[0] == &slice[0])
	fmt.Println("WARNING: After capacity exceeded, new backing array is on the HEAP, not arena!")

	// Best practice: Pre-allocate with sufficient capacity
	fmt.Println("\nBest Practice: Pre-allocate with exact or maximum expected capacity")
}

// ============================================================================
// Example 7: String Workaround
// ============================================================================
func example7StringWorkaround() {
	printHeader("Example 7: String Workaround")

	mem := arena.NewArena()
	defer mem.Free()

	// Strings cannot be directly allocated in arenas
	// Workaround: Allocate []byte, copy data, convert to string

	sourceStr := "Hello, Arena World!"

	// Allocate byte slice in arena
	bytes := arena.MakeSlice[byte](mem, len(sourceStr), len(sourceStr))
	copy(bytes, sourceStr)

	// Convert to string using unsafe (the data is in arena)
	//
	// Why arenaStr is "in the arena":
	// --------------------------------
	// A Go string is internally just a header with two fields:
	//   type stringHeader struct {
	//       Data uintptr  // pointer to underlying byte array
	//       Len  int
	//   }
	//
	// Both 'bytes' and 'arenaStr' share the same underlying data in arena.
	// The string header lives on the stack, but the character data is in arena.
	// This is why arenaStr becomes INVALID after mem.Free()!
	//
	arenaStr := unsafe.String(&bytes[0], len(bytes))

	fmt.Printf("Original string: %s\n", sourceStr)
	fmt.Printf("Arena string:    %s\n", arenaStr)

	// WARNING: arenaStr becomes invalid when arena is freed!
	// To keep it, use arena.Clone
	heapStr := arena.Clone(arenaStr)
	fmt.Printf("Heap cloned string: %s (safe after arena.Free)\n", heapStr)
}

// ============================================================================
// Example 8: Batch Processing Pattern
// ============================================================================
func example8BatchProcessing() {
	printHeader("Example 8: Batch Processing Pattern")

	// Simulate processing batches of data
	// Each batch uses its own arena, freed after processing

	batches := [][]int{
		{1, 2, 3, 4, 5},
		{10, 20, 30},
		{100, 200, 300, 400},
	}

	var results []int

	for batchNum, batch := range batches {
		// Create arena for this batch
		mem := arena.NewArena()

		// Process batch in arena
		processed := arena.MakeSlice[int](mem, len(batch), len(batch))
		for i, v := range batch {
			processed[i] = v * 2 // Double each value
		}

		// Extract results we need (clone to escape arena)
		batchResult := arena.Clone(processed)
		results = append(results, batchResult...)

		// Free arena immediately after batch is done
		mem.Free()

		fmt.Printf("Batch %d processed and arena freed\n", batchNum+1)
	}

	fmt.Printf("All results: %v\n", results)
	fmt.Println("Pattern: Per-batch arenas minimize GC pressure during processing")
}

// ============================================================================
// Example 9: Memory Usage Comparison
// ============================================================================
func example9MemoryComparison() {
	printHeader("Example 9: Memory Usage Comparison")

	const numAllocs = 10000

	// Force GC and get baseline
	runtime.GC()
	var baseStats runtime.MemStats
	runtime.ReadMemStats(&baseStats)

	// --- Without Arena (Normal Heap Allocation) ---
	heapPointers := make([]*Person, numAllocs)
	for i := 0; i < numAllocs; i++ {
		heapPointers[i] = &Person{ID: i, Name: "Test", Age: 25}
	}

	runtime.GC()
	var heapStats runtime.MemStats
	runtime.ReadMemStats(&heapStats)

	heapPointers = nil // Allow collection
	runtime.GC()

	// --- With Arena ---
	mem := arena.NewArena()
	arenaPointers := make([]*Person, numAllocs) // slice header on heap
	for i := 0; i < numAllocs; i++ {
		arenaPointers[i] = arena.New[Person](mem)
		arenaPointers[i].ID = i
		arenaPointers[i].Name = "Test"
		arenaPointers[i].Age = 25
	}

	var arenaStats runtime.MemStats
	runtime.ReadMemStats(&arenaStats)

	mem.Free() // All arena memory freed at once

	var afterFreeStats runtime.MemStats
	runtime.ReadMemStats(&afterFreeStats)

	// Print comparison
	fmt.Println("Allocation Statistics (approximate):")
	fmt.Printf("  Heap allocations:  NumGC=%d\n", heapStats.NumGC-baseStats.NumGC)
	fmt.Printf("  Arena allocations: NumGC=%d\n", arenaStats.NumGC-heapStats.NumGC)
	fmt.Println("\nNote: Arena allocations don't trigger GC for arena-internal memory")
	fmt.Println("Benefit: Reduced GC overhead for bulk allocations with known lifetime")
}

// ============================================================================
// Helper Functions
// ============================================================================

func printHeader(title string) {
	fmt.Println()
	fmt.Println("-" + repeatStr("-", 60))
	fmt.Println(title)
	fmt.Println("-" + repeatStr("-", 60))
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// ============================================================================
// API Summary
// ============================================================================
//
// arena.NewArena() *Arena
//   - Creates a new arena for memory allocation
//
// (a *Arena) Free()
//   - Frees the arena and ALL objects allocated from it
//   - Must not use the arena or any allocated objects after calling Free()
//
// arena.New[T any](a *Arena) *T
//   - Allocates a new value of type T in the arena
//   - Returns a pointer to the allocated value
//
// arena.MakeSlice[T any](a *Arena, len, cap int) []T
//   - Creates a slice backed by arena memory
//   - Note: Slice header is on stack/heap, but underlying array is in arena
//
// arena.Clone[T any](s T) T
//   - Makes a shallow copy on the heap (escapes arena lifetime)
//   - T must be a pointer, slice, or string
//   - Returns value untouched if not allocated from an arena
//
// reflect.ArenaNew(a *Arena, typ reflect.Type) reflect.Value
//   - Arena allocation with runtime type information
//
// ============================================================================
// Limitations
// ============================================================================
//
// 1. Maps cannot be directly allocated in arenas
// 2. append() beyond capacity moves slice to heap
// 3. Strings require workaround (allocate []byte, use unsafe.String)
// 4. Not thread-safe - don't use same arena from multiple goroutines
// 5. Use-after-free is undefined behavior (may panic with -asan flag)
//
// ============================================================================
// When to Use Arenas
// ============================================================================
//
// Good use cases:
// - Request-scoped allocations in HTTP handlers
// - Batch processing with clear lifecycle
// - Protobuf/JSON parsing of large messages
// - Temporary tree/graph structures
// - Performance-critical code with many allocations (1k-10k items)
//
// Poor use cases:
// - Long-lived objects
// - Objects with unclear or mixed lifetimes
// - Small number of allocations (overhead not worth it)
// - Production code (experimental feature!)
//
// ============================================================================
// Alternatives for Production
// ============================================================================
//
// 1. sync.Pool - Reuses objects, reduces allocations
// 2. Object pooling - Custom pools for specific types
// 3. Slice reuse - Reset and reuse slices instead of reallocating
// 4. Pre-allocation - Allocate upfront with known sizes
//
// ============================================================================
