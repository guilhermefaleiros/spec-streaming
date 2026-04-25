package worker

import (
	"context"

	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/videos"
)

type JobService interface {
	ClaimPending(ctx context.Context) (*jobs.Job, error)
	Update(ctx context.Context, job *jobs.Job) error
}

type VideoService interface {
	GetVideo(ctx context.Context, id string) (*videos.Video, error)
	Update(ctx context.Context, video *videos.Video) error
}

type Transcoder interface {
	Transcode(ctx context.Context, sourceKey string, videoID string) (string, error)
}

type Runner struct {
	Jobs   JobService
	Videos VideoService
	Codec  Transcoder
}

func (r *Runner) RunOnce(ctx context.Context) error {
	job, err := r.Jobs.ClaimPending(ctx)
	if err != nil || job == nil {
		return err
	}
	video, err := r.Videos.GetVideo(ctx, job.VideoID)
	if err != nil {
		return err
	}
	if err := video.MarkProcessing(); err != nil {
		return err
	}
	if err := r.Videos.Update(ctx, video); err != nil {
		return err
	}
	manifestKey, err := r.Codec.Transcode(ctx, video.SourceStorageKey, video.ID)
	if err != nil {
		video.MarkFailed(err.Error())
		_ = r.Videos.Update(ctx, video)
		job.Status = jobs.StatusFailed
		job.ErrorMessage = err.Error()
		return r.Jobs.Update(ctx, job)
	}
	if err := video.MarkReady(manifestKey); err != nil {
		return err
	}
	if err := r.Videos.Update(ctx, video); err != nil {
		return err
	}
	job.Status = jobs.StatusCompleted
	return r.Jobs.Update(ctx, job)
}
