package videos_test

import (
	"bytes"
	"context"
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
	service := videos.NewService(repo, storage)
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
	service := videos.NewService(repo, storage)
	handler := videos.NewHandler(service, storage)
	e := echo.New()
	e.POST("/videos", handler.Create)

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
