package dao

// HelloDao representa un DAO podria ser una base de datos o un acceso a API externo
type HelloDao struct {
}

// NewDao es el factory
func NewDao() *HelloDao {
	return new(HelloDao)
}

// Hello es nuestro metodo de negocio
func (d *HelloDao) Hello() string {
	return "Hello"
}
