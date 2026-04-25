package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/storage/local"
	"spec-streaming/backend/internal/videos"
	"spec-streaming/backend/internal/worker"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/spec_streaming"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	storage := local.New("tmp/storage")
	videoRepo := videos.NewPostgresRepository(pool)
	jobRepo := jobs.NewPostgresRepository(pool)

	jobService := jobs.NewService(jobRepo)
	videoService := videos.NewService(videoRepo, storage, jobService)

	transcoder := &dummyTranscoder{storage: storage}

	runner := &worker.Runner{
		Jobs:   jobService,
		Videos: videoService,
		Codec:  transcoder,
	}

	log.Println("Worker starting...")
	for {
		if err := runner.RunOnce(context.Background()); err != nil {
			log.Printf("runner error: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}

type dummyTranscoder struct {
	storage *local.Storage
}

func (t *dummyTranscoder) Transcode(ctx context.Context, sourceKey string, videoID string) (string, error) {
	manifestKey := "artifacts/" + videoID + "/manifest.mpd"
	manifestContent := `<?xml version="1.0"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static">
  <Period>
    <AdaptationSet mimeType="video/mp4">
      <Representation id="1" bandwidth="1000000">
        <BaseURL>segment.m4s</BaseURL>
        <SegmentBase indexRange="0-100"/>
      </Representation>
    </AdaptationSet>
  </Period>
</MPD>`
	if err := t.storage.SaveArtifact(manifestKey, bytes.NewBufferString(manifestContent)); err != nil {
		return "", err
	}
	return manifestKey, nil
}
