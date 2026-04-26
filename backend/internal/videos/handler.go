package videos

import (
	"io"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/storage"
)

type Handler struct {
	service *Service
	storage storage.Storage
}

type response struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	OriginalFilename string `json:"originalFilename"`
	Status           string `json:"status"`
	SourceStorageKey string `json:"sourceStorageKey"`
	ManifestKey      string `json:"manifestKey"`
	ErrorMessage     string `json:"errorMessage"`
}

func toResponse(video *Video) response {
	return response{
		ID:               video.ID,
		Title:            video.Title,
		OriginalFilename: video.OriginalFilename,
		Status:           string(video.Status),
		SourceStorageKey: video.SourceStorageKey,
		ManifestKey:      video.ManifestKey,
		ErrorMessage:     video.ErrorMessage,
	}
}

func toResponses(videos []Video) []response {
	items := make([]response, 0, len(videos))
	for i := range videos {
		items = append(items, toResponse(&videos[i]))
	}
	return items
}

func NewHandler(service *Service, storage storage.Storage) *Handler {
	return &Handler{service: service, storage: storage}
}

func (h *Handler) Create(c echo.Context) error {
	title := c.FormValue("title")
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot open file"})
	}
	defer src.Close()

	video, err := h.service.CreateVideo(c.Request().Context(), title, file.Filename, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, toResponse(video))
}

func (h *Handler) List(c echo.Context) error {
	videos, err := h.service.ListVideos(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, toResponses(videos))
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	video, err := h.service.GetVideo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, toResponse(video))
}

func (h *Handler) Status(c echo.Context) error {
	id := c.Param("id")
	video, err := h.service.GetVideo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": string(video.Status)})
}

func (h *Handler) Manifest(c echo.Context) error {
	id := c.Param("id")
	video, err := h.service.GetVideo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if video.Status != StatusReady {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "video not ready"})
	}

	rc, err := h.storage.OpenArtifact(video.ManifestKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot open manifest"})
	}
	defer rc.Close()

	c.Response().Header().Set(echo.HeaderContentType, "application/dash+xml")
	_, err = io.Copy(c.Response().Writer, rc)
	return err
}

func (h *Handler) Segment(c echo.Context) error {
	id := c.Param("id")
	video, err := h.service.GetVideo(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if video.Status != StatusReady {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "video not ready"})
	}

	segmentPath := path.Join(path.Dir(video.ManifestKey), c.Param("*"))
	rc, err := h.storage.OpenArtifact(segmentPath)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "segment not found"})
	}
	defer rc.Close()

	_, err = io.Copy(c.Response().Writer, rc)
	return err
}
