package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	httpHandler "github.com/valdir-alves3000/audio-to-text-goapi/internal/handler/http"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/process"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/utils"
)

func main() {
	if err := util.CheckDependencies(); err != nil {
		log.Fatalf("Dependency error: %v", err)
	}

	if err := os.MkdirAll("models", 0755); err != nil {
		log.Fatalf("Failed to create models directory: %v", err)
	}

	worker, err := process.GetWhisperWorker()
	if err != nil {
		log.Fatalf("Failed to start Whisper worker: %v", err)
	}

	r := gin.Default()

	httpHandler.RegisterRoutes(r, worker)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	url := fmt.Sprintf("http://localhost:%s", port)

	log.Println("====================================")
	log.Printf("🚀 Server running at %s\n", url)
	log.Println("====================================")

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}