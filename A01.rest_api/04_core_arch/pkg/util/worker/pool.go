package worker

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"example.com/rest-core-arch/pkg/util/log"
)

// Job defines the interface for work items
type Job interface {
	Execute(ctx context.Context) error
	Name() string
}

// Pool manages a pool of workers that process jobs
type Pool struct {
	workers   int
	jobQueue  chan Job
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a new worker pool that listens to the parent context
func NewPool(parentCtx context.Context, workers int, queueSize int) *Pool {
	// Create a child context to manage the pool's lifecycle
	ctx, cancel := context.WithCancel(parentCtx)

	return &Pool{
		workers:  workers,
		jobQueue: make(chan Job, queueSize),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start launches the worker goroutines
func (p *Pool) Start() {
	log.L().Info("Worker pool starting", zap.Int("workers", p.workers))
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				return // Channel closed, exit worker
			}
			
			// Execute job with the pool's context
			// In a real app, we might want a per-job timeout
			if err := job.Execute(p.ctx); err != nil {
				log.L().Error("Job failed", 
					zap.String("job", job.Name()), 
					zap.Int("worker", id), 
					zap.Error(err))
			} else {
				log.L().Debug("Job completed", 
					zap.String("job", job.Name()), 
					zap.Int("worker", id))
			}

		case <-p.ctx.Done():
			// Parent context cancelled, exit worker
			return 
		}
	}
}

// Submit adds a job to the queue. Returns error if queue is full or pool is shutting down.
func (p *Pool) Submit(job Job) error {
	select {
	case p.jobQueue <- job:
		return nil
	case <-p.ctx.Done():
		return errors.New("worker pool shutting down")
	default:
		return errors.New("job queue full")
	}
}

// Shutdown gracefully stops the worker pool
func (p *Pool) Shutdown(ctx context.Context) error {
	log.L().Info("Worker pool shutting down")
	
	// Signal workers to stop accepting new jobs via context
	p.cancel()

	// Close queue to ensure workers drain remaining items if logic allows
	// (Though p.ctx.Done() priority in select might skip them depending on runtime scheduling)
	// For strict draining, we wouldn't cancel p.ctx immediately, but we are following "Context-First"
	// which usually implies "Stop working now".
	// If we want to process pending items, we'd close the channel but NOT cancel context yet.
	// Let's stick to the guide: Context cancellation propagates shutdown.
	
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.L().Info("Worker pool stopped gracefully")
		return nil
	case <-ctx.Done():
		return ctx.Err() // Timeout
	}
}
