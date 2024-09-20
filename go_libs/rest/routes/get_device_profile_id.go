package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_libs/model/profile"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/:device/profile/:id",
		getProfile,
	)
}

func getProfile(c *gin.Context) {
	id := c.Param("id")

	data := profile.FetchProfile(id)

	c.JSON(http.StatusOK, gin.H{
		"login": data.Login(),
		"web":   data.Web(),
		"name":  data.Name(),
	})
}
