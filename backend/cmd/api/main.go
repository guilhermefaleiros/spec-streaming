package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"spec-streaming/backend/internal/platform/http"
)

func main() {
	e := echo.New()
	_ = http.NewRouter(nil) // placeholder until wiring is complete
	log.Println("API starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
