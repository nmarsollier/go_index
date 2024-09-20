package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_directories/model/profile"
	"github.com/nmarsollier/go_directories/model/user"
	"github.com/nmarsollier/go_directories/utils/errors"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/users/:id",
		validateUserName,
		fetchUser,
		fetchProfile,
		build,
	)
}

// Hacemos las validaciones que nos aseguran que id es valido
func validateUserName(c *gin.Context) {
	id := c.Param("id")

	if len(id) < 1 {
		c.Error(errors.NewCustomError(400, "userName debe tener al menos 1 caracteres"))
		c.Abort()
		return
	}
}

func fetchUser(c *gin.Context) {
	c.Set("user", user.FetchUser())
	c.Next()
}

func fetchProfile(c *gin.Context) {
	c.Set("profile", profile.FetchProfile())
	c.Next()
}

func build(c *gin.Context) {
	user := c.MustGet("user").(*user.User)
	profile := c.MustGet("profile").(*profile.Profile)

	c.JSON(http.StatusOK, gin.H{
		"login":  user.Login,
		"access": user.Access,
		"name":   profile.Name,
	})
}
