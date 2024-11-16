<!-- cSpell:language es -->

[English version](README_en.md)

# Las librerías compartidas

Cuando implementamos un microservicio puntual, normalmente esta implementación debería resolver un problema puntual de nuestro negocio, y su interfaz con el mundo exterior esta definida a través de rest, message broker o algo similar. Su definición interna no es muy importante a los ojos de los usuarios, y podemos darnos el lujo de ser muy flexibles como lo venimos viendo en esta guía.

Pero cuando implementamos librerías compartidas, como pueden ser:

- Factories para servicios remotos
- Log de errores
- Apis de acceso (brokers de mensajería o bases de datos)
- Acceso internacionalizado
- Acceso a imágenes y recursos compartidos
- Trazabilidad
- Utilidades en general
- Utilidades de testing (http o bases de datos mock)

Estas librerías, tienen como funcionalidad no repetir código en los microservicios y ayudan a los equipos de desarrollo a centralizar decisiones y controlar lógica común.

Pero no se desarrollan como el código que hemos venido trabajando en esta guía, sino que debemos hacerlas mas robustas para que puedan ser mas flexibles y adaptarse a nuevas implementaciones.

## Interfaces y estructuras

Las estructuras en go definen datos relacionados en un contexto puntual.

Por ejemplo

```go
type Profile struct {
	ID    string
	Login string
	Name  string
	Web   string
}
```

Profile representa el perfil de un usuario.

Las interfaces nos permiten definir comportamientos en forma desacoplada de las implementaciones. Por ejemplo, si definimos :

```go
type Profile interface {
	ID() string
	Login() string
	Name() string
	Web() string
}

func FetchProfile(id string) Profile {
   ...
}
```

Estamos definiendo una serie de comportamientos que se deben cumplir en las implementaciones, y que existe una función FetchProfile, que retorna esos comportamientos.

Las implementaciones de las interfaces se hacen con estructuras y métodos asociadas a ellas :

```go
type profile struct {
	id    string
	login string
	name  string
	web   string
}

func (p *profile) ID() string {
	return p.id
}

func (p *profile) Login() string {
	return p.login
}

func (p *profile) Name() string {
	return p.name
}

func (p *profile) Web() string {
	return p.web
}
```

Entonces podemos definir la funcion FetchProfile como :

```go
func FetchProfile(id string) Profile {
	return &profile{
		id:    id,
		login: "nmarsollier",
		name:  "Nestor Marsollier",
		web:   "https://github.com/nmarsollier/profile",
	}
}
```

> El uso de interfaces es lo que podríamos denominar polimorfismo en Go, ya que podemos implementar diferentes estructuras con diferentes comportamientos, pero respetando un set de funciones determinado.

## Y todo esto para que ?

El esquema anterior es muy importante en el uso de librerías compartidas, ya que nos permite definir un comportamiento esperado, pero no los detalles de implementación interna, esto nos deja cierta flexibilidad a la hora de adaptarnos a cambios y a su vez a la hora de implementaciones personalizadas de cada negocio.

Encapsulamos yu desacoplamos, abrimos posibilidades a otras implementaciones, podemos lograr mejor testing pudiendo mockear mejor llamadas y resultados.

## Atención a las interfaces implícitas

En go las interfaces son implícitas, lo que quiere decir que no debemos ser específicos de que estructura implementa una interfaz, simplemente con que la estructura defina todos los métodos de la interfaz, ya la esta cumpliendo

> Debemos ser específicos con los nombres de funciones en las interfaces. Si nuestra función se llama 'Serialize() string' en dos interfaces diferentes XMLSerializer y JSONSerializer, podemos incurrir en un error, porque para el compilador si una estructura implementa Serialize, ya esta implementando ambas interfaces.

## Donde definimos las interfaces ?

Normalmente se definen en las librerías, a veces en las implementaciones de los microservicios necesitamos ser mas flexibles, y requieren de una definición de que es lo que necesitan para operar, que básicamente es un subset muy reducido de lo que la interfaz de la librería expone, y queremos generalizar para que pueda operar con otras opciones.

Por ejemplo, un cliente de base de datos o cliente http, podría exponer muchos métodos, pero a nuestros DAOs no queremos habilitar todos los métodos que pueda exponer un cliente de bases de datos, por lo que creamos una interfaz con un subset de lo que necesitamos en nuestra app, y operamos con ese subset, para evitar usar métodos que no deseamos utilizar del cliente original.

> En algunas bibliografías se recomienda definir la interfaz en cada función, de modo que en cada punto que se define una función se definen las interfaces necesarias, si bien esto parece que nos da flexibilidad, en realidad, terminamos teniendo demasiadas interfaces y muchos código que mantener, no lo recomiendo en lo personal.

## Reglas de las interfaces

> Robustness principle: Be conservative with what you do, be liberal with you accept

Nuestras funciones deberían poder funcionar correctamente con cualquier estructura de código que implemente una interfaz, sin embargo, debemos usar solo lo que nuestra interfaz define. Nunca debemos hacer interfaces que definan mas de lo que usamos.

> Accept interfaces, return structs

Las funciones deben definirse con parámetros y retornos en Interfaces, sin embargo nuestro código retorna una implementación especifica de una estructura

```go
func FetchProfile(id string) Profile {
	return &profile{
		id:    id,
		login: "nmarsollier",
		name:  "Nestor Marsollier",
		web:   "https://github.com/nmarsollier/profile",
	}
}
```

> The empty interface says nothing

Limitamos el uso de interface{} al mínimo posible, solo lo usamos cuando no podemos mapear los resultados de una librería de terceros, o bien en alguna implementación muy puntual como mapas para cache.

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index/blob/main/README.md)
