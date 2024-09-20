package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Internal configure ping/pong service
func init() {
	getRouter().GET("/ping", pingHandler)
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"answer": "pong",
	})
}
