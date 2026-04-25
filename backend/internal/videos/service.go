package videos

import (
	"context"
	"io"
	"math/rand"

	"spec-streaming/backend/internal/storage"
)

type Service struct {
	repo    Repository
	storage storage.Storage
}

func NewService(repo Repository, storage storage.Storage) *Service {
	return &Service{repo: repo, storage: storage}
}

func (s *Service) CreateVideo(ctx context.Context, title string, filename string, file io.Reader) (*Video, error) {
	video := &Video{
		ID:               generateID(),
		Title:            title,
		OriginalFilename: filename,
		Status:           StatusUploaded,
		SourceStorageKey: "sources/" + generateID(),
	}
	if err := s.storage.SaveSource(video.SourceStorageKey, file); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, video); err != nil {
		return nil, err
	}
	return video, nil
}

func (s *Service) ListVideos(ctx context.Context) ([]Video, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetVideo(ctx context.Context, id string) (*Video, error) {
	return s.repo.GetByID(ctx, id)
}

func generateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
