package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_rest_controller/model/hello/service"
)

// Internal configure ping/pong service
func init() {
	getRouter().GET("/hello/:userName", sayHelloHandler)
}

func sayHelloHandler(c *gin.Context) {
	userName := c.Param("userName")

	c.JSON(http.StatusOK, gin.H{
		"answer": service.SayHello(userName),
	})
}
