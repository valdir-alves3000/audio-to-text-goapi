package http

import (
	"github.com/gin-gonic/gin"
)

func noCache(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Surrogate-Control", "no-store")
	c.Next()
}

func RegisterRoutes(r *gin.Engine) {
	r.Use(noCache)
	r.Static("/static", "./web/static")
	r.GET("/", HomePageHandler)
	r.GET("/docs", DocsHandler)

	api := r.Group("/api")
	{
		api.POST("/transcribe", TranscribeHandler)
	}
}
