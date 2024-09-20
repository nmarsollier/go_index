package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_declarative/service"
	"github.com/nmarsollier/go_declarative/utils/errors"
	"github.com/nmarsollier/go_declarative/utils/gu"
)

// Internal configure ping/pong service
func init() {
	router().GET(
		"/hello/:userName",
		validateUserName,
		sayHello,
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

func sayHello(c *gin.Context) {
	gu.SendJSONAnswer(c, service.SayHello(c.Param("userName")))
}
