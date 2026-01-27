package example

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"example.com/rest-core-arch/pkg/util/log"
	"example.com/rest-core-arch/pkg/util/worker"
)

type Service struct {
	workerPool *worker.Pool
}

func NewService(pool *worker.Pool) *Service {
	return &Service{
		workerPool: pool,
	}
}

// Create simulates a business operation.
// It performs a synchronous check and then offloads work to the background worker.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	// 1. Synchronous Validation
	if req.Name == "error" {
		return nil, fmt.Errorf("invalid name")
	}

	id := uuid.New().String()

	// 2. Async Background Work
	job := &ExampleJob{
		ID:       id,
		JobName:  req.Name,
		Duration: time.Duration(req.Duration) * time.Millisecond,
	}

	if err := s.workerPool.Submit(job); err != nil {
		log.FromContext(ctx).Warn("Failed to submit job", zap.Error(err))
		return nil, fmt.Errorf("system overloaded, try again later")
	}

	return &CreateResponse{
		ID:      id,
		Message: "Request accepted",
		Status:  "processing",
	}, nil
}

// ExampleJob implements worker.Job interface
type ExampleJob struct {
	ID       string
	JobName  string
	Duration time.Duration
}

func (j *ExampleJob) Name() string {
	return fmt.Sprintf("job-%s", j.ID)
}

func (j *ExampleJob) Execute(ctx context.Context) error {
	log.L().Info("Processing job started", zap.String("id", j.ID), zap.String("name", j.JobName))
	
	select {
	case <-time.After(j.Duration):
		log.L().Info("Processing job finished", zap.String("id", j.ID))
		return nil
	case <-ctx.Done():
		log.L().Warn("Processing job cancelled", zap.String("id", j.ID))
		return ctx.Err()
	}
}
