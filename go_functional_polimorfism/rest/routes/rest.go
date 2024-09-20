package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_functional_polimorfism/rest/middlewares"
)

// Start server in 8080 port
func Start() {
	router().Run(":8080")
}

var engine *gin.Engine = nil

func router() *gin.Engine {
	if engine == nil {
		engine = gin.Default()
		engine.Use(middlewares.ErrorHandler)
	}

	return engine
}
