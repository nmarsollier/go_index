package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_cache/model/profile"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/profile",
		getProfile,
	)
}

func getProfile(c *gin.Context) {
	data := profile.FetchProfile("123")

	if data == nil {
		c.AbortWithError(500, errors.New("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login": data.Login,
		"web":   data.Web,
		"name":  data.Name,
	})
}
