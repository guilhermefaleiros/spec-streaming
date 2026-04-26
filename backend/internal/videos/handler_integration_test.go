package videos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/storage/local"
	"spec-streaming/backend/internal/videos"
)

type fakeVideoRepo struct {
	videos map[string]*videos.Video
}

func newFakeVideoRepo() *fakeVideoRepo {
	return &fakeVideoRepo{videos: make(map[string]*videos.Video)}
}

func (r *fakeVideoRepo) Create(ctx context.Context, v *videos.Video) error {
	r.videos[v.ID] = v
	return nil
}

func (r *fakeVideoRepo) List(ctx context.Context) ([]videos.Video, error) {
	var list []videos.Video
	for _, v := range r.videos {
		list = append(list, *v)
	}
	return list, nil
}

func (r *fakeVideoRepo) GetByID(ctx context.Context, id string) (*videos.Video, error) {
	if v, ok := r.videos[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (r *fakeVideoRepo) Update(ctx context.Context, v *videos.Video) error {
	r.videos[v.ID] = v
	return nil
}

func TestServeManifestForReadyVideo(t *testing.T) {
	repo := newFakeVideoRepo()
	storage := local.New(t.TempDir())
	service := videos.NewService(repo, storage, nil)
	handler := videos.NewHandler(service, storage)
	e := echo.New()
	e.GET("/videos/:id/stream/manifest.mpd", handler.Manifest)

	// Arrange a ready video with manifest key
	video := &videos.Video{
		ID:          "vid-1",
		Status:      videos.StatusReady,
		ManifestKey: "videos/vid-1/manifest.mpd",
	}
	repo.videos[video.ID] = video

	// Arrange local storage containing a simple manifest file
	if err := storage.SaveArtifact(video.ManifestKey, bytes.NewBufferString("<MPD>manifest</MPD>")); err != nil {
		t.Fatalf("save artifact: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/videos/vid-1/stream/manifest.mpd", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("<MPD>")) {
		t.Fatalf("expected MPD markup, got %s", rec.Body.String())
	}
}

func TestUploadAndListVideos(t *testing.T) {
	repo := newFakeVideoRepo()
	storage := local.New(t.TempDir())
	service := videos.NewService(repo, storage, nil)
	handler := videos.NewHandler(service, storage)
	e := echo.New()
	e.POST("/videos", handler.Create)
	e.GET("/videos", handler.List)

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

	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}
	if _, ok := created["id"]; !ok {
		t.Fatalf("expected camelCase create response, got %#v", created)
	}
	if _, ok := created["ID"]; ok {
		t.Fatalf("expected uppercase create keys to be absent, got %#v", created)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/videos", nil)
	listRec := httptest.NewRecorder()
	e.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from list, got %d", listRec.Code)
	}

	var listed []map[string]any
	if err := json.Unmarshal(listRec.Body.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list response: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected 1 listed video, got %d", len(listed))
	}
	if _, ok := listed[0]["status"]; !ok {
		t.Fatalf("expected camelCase list response, got %#v", listed[0])
	}
	if _, ok := listed[0]["Status"]; ok {
		t.Fatalf("expected uppercase list keys to be absent, got %#v", listed[0])
	}
}

func TestListVideosReturnsCamelCaseJSON(t *testing.T) {
	repo := newFakeVideoRepo()
	storage := local.New(t.TempDir())
	service := videos.NewService(repo, storage, nil)
	handler := videos.NewHandler(service, storage)
	e := echo.New()
	e.GET("/videos", handler.List)

	repo.videos["vid-1"] = &videos.Video{
		ID:               "vid-1",
		Title:            "Trailer",
		OriginalFilename: "trailer.mp4",
		Status:           videos.StatusReady,
		SourceStorageKey: "sources/vid-1",
		ManifestKey:      "artifacts/vid-1/manifest.mpd",
		ErrorMessage:     "",
	}

	req := httptest.NewRequest(http.MethodGet, "/videos", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(payload) != 1 {
		t.Fatalf("expected 1 video, got %d", len(payload))
	}

	video := payload[0]
	if video["id"] != "vid-1" {
		t.Fatalf("expected camelCase id, got %#v", video)
	}
	if video["title"] != "Trailer" {
		t.Fatalf("expected camelCase title, got %#v", video)
	}
	if video["status"] != "ready" {
		t.Fatalf("expected camelCase status, got %#v", video)
	}
	if _, ok := video["ID"]; ok {
		t.Fatalf("expected uppercase keys to be absent, got %#v", video)
	}
}

func TestGetVideoReturnsCamelCaseJSON(t *testing.T) {
	repo := newFakeVideoRepo()
	storage := local.New(t.TempDir())
	service := videos.NewService(repo, storage, nil)
	handler := videos.NewHandler(service, storage)
	e := echo.New()
	e.GET("/videos/:id", handler.Get)

	repo.videos["vid-1"] = &videos.Video{
		ID:               "vid-1",
		Title:            "Trailer",
		OriginalFilename: "trailer.mp4",
		Status:           videos.StatusReady,
		SourceStorageKey: "sources/vid-1",
		ManifestKey:      "artifacts/vid-1/manifest.mpd",
	}

	req := httptest.NewRequest(http.MethodGet, "/videos/vid-1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal get response: %v", err)
	}

	if payload["id"] != "vid-1" {
		t.Fatalf("expected camelCase id, got %#v", payload)
	}
	if payload["status"] != "ready" {
		t.Fatalf("expected camelCase status, got %#v", payload)
	}
	if _, ok := payload["ID"]; ok {
		t.Fatalf("expected uppercase keys to be absent, got %#v", payload)
	}
}
