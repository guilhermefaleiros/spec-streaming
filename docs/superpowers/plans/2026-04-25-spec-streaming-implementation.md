# Spec Streaming Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a streaming prototype that accepts MP4 uploads, processes them asynchronously into MPEG-DASH, lists uploaded videos, and plays processed videos through a React frontend.

**Architecture:** Implement a modular monolith with a Go API process and a separate Go worker process sharing PostgreSQL and a storage abstraction. The frontend is a React/Vite app that uploads files, polls processing status, lists videos, and opens a player page that consumes DASH assets served by the API.

**Tech Stack:** Go, Echo, PostgreSQL, ffmpeg, React, Vite, TypeScript, shadcn/ui, Playwright

---

## File Structure

### Back-end

- Create: `backend/go.mod`
- Create: `backend/cmd/api/main.go`
- Create: `backend/cmd/worker/main.go`
- Create: `backend/internal/app/config.go`
- Create: `backend/internal/platform/postgres/postgres.go`
- Create: `backend/internal/platform/http/router.go`
- Create: `backend/internal/platform/http/errors.go`
- Create: `backend/internal/storage/storage.go`
- Create: `backend/internal/storage/local/local.go`
- Create: `backend/internal/storage/s3/s3.go`
- Create: `backend/internal/videos/video.go`
- Create: `backend/internal/videos/repository.go`
- Create: `backend/internal/videos/service.go`
- Create: `backend/internal/videos/handler.go`
- Create: `backend/internal/jobs/job.go`
- Create: `backend/internal/jobs/repository.go`
- Create: `backend/internal/jobs/service.go`
- Create: `backend/internal/transcoding/transcoder.go`
- Create: `backend/internal/transcoding/ffmpeg.go`
- Create: `backend/internal/worker/runner.go`
- Create: `backend/migrations/0001_create_videos_and_jobs.sql`

### Back-end Tests

- Create: `backend/internal/videos/service_test.go`
- Create: `backend/internal/storage/local/local_test.go`
- Create: `backend/internal/videos/handler_integration_test.go`
- Create: `backend/internal/worker/runner_integration_test.go`

### Front-end

- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/App.tsx`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/types.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/components/upload-form.tsx`
- Create: `frontend/src/components/video-list.tsx`
- Create: `frontend/src/components/status-badge.tsx`
- Create: `frontend/src/components/video-player.tsx`
- Create: `frontend/src/pages/home-page.tsx`
- Create: `frontend/src/pages/video-page.tsx`

### Front-end Tests

- Create: `frontend/src/components/upload-form.test.tsx`
- Create: `frontend/src/components/video-list.test.tsx`
- Create: `frontend/playwright.config.ts`
- Create: `frontend/tests/e2e/upload-and-play.spec.ts`
- Create: `frontend/tests/fixtures/sample.mp4`

### Root Tooling

- Create: `docker-compose.yml`
- Create: `Makefile`
- Create: `.gitignore`

### Suggested Dependency Boundaries

- `videos` owns the `Video` entity, validation, status rules, and HTTP use cases.
- `jobs` owns job lifecycle, locking, and worker-facing transitions.
- `storage` hides local and S3-compatible persistence.
- `transcoding` wraps ffmpeg execution and artifact layout.
- `worker` orchestrates pending-job polling and processing.

### Task 1: Bootstrap the repository skeleton

**Files:**
- Create: `.gitignore`
- Create: `Makefile`
- Create: `docker-compose.yml`
- Create: `backend/go.mod`
- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`

- [ ] **Step 1: Write the failing smoke checks for expected project structure**

Create `backend/internal/videos/service_test.go` with an initial package compile target:

```go
package videos_test

import "testing"

func TestBackendPackageBootstraps(t *testing.T) {}
```

Create `frontend/src/components/upload-form.test.tsx` with a basic frontend compile target:

```tsx
import { describe, it, expect } from 'vitest'

