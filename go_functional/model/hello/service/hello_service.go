package service

import "github.com/nmarsollier/go_di_ioc/go_funcional/model/hello/dao"

// Nos va a permitir mockear respuestas para los tests
var sayHelloFunc func() string = dao.Hello

// SayHello es nuestro negocio
func SayHello() string {
	return sayHelloFunc()
}
