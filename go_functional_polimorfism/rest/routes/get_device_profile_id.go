package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_functional_polimorfism/model/profile"
	"github.com/nmarsollier/go_functional_polimorfism/utils/errors"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/:device/profile/:id",
		validateUserName,
		getProfile,
	)
}

// Hacemos las validaciones que nos aseguran que id es valido
func validateUserName(c *gin.Context) {
	device := c.Param("device")

	if !profile.IsValidDevice(device) {
		c.Error(errors.NewCustomError(400, "device debe ser mobile o web"))
		c.Abort()
		return
	}
}

func getProfile(c *gin.Context) {
	device := c.Param("device")
	id := c.Param("id")

	data := profile.FetchProfile(id)
	image := profile.GetImage[device](id)

	c.JSON(http.StatusOK, gin.H{
		"login": data.Login,
		"web":   data.Web,
		"name":  data.Name,
		"image": image,
	})
}
