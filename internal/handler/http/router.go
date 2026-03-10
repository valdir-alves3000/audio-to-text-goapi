package http

import (
	"github.com/gin-gonic/gin"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/process"
)

func noCache(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Surrogate-Control", "no-store")
	c.Next()
}

func RegisterRoutes(r *gin.Engine, worker *process.WhisperWorker) {

	r.Use(noCache)

	r.Static("/static", "./web/static")
	r.GET("/", HomePageHandler)
	r.GET("/docs", DocsHandler)

	api := r.Group("/api")
	{
		api.POST("/transcribe", NewTranscribeHandler(worker))
	}
}