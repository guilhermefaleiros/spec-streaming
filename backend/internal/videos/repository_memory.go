package videos

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	videos map[string]*Video
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{videos: make(map[string]*Video)}
}

func (r *MemoryRepository) Create(ctx context.Context, v *Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.videos[v.ID] = v
	return nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []Video
	for _, v := range r.videos {
		list = append(list, *v)
	}
	return list, nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.videos[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (r *MemoryRepository) Update(ctx context.Context, v *Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.videos[v.ID] = v
	return nil
}
