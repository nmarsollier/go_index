package dao

// HelloDao Es nuestra implementacion de Dao
type HelloDao struct {
}

// NewDao es el factory
func NewDao() *HelloDao {
	return new(HelloDao)
}

// Hello es nuestro metodo de negocio
func (d *HelloDao) Hello() string {
	return "Holiiis"
}
