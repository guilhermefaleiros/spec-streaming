package main

import (
	"log"
	"os"

	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/platform/http"
	"spec-streaming/backend/internal/storage/local"
	"spec-streaming/backend/internal/videos"
)

func main() {
	storage := local.New("tmp/storage")
	videoRepo := videos.NewMemoryRepository()
	jobRepo := jobs.NewMemoryRepository()

	jobService := jobs.NewService(jobRepo)
	videoService := videos.NewService(videoRepo, storage, jobService)

	videoHandler := videos.NewHandler(videoService, storage)

	e := http.NewRouter(videoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API starting on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
