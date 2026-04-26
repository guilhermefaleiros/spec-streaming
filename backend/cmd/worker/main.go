package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	transcoder := &ffmpegTranscoder{storage: storage}

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

type ffmpegTranscoder struct {
	storage *local.Storage
}

func (t *ffmpegTranscoder) Transcode(ctx context.Context, sourceKey string, videoID string) (string, error) {
	inputPath := t.storage.SourcePath(sourceKey)
	outDir := filepath.Join(os.TempDir(), "spec-streaming-transcode", videoID)
	manifestKey := "artifacts/" + videoID + "/manifest.mpd"

	// Ensure output directory exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}
	defer os.RemoveAll(outDir) // Clean up temp dir when done

	outputPath := filepath.Join(outDir, "manifest.mpd")

	// Build ffmpeg command for DASH
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "dash",
		"-init_seg_name", "init_$RepresentationID$.m4s",
		"-media_seg_name", "chunk_$RepresentationID$_$Number%05d$.m4s",
		"-use_template", "1",
		"-use_timeline", "1",
		outputPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("[ffmpeg] starting transcoding for video %s", videoID)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}
	log.Printf("[ffmpeg] transcoding complete for video %s", videoID)

	// Copy generated files from temp dir to storage as artifacts
	entries, err := os.ReadDir(outDir)
	if err != nil {
		return "", fmt.Errorf("read output dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		srcPath := filepath.Join(outDir, entry.Name())
		artifactKey := "artifacts/" + videoID + "/" + entry.Name()

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", fmt.Errorf("read artifact %s: %w", entry.Name(), err)
		}

		if err := t.storage.SaveArtifact(artifactKey, strings.NewReader(string(data))); err != nil {
			return "", fmt.Errorf("save artifact %s: %w", entry.Name(), err)
		}
		log.Printf("[ffmpeg] saved artifact: %s", artifactKey)
	}

	return manifestKey, nil
}
