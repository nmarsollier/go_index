package service

import (
	"fmt"

	"github.com/nmarsollier/go_declarative/dao"
)

// Nos va a permitir mockear respuestas para los tests
// en producción deberia funcionar muy bien, porque solo
// los tests pueden cambiar este valor
var daoGetHello func() string = dao.GetHello

// SayHello es nuestro negocio
func SayHello(userName string) string {
	return fmt.Sprintf("%s %s", daoGetHello(), userName)
}
