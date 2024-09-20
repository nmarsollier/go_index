package routes

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_router_design/utils/test"
	"gopkg.in/go-playground/assert.v1"
)

func TestValidateUserName(t *testing.T) {
	response := test.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/hello/abc", nil)

	validateUserName(context)

	response.Assert(0, "")
	assert.Equal(t, context.Errors.Last().Error(), "userName debe tener al menos 5 caracteres")
}
