package main

import (
	"log"
	"time"

	"spec-streaming/backend/internal/worker"
)

func main() {
	runner := &worker.Runner{}
	log.Println("Worker starting...")
	for {
		if err := runner.RunOnce(nil); err != nil {
			log.Printf("runner error: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}
