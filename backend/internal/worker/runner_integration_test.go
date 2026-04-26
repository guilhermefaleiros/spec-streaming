package worker_test

import (
	"context"
	"testing"

	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/videos"
	"spec-streaming/backend/internal/worker"
)

type fakeJobService struct {
	repo map[string]*jobs.Job
}

func newFakeJobService() *fakeJobService {
	return &fakeJobService{repo: make(map[string]*jobs.Job)}
}

func (s *fakeJobService) ClaimPending(ctx context.Context) (*jobs.Job, error) {
	for _, job := range s.repo {
		if job.Status == jobs.StatusPending {
			job.Status = jobs.StatusProcessing
			return job, nil
		}
	}
	return nil, nil
}

func (s *fakeJobService) Update(ctx context.Context, job *jobs.Job) error {
	s.repo[job.ID] = job
	return nil
}

type fakeVideoService struct {
	repo map[string]*videos.Video
}

func newFakeVideoService() *fakeVideoService {
	return &fakeVideoService{repo: make(map[string]*videos.Video)}
}

func (s *fakeVideoService) GetVideo(ctx context.Context, id string) (*videos.Video, error) {
	return s.repo[id], nil
}

func (s *fakeVideoService) Update(ctx context.Context, video *videos.Video) error {
	s.repo[video.ID] = video
	return nil
}

type fakeTranscoder struct {
	manifestKey string
}

func (t *fakeTranscoder) Transcode(ctx context.Context, sourceKey string, videoID string) (string, error) {
	return t.manifestKey, nil
}

func TestRunnerProcessesPendingJob(t *testing.T) {
	jobSvc := newFakeJobService()
	videoSvc := newFakeVideoService()
	transcoder := &fakeTranscoder{manifestKey: "videos/1/manifest.mpd"}

	video := &videos.Video{
		ID:               "vid-1",
		Status:           videos.StatusUploaded,
		SourceStorageKey: "sources/vid-1",
	}
	videoSvc.repo[video.ID] = video

	job := &jobs.Job{
		ID:      "job-1",
		VideoID: video.ID,
		Status:  jobs.StatusPending,
	}
	jobSvc.repo[job.ID] = job

	runner := &worker.Runner{
		Jobs:   jobSvc,
		Videos: videoSvc,
		Codec:  transcoder,
	}

	if err := runner.RunOnce(context.Background()); err != nil {
		t.Fatalf("run once: %v", err)
	}

	if video.Status != videos.StatusReady {
		t.Fatalf("expected video status ready, got %s", video.Status)
	}
	if video.ManifestKey != "videos/1/manifest.mpd" {
		t.Fatalf("expected manifest key, got %s", video.ManifestKey)
	}
	if job.Status != jobs.StatusCompleted {
		t.Fatalf("expected job status completed, got %s", job.Status)
	}
}
