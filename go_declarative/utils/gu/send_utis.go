package gu

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendJSONAnswer sends an answer as JSON
func SendJSONAnswer(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"answer": data,
	})
}
