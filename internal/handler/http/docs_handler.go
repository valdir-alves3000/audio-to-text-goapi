package http

import "github.com/gin-gonic/gin"

func DocsHandler(c *gin.Context) {
	c.File("./web/docs.html")
}
