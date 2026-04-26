package main_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/storage/local"
	"spec-streaming/backend/internal/videos"
	"spec-streaming/backend/internal/worker"
)

func TestFullTranscodingFlow(t *testing.T) {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/spec_streaming"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}
	defer pool.Close()

	// Clean up tables
	pool.Exec(ctx, "DELETE FROM transcoding_jobs")
	pool.Exec(ctx, "DELETE FROM videos")

	storage := local.New(t.TempDir())
	videoRepo := videos.NewPostgresRepository(pool)
	jobRepo := jobs.NewPostgresRepository(pool)

	jobService := jobs.NewService(jobRepo)
	videoService := videos.NewService(videoRepo, storage, jobService)
	videoHandler := videos.NewHandler(videoService, storage)

	// 1. Upload a video via API
	e := echo.New()
	e.POST("/videos", videoHandler.Create)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Test Video")
	part, _ := writer.CreateFormFile("file", "test.mp4")
	_, _ = part.Write([]byte("fake-mp4-content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/videos", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Verify job was created by querying directly
	var jobID, videoID string
	err = pool.QueryRow(ctx, "SELECT id, video_id FROM transcoding_jobs WHERE status = 'pending'").Scan(&jobID, &videoID)
	if err != nil {
		t.Fatalf("expected a pending job: %v", err)
	}

	// 3. Run worker to process the job
	transcoder := &dummyTranscoder{storage: storage}
	runner := &worker.Runner{
		Jobs:   jobService,
		Videos: videoService,
		Codec:  transcoder,
	}

	if err := runner.RunOnce(ctx); err != nil {
		t.Fatalf("runner error: %v", err)
	}

	// 4. Verify video is ready
	video, err := videoService.GetVideo(ctx, videoID)
	if err != nil {
		t.Fatalf("get video: %v", err)
	}
	if video.Status != videos.StatusReady {
		t.Fatalf("expected video status ready, got %s", video.Status)
	}
	if video.ManifestKey == "" {
		t.Fatal("expected manifest key to be set")
	}

	// 5. Verify manifest is served
	e2 := echo.New()
	e2.GET("/videos/:id/stream/manifest.mpd", videoHandler.Manifest)
	req2 := httptest.NewRequest(http.MethodGet, "/videos/"+video.ID+"/stream/manifest.mpd", nil)
	rec2 := httptest.NewRecorder()
	e2.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}
	if !bytes.Contains(rec2.Body.Bytes(), []byte("MPD")) {
		t.Fatalf("expected MPD content, got %s", rec2.Body.String())
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
      </Representation>
    </AdaptationSet>
  </Period>
</MPD>`
	if err := t.storage.SaveArtifact(manifestKey, bytes.NewBufferString(manifestContent)); err != nil {
		return "", err
	}
	return manifestKey, nil
}
