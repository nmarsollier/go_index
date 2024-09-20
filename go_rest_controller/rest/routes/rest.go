package routes

import (
	"github.com/gin-gonic/gin"
)

// Start server in 8080 port
func Start() {
	getRouter().Run(":8080")
}

var router *gin.Engine = nil

func getRouter() *gin.Engine {
	if router == nil {
		router = gin.Default()
	}

	return router
}
