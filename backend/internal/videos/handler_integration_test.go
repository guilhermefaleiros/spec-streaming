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

func TestUploadAndListVideos(t *testing.T) {
	repo := newFakeVideoRepo()
	storage := local.New(t.TempDir())
	service := videos.NewService(repo, storage)
	handler := videos.NewHandler(service)
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
