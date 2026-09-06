package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"
)

/*
> go test -bench=BenchmarkSubmit -run=^$ -benchmem -benchtime=10s

> go mod tidy
> go run ./10_rest_pattern.go
*/

/*
# Run for longer duration (more accurate results)
go test -bench=BenchmarkSubmit -run=^$ -benchmem -benchtime=10s

# Run with CPU profiling
go test -bench=BenchmarkSubmit -run=^$ -cpuprofile=cpu.prof

# Run with memory profiling
go test -bench=BenchmarkSubmit -run=^$ -memprofile=mem.prof

# View profile after running
go tool pprof cpu.prof
*/

func BenchmarkSubmit(b *testing.B) {
	client := &http.Client{Timeout: 5 * time.Second}
	b.ResetTimer()

	for b.Loop() {
		var wg sync.WaitGroup
		wg.Add(100) // number of concurrent clients to simulate
		for range 100 {
			go func() {
				defer wg.Done()
				job := Job{ID: "x", Data: "payload"}
				body, _ := json.Marshal(job)
				_, _ = client.Post("http://localhost:8080/submit", "application/json", bytes.NewReader(body))
			}()
		}
		wg.Wait()
	}
}
