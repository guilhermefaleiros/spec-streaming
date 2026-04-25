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
	e.GET("/videos/:id/stream/manifest.mpd", videoHandler.Manifest)
	e.GET("/videos/:id/stream/*", videoHandler.Segment)
	return e
}
