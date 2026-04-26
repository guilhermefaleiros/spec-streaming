package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"spec-streaming/backend/internal/jobs"
	"spec-streaming/backend/internal/platform/http"
	"spec-streaming/backend/internal/storage/local"
	"spec-streaming/backend/internal/videos"
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
