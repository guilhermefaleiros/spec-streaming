package jobs

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateJob(ctx context.Context, videoID string) (*Job, error) {
	job := &Job{
		ID:      videoID + "-job",
		VideoID: videoID,
		Status:  StatusPending,
	}
	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *Service) ClaimPending(ctx context.Context) (*Job, error) {
	return s.repo.ClaimPending(ctx)
}

func (s *Service) Update(ctx context.Context, job *Job) error {
	return s.repo.Update(ctx, job)
}