describe('frontend bootstrap', () => {
  it('loads the test runner', () => {
    expect(true).toBe(true)
  })
})
```

- [ ] **Step 2: Run the smoke checks to verify the repo is not ready yet**

Run: `make test-backend`
Expected: FAIL with `make: *** No rule to make target 'test-backend'`

Run: `cd frontend && npm test`
Expected: FAIL because `package.json` does not exist yet

- [ ] **Step 3: Add minimal repository tooling and package manifests**

Create `.gitignore`:

```gitignore
.DS_Store
backend/.cache/
backend/bin/
frontend/node_modules/
frontend/dist/
frontend/playwright-report/
frontend/test-results/
tmp/
```

Create `Makefile`:

```make
.PHONY: test-backend test-frontend test-e2e

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

test-e2e:
	cd frontend && npx playwright test
```

Create `backend/go.mod`:

```go
module spec-streaming/backend

go 1.24

require (
	github.com/labstack/echo/v4 v4.13.3
	github.com/jackc/pgx/v5 v5.7.2
)
```

Create `frontend/package.json`:

```json
{
  "name": "spec-streaming-frontend",
  "private": true,
  "version": "0.0.1",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "test": "vitest",
    "test:e2e": "playwright test"
  },
  "dependencies": {
    "react": "^19.1.0",
    "react-dom": "^19.1.0",
    "react-router-dom": "^7.6.0"
  },
  "devDependencies": {
    "@playwright/test": "^1.51.0",
    "@testing-library/react": "^16.3.0",
    "@types/react": "^19.1.2",
    "@types/react-dom": "^19.1.2",
    "@vitejs/plugin-react": "^4.4.1",
    "typescript": "^5.8.3",
    "vite": "^6.3.5",
    "vitest": "^3.1.4"
  }
}
```

- [ ] **Step 4: Run the smoke checks again**

Run: `make test-backend`
Expected: PASS or `ok` for the placeholder backend test

Run: `make test-frontend`
Expected: FAIL until the frontend test config exists

- [ ] **Step 5: Commit the bootstrap**

```bash
git add .gitignore Makefile docker-compose.yml backend/go.mod frontend/package.json frontend/vite.config.ts frontend/tsconfig.json backend/internal/videos/service_test.go frontend/src/components/upload-form.test.tsx
git commit -m "chore: bootstrap project structure"
```

### Task 2: Implement the database schema and core domain rules

**Files:**
- Create: `backend/migrations/0001_create_videos_and_jobs.sql`
- Create: `backend/internal/videos/video.go`
- Create: `backend/internal/jobs/job.go`
- Create: `backend/internal/videos/service.go`
- Create: `backend/internal/jobs/service.go`
- Test: `backend/internal/videos/service_test.go`

- [ ] **Step 1: Write failing tests for valid status transitions**

Update `backend/internal/videos/service_test.go`:

```go
package videos_test

import (
	"testing"

	"spec-streaming/backend/internal/videos"
)

func TestVideoStatusTransitions(t *testing.T) {
	video := videos.Video{Status: videos.StatusUploaded}

	if err := video.MarkProcessing(); err != nil {
		t.Fatalf("expected uploaded -> processing to pass: %v", err)
	}

	if err := video.MarkReady("videos/123/manifest.mpd"); err != nil {
		t.Fatalf("expected processing -> ready to pass: %v", err)
	}
}

func TestVideoRejectsInvalidTransition(t *testing.T) {
	video := videos.Video{Status: videos.StatusUploaded}

	if err := video.MarkReady("videos/123/manifest.mpd"); err == nil {
		t.Fatal("expected uploaded -> ready to fail")
	}
}
```

- [ ] **Step 2: Run the targeted tests to verify failure**

Run: `cd backend && go test ./internal/videos -run TestVideo -v`
Expected: FAIL with undefined `videos.Video` and status methods

- [ ] **Step 3: Add the domain entities and migration**

Create `backend/internal/videos/video.go`:

```go
package videos

