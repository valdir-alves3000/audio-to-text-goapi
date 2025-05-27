package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Static("/static", "./web/static")
	r.GET("/", HomePageHandler)

	api := r.Group("/api")
	{
		api.POST("/transcribe", TranscribeHandler)
	}
}
