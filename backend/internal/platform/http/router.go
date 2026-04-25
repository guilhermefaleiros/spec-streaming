package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"spec-streaming/backend/internal/videos"
)

func NewRouter(videoHandler *videos.Handler) *echo.Echo {
	e := echo.New()

	// CORS for frontend development
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.POST("/videos", videoHandler.Create)
	e.GET("/videos", videoHandler.List)
	e.GET("/videos/:id", videoHandler.Get)
	e.GET("/videos/:id/status", videoHandler.Status)
	e.GET("/videos/:id/stream/manifest.mpd", videoHandler.Manifest)
	e.GET("/videos/:id/stream/*", videoHandler.Segment)
	return e
}