import "fmt"

type Status string

const (
	StatusUploaded   Status = "uploaded"
	StatusProcessing Status = "processing"
	StatusReady      Status = "ready"
	StatusFailed     Status = "failed"
)

type Video struct {
	ID               string
	Title            string
	OriginalFilename string
	Status           Status
	SourceStorageKey string
	ManifestKey      string
	ErrorMessage     string
}

func (v *Video) MarkProcessing() error {
	if v.Status != StatusUploaded {
		return fmt.Errorf("cannot move %s to processing", v.Status)
	}
	v.Status = StatusProcessing
	return nil
}

func (v *Video) MarkReady(manifestKey string) error {
	if v.Status != StatusProcessing {
		return fmt.Errorf("cannot move %s to ready", v.Status)
	}
	v.Status = StatusReady
	v.ManifestKey = manifestKey
	v.ErrorMessage = ""
	return nil
}

func (v *Video) MarkFailed(message string) {
	v.Status = StatusFailed
	v.ErrorMessage = message
}
```

Create `backend/internal/jobs/job.go`:

```go
package jobs

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID           string
	VideoID      string
	Status       Status
	Attempts     int
	ErrorMessage string
}
```

Create `backend/migrations/0001_create_videos_and_jobs.sql`:

```sql
CREATE TABLE videos (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  original_filename TEXT NOT NULL,
  status TEXT NOT NULL,
  source_storage_key TEXT NOT NULL,
  manifest_storage_key TEXT,
  duration_seconds INTEGER,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transcoding_jobs (
  id TEXT PRIMARY KEY,
  video_id TEXT NOT NULL REFERENCES videos(id),
  status TEXT NOT NULL,
  attempts INTEGER NOT NULL DEFAULT 0,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- [ ] **Step 4: Run the targeted tests again**

Run: `cd backend && go test ./internal/videos -run TestVideo -v`
Expected: PASS

- [ ] **Step 5: Commit the core domain**

```bash
git add backend/migrations/0001_create_videos_and_jobs.sql backend/internal/videos/video.go backend/internal/jobs/job.go backend/internal/videos/service_test.go
git commit -m "feat: add video and job domain model"
```

### Task 3: Add the storage abstraction and local implementation

**Files:**
- Create: `backend/internal/storage/storage.go`
- Create: `backend/internal/storage/local/local.go`
- Create: `backend/internal/storage/s3/s3.go`
- Test: `backend/internal/storage/local/local_test.go`

- [ ] **Step 1: Write a failing test for local source and artifact persistence**

Create `backend/internal/storage/local/local_test.go`:

```go
package local_test

import (
	"bytes"
	"io"
	"testing"

	localstorage "spec-streaming/backend/internal/storage/local"
)

func TestLocalStorageSavesAndReadsArtifacts(t *testing.T) {
	store := localstorage.New(t.TempDir())

	if err := store.SaveArtifact("videos/1/manifest.mpd", bytes.NewBufferString("manifest")); err != nil {
		t.Fatalf("save artifact: %v", err)
	}

	rc, err := store.OpenArtifact("videos/1/manifest.mpd")
	if err != nil {
		t.Fatalf("open artifact: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}

	if string(data) != "manifest" {
		t.Fatalf("unexpected data: %s", string(data))
	}
}
```

- [ ] **Step 2: Run the local storage test to verify failure**

Run: `cd backend && go test ./internal/storage/local -run TestLocalStorageSavesAndReadsArtifacts -v`
Expected: FAIL because the package and constructor do not exist

- [ ] **Step 3: Add the storage contract and local adapter**

Create `backend/internal/storage/storage.go`:

```go
package storage

import "io"

type Storage interface {
	SaveSource(key string, reader io.Reader) error
	OpenSource(key string) (io.ReadCloser, error)
	SaveArtifact(key string, reader io.Reader) error
	OpenArtifact(key string) (io.ReadCloser, error)
	ArtifactExists(key string) (bool, error)
}
```

Create `backend/internal/storage/local/local.go`:

```go
package local

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	root string
}

func New(root string) *Storage {
	return &Storage{root: root}
}

func (s *Storage) SaveSource(key string, reader io.Reader) error { return s.writeFile(key, reader) }
func (s *Storage) OpenSource(key string) (io.ReadCloser, error) { return os.Open(filepath.Join(s.root, key)) }
func (s *Storage) SaveArtifact(key string, reader io.Reader) error { return s.writeFile(key, reader) }
func (s *Storage) OpenArtifact(key string) (io.ReadCloser, error) { return os.Open(filepath.Join(s.root, key)) }

func (s *Storage) ArtifactExists(key string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.root, key))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Storage) writeFile(key string, reader io.Reader) error {
	path := filepath.Join(s.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}
```

Create `backend/internal/storage/s3/s3.go` as an explicit MVP stub:

```go
package s3

import (
	"fmt"
	"io"
)

type Storage struct{}

func New() *Storage { return &Storage{} }

func (s *Storage) SaveSource(string, io.Reader) error            { return fmt.Errorf("not implemented") }
func (s *Storage) OpenSource(string) (io.ReadCloser, error)      { return nil, fmt.Errorf("not implemented") }
func (s *Storage) SaveArtifact(string, io.Reader) error          { return fmt.Errorf("not implemented") }
func (s *Storage) OpenArtifact(string) (io.ReadCloser, error)    { return nil, fmt.Errorf("not implemented") }
func (s *Storage) ArtifactExists(string) (bool, error)           { return false, fmt.Errorf("not implemented") }
```

- [ ] **Step 4: Run the storage tests again**

Run: `cd backend && go test ./internal/storage/... -v`
Expected: PASS for local storage tests

- [ ] **Step 5: Commit the storage layer**

```bash
git add backend/internal/storage/storage.go backend/internal/storage/local/local.go backend/internal/storage/s3/s3.go backend/internal/storage/local/local_test.go
git commit -m "feat: add storage abstraction"
```

### Task 4: Build the upload, list, detail, and status API

**Files:**
- Create: `backend/internal/platform/postgres/postgres.go`
- Create: `backend/internal/platform/http/router.go`
- Create: `backend/internal/platform/http/errors.go`
- Create: `backend/internal/videos/repository.go`
- Create: `backend/internal/jobs/repository.go`
- Create: `backend/internal/videos/service.go`
- Create: `backend/internal/jobs/service.go`
- Create: `backend/internal/videos/handler.go`
- Create: `backend/cmd/api/main.go`
- Test: `backend/internal/videos/handler_integration_test.go`

- [ ] **Step 1: Write a failing integration test for `POST /videos` and `GET /videos`**

Create `backend/internal/videos/handler_integration_test.go`:

```go
package videos_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/videos"
)

func TestUploadAndListVideos(t *testing.T) {
	e := echo.New()
	_ = bytes.NewBuffer(nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Trailer")
	part, _ := writer.CreateFormFile("file", "trailer.mp4")
	_, _ = part.Write([]byte("fake-mp4"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/videos", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}
```

- [ ] **Step 2: Run the integration test to verify failure**

Run: `cd backend && go test ./internal/videos -run TestUploadAndListVideos -v`
Expected: FAIL because no router or handlers are wired

- [ ] **Step 3: Implement repositories, services, and handlers minimally**

Create `backend/internal/videos/repository.go` with methods:

```go
package videos

import "context"

type Repository interface {
	Create(context.Context, *Video) error
	List(context.Context) ([]Video, error)
	GetByID(context.Context, string) (*Video, error)
	Update(context.Context, *Video) error
}
```

Create `backend/internal/jobs/repository.go` with methods:

```go
package jobs

import "context"

type Repository interface {
	Create(context.Context, *Job) error
	ClaimPending(context.Context) (*Job, error)
	Update(context.Context, *Job) error
}
```

Create `backend/internal/videos/service.go` with a use case shaped like:

```go
func (s *Service) CreateVideo(ctx context.Context, title string, filename string, file io.Reader) (*videos.Video, error)
func (s *Service) ListVideos(ctx context.Context) ([]videos.Video, error)
func (s *Service) GetVideo(ctx context.Context, id string) (*videos.Video, error)
```

Create `backend/internal/videos/handler.go` with handlers shaped like:

```go
func (h *Handler) Create(c echo.Context) error
func (h *Handler) List(c echo.Context) error
func (h *Handler) Get(c echo.Context) error
func (h *Handler) Status(c echo.Context) error
```

Create `backend/internal/platform/http/router.go`:

```go
package http

import (
	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/videos"
)

func NewRouter(videoHandler *videos.Handler) *echo.Echo {
	e := echo.New()
	e.POST("/videos", videoHandler.Create)
	e.GET("/videos", videoHandler.List)
	e.GET("/videos/:id", videoHandler.Get)
	e.GET("/videos/:id/status", videoHandler.Status)
	return e
}
```

- [ ] **Step 4: Run the backend integration tests**

Run: `cd backend && go test ./internal/videos -run TestUploadAndListVideos -v`
Expected: PASS

Run: `cd backend && go test ./...`
Expected: PASS

- [ ] **Step 5: Commit the API slice**

```bash
git add backend/internal/platform/postgres/postgres.go backend/internal/platform/http/router.go backend/internal/platform/http/errors.go backend/internal/videos/repository.go backend/internal/jobs/repository.go backend/internal/videos/service.go backend/internal/jobs/service.go backend/internal/videos/handler.go backend/internal/videos/handler_integration_test.go backend/cmd/api/main.go
git commit -m "feat(api): add video upload and listing endpoints"
```

### Task 5: Implement the worker and ffmpeg transcoding flow

**Files:**
- Create: `backend/internal/transcoding/transcoder.go`
- Create: `backend/internal/transcoding/ffmpeg.go`
- Create: `backend/internal/worker/runner.go`
- Create: `backend/cmd/worker/main.go`
- Test: `backend/internal/worker/runner_integration_test.go`

- [ ] **Step 1: Write a failing worker test for pending job processing**

Create `backend/internal/worker/runner_integration_test.go`:

```go
package worker_test

import "testing"

func TestRunnerProcessesPendingJob(t *testing.T) {
	t.Fatal("implement runner test")
}
```

- [ ] **Step 2: Run the worker test to verify failure**

Run: `cd backend && go test ./internal/worker -run TestRunnerProcessesPendingJob -v`
Expected: FAIL with `implement runner test`

- [ ] **Step 3: Replace the placeholder with a real runner test and implement the runner**

Replace `backend/internal/worker/runner_integration_test.go` with a test that asserts:

```go
func TestRunnerProcessesPendingJob(t *testing.T) {
	// Arrange an in-memory fake job repo with one pending job.
	// Arrange a fake video repo returning a matching uploaded video.
	// Arrange a fake transcoder returning manifest key "videos/1/manifest.mpd".
	// Run one cycle of the worker.
	// Assert the job becomes completed and the video becomes ready.
}
```

Create `backend/internal/transcoding/transcoder.go`:

```go
package transcoding

import "context"

type Service interface {
	Transcode(ctx context.Context, sourceKey string, videoID string) (string, error)
}
```

Create `backend/internal/worker/runner.go`:

```go
package worker

import "context"

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
	manifestKey, err := r.Codec.Transcode(ctx, video.SourceStorageKey, video.ID)
	if err != nil {
		video.MarkFailed(err.Error())
		return r.Videos.Update(ctx, video)
	}
	if err := video.MarkReady(manifestKey); err != nil {
		return err
	}
	return r.Videos.Update(ctx, video)
}
```

Create `backend/internal/transcoding/ffmpeg.go` with a shell-out implementation shape:

```go
package transcoding

import "os/exec"

func buildFFmpegCommand(input string, outDir string) *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-i", input,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-f", "dash",
		outDir+"/manifest.mpd",
	)
}
```

- [ ] **Step 4: Run all worker tests**

Run: `cd backend && go test ./internal/worker ./internal/transcoding -v`
Expected: PASS

- [ ] **Step 5: Commit the worker path**

```bash
git add backend/internal/transcoding/transcoder.go backend/internal/transcoding/ffmpeg.go backend/internal/worker/runner.go backend/internal/worker/runner_integration_test.go backend/cmd/worker/main.go
git commit -m "feat(worker): process transcoding jobs"
```

### Task 6: Serve DASH assets from the API

**Files:**
- Modify: `backend/internal/videos/handler.go`
- Modify: `backend/internal/platform/http/router.go`
- Test: `backend/internal/videos/handler_integration_test.go`

- [ ] **Step 1: Write a failing test for manifest serving**

Add to `backend/internal/videos/handler_integration_test.go`:

```go
func TestServeManifestForReadyVideo(t *testing.T) {
	// Arrange a ready video with manifest key videos/1/manifest.mpd.
	// Arrange local storage containing a simple manifest file.
	// Issue GET /videos/1/stream/manifest.mpd.
	// Expect 200 and body containing MPD markup.
}
```

- [ ] **Step 2: Run the manifest test to verify failure**

Run: `cd backend && go test ./internal/videos -run TestServeManifestForReadyVideo -v`
Expected: FAIL because stream routes do not exist

- [ ] **Step 3: Implement stream handlers and route registration**

Add to `backend/internal/videos/handler.go`:

```go
func (h *Handler) Manifest(c echo.Context) error
func (h *Handler) Segment(c echo.Context) error
```

Implementation rules:

- `Manifest` loads the video by ID and rejects non-`ready` videos.
- `Manifest` reads `manifest_storage_key` from storage and returns `application/dash+xml`.
- `Segment` resolves the requested relative path under the video artifact prefix and streams the file bytes.

Update `backend/internal/platform/http/router.go`:

```go
e.GET("/videos/:id/stream/manifest.mpd", videoHandler.Manifest)
e.GET("/videos/:id/stream/*", videoHandler.Segment)
```

- [ ] **Step 4: Run the streaming tests**

Run: `cd backend && go test ./internal/videos -run TestServeManifestForReadyVideo -v`
Expected: PASS

Run: `cd backend && go test ./...`
Expected: PASS

- [ ] **Step 5: Commit the streaming endpoints**

```bash
git add backend/internal/videos/handler.go backend/internal/platform/http/router.go backend/internal/videos/handler_integration_test.go
git commit -m "feat(api): serve dash manifests and segments"
```

### Task 7: Build the React upload, list, polling, and player UI

**Files:**
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/App.tsx`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/types.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/components/upload-form.tsx`
- Create: `frontend/src/components/video-list.tsx`
- Create: `frontend/src/components/status-badge.tsx`
- Create: `frontend/src/components/video-player.tsx`
- Create: `frontend/src/pages/home-page.tsx`
- Create: `frontend/src/pages/video-page.tsx`
- Test: `frontend/src/components/upload-form.test.tsx`
- Test: `frontend/src/components/video-list.test.tsx`

- [ ] **Step 1: Write failing component tests for upload and list rendering**

Create `frontend/src/components/video-list.test.tsx`:

```tsx
import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { VideoList } from './video-list'

describe('VideoList', () => {
  it('renders ready and processing videos', () => {
    render(
      <VideoList
        videos={[
          { id: '1', title: 'Trailer', status: 'ready' },
          { id: '2', title: 'Interview', status: 'processing' },
        ]}
      />,
    )

    expect(screen.getByText('Trailer')).toBeInTheDocument()
    expect(screen.getByText('Interview')).toBeInTheDocument()
  })
})
```

Replace `frontend/src/components/upload-form.test.tsx` with:

```tsx
import { describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import { UploadForm } from './upload-form'

describe('UploadForm', () => {
  it('submits title and file', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined)
    render(<UploadForm onSubmit={onSubmit} />)

    fireEvent.change(screen.getByLabelText(/title/i), {
      target: { value: 'Trailer' },
    })

    const file = new File(['video'], 'trailer.mp4', { type: 'video/mp4' })
    fireEvent.change(screen.getByLabelText(/file/i), {
      target: { files: [file] },
    })

    fireEvent.click(screen.getByRole('button', { name: /upload/i }))

    expect(onSubmit).toHaveBeenCalledWith({ title: 'Trailer', file })
  })
})
```

- [ ] **Step 2: Run the frontend component tests to verify failure**

Run: `cd frontend && npm test -- --run`
Expected: FAIL because the components do not exist yet

- [ ] **Step 3: Implement the frontend routes and components minimally**

Create `frontend/src/lib/types.ts`:

```ts
export type VideoStatus = 'uploaded' | 'processing' | 'ready' | 'failed'

export interface VideoItem {
  id: string
  title: string
  status: VideoStatus
  errorMessage?: string
}
```

Create `frontend/src/lib/api.ts` with methods:

```ts
export async function uploadVideo(input: { title: string; file: File }): Promise<VideoItem>
export async function listVideos(): Promise<VideoItem[]>
export async function getVideo(id: string): Promise<VideoItem>
```

Create `frontend/src/components/upload-form.tsx` with a controlled title field, file input, and `onSubmit` callback.

Create `frontend/src/components/video-list.tsx` with a prop shape:

```tsx
type Props = {
  videos: Array<{ id: string; title: string; status: string }>
}
```

Create `frontend/src/components/video-player.tsx` with a prop:

```tsx
type Props = { manifestUrl: string }
```

Implementation rule: start with a plain `<video controls src={manifestUrl}>` wrapper and defer DASH-specific player integration to a follow-up task only if native support proves insufficient.

- [ ] **Step 4: Add polling on the home page and verify tests pass**

Create `frontend/src/pages/home-page.tsx` with a `useEffect` loop that:

```tsx
useEffect(() => {
  let active = true

  async function refresh() {
    const next = await listVideos()
    if (active) setVideos(next)
  }

  refresh()
  const id = window.setInterval(refresh, 3000)
  return () => {
    active = false
    window.clearInterval(id)
  }
}, [])
```

Run: `cd frontend && npm test -- --run`
Expected: PASS

- [ ] **Step 5: Commit the frontend UI**

```bash
git add frontend/src/main.tsx frontend/src/App.tsx frontend/src/lib/api.ts frontend/src/lib/types.ts frontend/src/lib/format.ts frontend/src/components/upload-form.tsx frontend/src/components/video-list.tsx frontend/src/components/status-badge.tsx frontend/src/components/video-player.tsx frontend/src/pages/home-page.tsx frontend/src/pages/video-page.tsx frontend/src/components/upload-form.test.tsx frontend/src/components/video-list.test.tsx
git commit -m "feat(frontend): add upload list and player views"
```

### Task 8: Add Playwright E2E coverage for the happy path

**Files:**
- Create: `frontend/playwright.config.ts`
- Create: `frontend/tests/e2e/upload-and-play.spec.ts`
- Create: `frontend/tests/fixtures/sample.mp4`

- [ ] **Step 1: Write the failing end-to-end test**

Create `frontend/tests/e2e/upload-and-play.spec.ts`:

```ts
import { expect, test } from '@playwright/test'
import path from 'node:path'

test('uploads a video and opens the player page', async ({ page }) => {
  await page.goto('/')

  await page.getByLabel('Title').fill('Trailer')
  await page.getByLabel('File').setInputFiles(
    path.join(process.cwd(), 'tests/fixtures/sample.mp4'),
  )
  await page.getByRole('button', { name: 'Upload' }).click()

  await expect(page.getByText('Trailer')).toBeVisible()
  await expect(page.getByText('ready')).toBeVisible({ timeout: 30000 })

  await page.getByRole('link', { name: 'Trailer' }).click()
  await expect(page).toHaveURL(/\/videos\//)
  await expect(page.locator('video')).toBeVisible()
})
```

- [ ] **Step 2: Run Playwright to verify failure**

Run: `cd frontend && npx playwright test`
Expected: FAIL because the config, app boot, or backend dependencies are not ready yet

- [ ] **Step 3: Add the Playwright config with web server integration**

Create `frontend/playwright.config.ts`:

```ts
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  reporter: 'html',
  use: {
    baseURL: 'http://127.0.0.1:4173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'npm run dev -- --host 127.0.0.1 --port 4173',
    url: 'http://127.0.0.1:4173',
    reuseExistingServer: !process.env.CI,
  },
})
```

Implementation note: if the E2E flow needs the API and worker up as well, add a root-level helper command in `Makefile` that starts PostgreSQL, API, and worker before invoking `npx playwright test`.

- [ ] **Step 4: Run the full end-to-end suite**

Run: `make test-e2e`
Expected: PASS with one Chromium scenario and Playwright HTML report generated

- [ ] **Step 5: Commit the E2E coverage**

```bash
git add frontend/playwright.config.ts frontend/tests/e2e/upload-and-play.spec.ts frontend/tests/fixtures/sample.mp4 Makefile
git commit -m "test(e2e): add upload and playback flow"
```

### Task 9: Wire local development orchestration and verification commands

**Files:**
- Modify: `docker-compose.yml`
- Modify: `Makefile`
- Modify: `backend/cmd/api/main.go`
- Modify: `backend/cmd/worker/main.go`

- [ ] **Step 1: Write down the failing verification commands**

Run: `docker compose up -d`
Expected: FAIL because `docker-compose.yml` does not define PostgreSQL yet

Run: `make test-backend && make test-frontend && make test-e2e`
Expected: FAIL until full orchestration exists

- [ ] **Step 2: Add local runtime orchestration**

Create `docker-compose.yml`:

```yaml
services:
  postgres:
    image: postgres:17
    environment:
      POSTGRES_DB: spec_streaming
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - '5432:5432'
```

Extend `Makefile` with commands:

```make
run-api:
	cd backend && go run ./cmd/api

run-worker:
	cd backend && go run ./cmd/worker

verify:
	$(MAKE) test-backend
	$(MAKE) test-frontend
	$(MAKE) test-e2e
```

- [ ] **Step 3: Run the final verification commands**

Run: `docker compose up -d`
Expected: PASS with PostgreSQL container running

Run: `make verify`
Expected: PASS

- [ ] **Step 4: Capture the manual happy path**

Run:

```bash
docker compose up -d
make run-api
make run-worker
cd frontend && npm run dev
```

Expected manual result:

- upload page loads,
- MP4 upload returns created video,
- status transitions to `ready`,
- clicking the item opens the player page,
- the browser requests `manifest.mpd` via the API.

- [ ] **Step 5: Commit the local orchestration**

```bash
git add docker-compose.yml Makefile backend/cmd/api/main.go backend/cmd/worker/main.go
git commit -m "chore: add local development workflow"
```

## Self-Review

- Spec coverage: upload, transcoding worker, MPEG-DASH serving, video list, player page, polling, local and S3-compatible storage abstraction, and Playwright E2E are each mapped to at least one task.
- Placeholder scan: the only intentional stub is the S3 adapter implementation, which is explicitly scoped as a non-MVP adapter shell while preserving the interface promised by the spec.
- Type consistency: `Video`, `Job`, status names, storage methods, and stream route paths are kept consistent across tasks.
