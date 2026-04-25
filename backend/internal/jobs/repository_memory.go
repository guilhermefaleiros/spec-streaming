package jobs

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{jobs: make(map[string]*Job)}
}

func (r *MemoryRepository) Create(ctx context.Context, j *Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}

func (r *MemoryRepository) ClaimPending(ctx context.Context) (*Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, job := range r.jobs {
		if job.Status == StatusPending {
			job.Status = StatusProcessing
			job.Attempts++
			return job, nil
		}
	}
	return nil, nil
}

func (r *MemoryRepository) Update(ctx context.Context, j *Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}
