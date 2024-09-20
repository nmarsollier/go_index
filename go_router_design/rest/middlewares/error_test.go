package middlewares

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	errutils "github.com/nmarsollier/go_router_design/utils/errors"
	"github.com/nmarsollier/go_router_design/utils/test"
)

func TestCustomError(t *testing.T) {
	response := test.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)

	context.Error(errutils.NewCustomError(400, "Custom Test"))
	handleErrorIfNeeded(context)

	response.Assert(400, "{\"error\":\"Custom Test\"}")
}

func TestError(t *testing.T) {
	response := test.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)

	context.Error(errors.New("Error Test"))
	handleErrorIfNeeded(context)

	response.Assert(500, "{\"error\":\"Error Test\"}")
}
