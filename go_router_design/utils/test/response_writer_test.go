package test

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseWriter(t *testing.T) {
	response := ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	context.JSON(500, gin.H{"error": "Internal server error"})
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}
