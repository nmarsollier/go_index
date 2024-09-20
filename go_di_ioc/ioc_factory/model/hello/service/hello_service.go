package service

import "github.com/nmarsollier/go_di_ioc/ioc_factory/model/hello/dao"

// IHelloDao interface DAO necesaria a inyectar en el service
type IHelloDao interface {
	Hello() string
}

// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}

// NewService es una función que puede mockearse
func NewService() *HelloService {
	return &HelloService{
		dao.NewDao(),
	}
}

// SayHello es nuestro metodo de negocio
func (s *HelloService) SayHello() string {
	return s.dao.Hello()
}
