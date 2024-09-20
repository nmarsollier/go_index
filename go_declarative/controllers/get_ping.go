package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_declarative/utils/gu"
)

// Internal configure ping/pong service
func init() {
	router().GET("/ping", ping)
}

func ping(c *gin.Context) {
	gu.SendJSONAnswer(c, "pong")
}
