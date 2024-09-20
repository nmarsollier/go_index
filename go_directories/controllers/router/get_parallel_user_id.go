package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_directories/model/profile"
	"github.com/nmarsollier/go_directories/model/user"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/parallel/users/:id",
		validateUserName,
		inParallel(
			fetchUserInParallel,
			fetchProfileInParallel,
		),
		build,
	)
}

func fetchUserInParallel(c *gin.Context) {
	c.Set("user", user.FetchUser())
}

func fetchProfileInParallel(c *gin.Context) {
	c.Set("profile", profile.FetchProfile())
}

func inParallel(handlers ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(handlers))

		for _, handler := range handlers {
			go func(handlerFunc gin.HandlerFunc) {
				defer waitGroup.Done()
				handlerFunc(c)
			}(handler)
		}

		waitGroup.Wait()

		c.Next()
	}
}
