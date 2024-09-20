package service

// IHelloDao interface DAO necesaria a inyectar en el service
type IHelloDao interface {
	Hello() string
}

// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}

// NewService Es un factory del servicio de Negocio , depende de IHelloDao
func NewService(dao IHelloDao) *HelloService {
	return &HelloService{dao}
}

// SayHello es nuestro metodo de negocio
func (s *HelloService) SayHello() string {
	return s.dao.Hello()
}
