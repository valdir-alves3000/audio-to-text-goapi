package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/handler/http"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/utils"
)

func main() {
	if err := util.CheckDependencies(); err != nil {
		log.Fatalf("Dependency error: %v", err)
	}

	if err := os.MkdirAll("models", 0755); err != nil {
		log.Fatalf("Failed to create models directory: %v", err)
	}

	r := gin.Default()
	http.RegisterRoutes(r)

	log.Println("Server running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
