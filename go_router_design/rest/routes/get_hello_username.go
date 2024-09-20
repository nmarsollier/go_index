package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_router_design/model/hello/service"
	"github.com/nmarsollier/go_router_design/utils/errors"
)

// Internal configure ping/pong service
func init() {
	router().GET(
		"/hello/:userName",
		validateUserName,
		sayHelloHandler,
	)
}

// validamos que el parametro userName tenga al menos 5 caracteres
func validateUserName(c *gin.Context) {
	userName := c.Param("userName")

	if len(userName) < 5 {
		c.Error(errors.NewCustomError(400, "userName debe tener al menos 5 caracteres"))
		c.Abort()
		return
	}
}

func sayHelloHandler(c *gin.Context) {
	userName := c.Param("userName")

	c.JSON(http.StatusOK, gin.H{
		"answer": service.SayHello(userName),
	})
}
