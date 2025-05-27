package http

import (
	"github.com/gin-gonic/gin"
)


func HomePageHandler(c *gin.Context) {
	c.File("./web/index.html")
}