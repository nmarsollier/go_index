<!-- cSpell:language es -->

[English version](README_en.md)

# DI e IoC en GO

Este repositorio plantea alternativas de manejo de dependencias, a la programación tradicional de un proyecto Go.

Mas adelante se encuentra un ejemplo de como manejar un service locator de forma adecuada, cuando tiene sentido.

## Mal uso de Inyección de Dependencias

En ese uso tradicional de IoC que nos permite proveer dependencias en una estructura de datos para que sean usadas internamente.

En la carpeta [ejemplo_tradicional](./ejemplo_tradicional/) tenemos los ejemplos de código.

La mayoría de los programadores recomiendan inyección de dependencias para separar capas lógicas de código, y desacoplar los servicios de los clientes.

En Go la estrategia mas común es la de Inyección de Dependencias pasada por Constructor.

Nuestro código luce como el siguiente:

```go
srv := service.NewService(dao.NewDao())
fmt.Println(srv.SayHello())
```

donde Service es algo como lo siguiente :

```go
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

// SayHello es nuestro método de negocio
func (s HelloService) SayHello() string {
	return s.dao.Hello()
}
```

Y esa es una buena practica, pero cuando las dependencias son realmente implementaciones que se ajusten a casos donde si se debe usar inyección de dependencias (ejemplo: Strategy pattern), aclarando que esta estrategia no está mal, el problema es cuando abusamos de ella en casos donde no es correcto usarse, en nuestros sistemas, la gran mayoría de dependencias no necesitan usar DI.

Un uso como el anterior no desacopla realmente, todo lo contrario, terminamos acoplando mucho mas, nuestro código debe definir métodos bootstrap en lugares donde no deberían estar, acoplando todo el negocio en un archivo main.go por ejemplo.

Ademas de acoplar, exponemos los detalles internos de nuestras implementaciones, algo que debería estar encapsulado.

Podemos revisar las raíces de programación orientada a objetos y darnos cuenta que incluso iría en contra del principio de [Information Expert](<https://en.wikipedia.org/wiki/GRASP_(object-oriented_design)#Information_expert>), que nos dice que las dependencias deben ser creadas en el lugar donde se tiene la información, de modo tal que por ejemplo crear las dependencias de un servicio en la capa de controladores, genera acoplamiento y nuestro código ya no seria cohesivo.

La realidad, es que se usa porque nos da la posibilidad de mockear esas dependencias para realizar buenos testing. Sin embargo existen muchas alternativas que podemos usar sin llegar a usar DI.

## Uso de Factory Methods como IoC

Veamos como podemos mejorar la situación anterior.

Inversión de control básicamente significa tener un framework, que cuando necesite un recurso se lo pida a él, y el recurso se obtiene del contexto.

Un service locator, es básicamente un framework que conoce nuestras dependencias, y las inyecta en donde sea necesario, pero tiene el mismo problema que el bootstrap anterior, acopla todos los servicios en un solo lugar.

Si partimos los patrones generales de asignación de responsabilidades GRASP, una de las formas clásicas y adecuadas de construir objetos es el uso de Factory Methods.

Pensemos en ese factory method, como parte de un framework de inyección de dependencias, que dependiendo del contexto (parámetros en el factory u otro método) nos va a retornar la instancia correcta del servicio que necesitemos. Lo que tiene de adecuado este patrón, es que la estrategia de creación, se escribe junto a los servicios, por lo que queda mucho mas claro el funcionamiento del mismo.

Esta estrategia nos permite evitar inyectar las dependencias en los constructores y delegar la creación a funciones factory que están en el lugar adecuado.

Este ejemplo lo encontramos en [ioc_factory](./ioc_factory/)

Como vemos en la función main: la creación del service no esta acoplada a la creación del dao.

```go
srv := service.NewService()

fmt.Println(srv.SayHello())
```

Sino mas bien el mismo service se encarga de crear el dao que corresponda según el contexto.

Esto esta muy en linea con el patrón experto.

```go
// NewService es una función que puede mockearse
func NewService() *HelloService {
	return &HelloService{
		dao.NewDao(),
	}
}
```

Donde dao.NewDao() es exactamente esta función que nos devuelve una dependencia, haciendo posible la inversión de control.

```go
// NewDao es el factory
func NewDao() *HelloDao {
	return new(HelloDao)
}
```

Si existe una estrategia de construcción, digamos, singleton, pool de objetos, instancias individuales, o la que sea, esa función se hará cargo.
A su vez, no necesariamente deba existir una sola función, podrían existir varios factories, algo que quedaría bastante bien organizado, y sobre todo bien encapsulado.
También nos permite definir un contexto pasado por parámetros para que nos devuelva la instancia ideal.

> Si bien podríamos pensar que esta estrategia es un antipatrón (Dependency Freak), no lo es, porque la realidad es que HelloDao es interno a nuestro módulo, y no debería necesitar inyección de dependencias para crearlo. En el articulo (Introduction to Dependecy Injection)[https://kariera.future-processing.pl/blog/introduction-to-dependency-injection/], podemos notar aclaraciones puntuales sobre cuando debemos y cuando no debemos usar inyección de dependencias.

Para realizar mocks en los tests solo tenemos que definir un valor para mockedDao

```go
func TestSayHelo(t *testing.T) {
	// Mockeamos
	mockedDao := new(daoMock)

	s := HelloService{
		mockedDao,
	}
	assert.Equal(t, "Hello", s.SayHello())
}
```

Siguiendo los lineamientos de no realizar estrategias donde no es necesario, el dao, no expone interfaces, es solo una estructura, se mockea fácilmente sin necesidad de mayores artefactos.

> Si algo puede ser hacky son los tests, el código productivo debe ser legible y fácil de mantener

Ventajas:

- Permite encapsular el código de forma correcta, definiendo los servicios que se necesitan en el lugar donde se usan.
- No expone los detalles internos de implementación.
- Permite reducir complejidad de constructores.
- Nos evita tener que tener todos los constructores acoplados en un bootstrap.
- Podemos utilizar el patrón experto de forma más clara y concisa.
- Escribimos las estrategias de factories en el archivo adecuado.

## Ahora veamos los fundamentos

En realidad, hacer inyección de dependencias es una practica simpática, el problema es la forma en que se hace, se exponen muchas veces estrategias en los libros y los ejemplos son simples, y funcionan para ese ejemplo, pero no escalan, porque terminan repartiendo responsabilidades incorrectamente y exponiendo inadecuadamente los internals de las implementaciones. (GRASP patterns)

### El patrón Strategy

Y uno de los mayores vicios que vemos en las implementaciones es el abuso de creación de interfaces que no hacen nada.

El patrón estrategia se fortaleció en la programación OO.
Permite establecer diferentes estrategias de resolución de un problema a través de una interfaz y múltiples implementaciones.

Lo cierto es que la existencia de Strategy, es lo que le da sentido a la inyección de dependencias por constructores.

No debemos usar DI por constructor cuando no tenemos strategy. O sea, si realmente existe una interfaz y el usuario de nuestra librería tiene la libertad de implementar el comportamiento, esta perfecto.

Pero si las opciones son limitadas, o bien existe una única opción, es preferible usa Factory Methods.

Porque digo esto ? Porque es muy común observar las siguientes conductas a la hora de programar :

- Implementar interfaces si o si, para separar capas
- Implementar interfaces cuando solo existe una sola implementación
- Utilizar interfaces para poder mockear tests, cuando en realidad existe una sola implementación
- O simplemente porque es la forma que todos dicen

### Lo que realmente deberíamos considerar es que :

- No debemos usar strategy cuando no hay polimorfismo.
- Tampoco cuando las opciones de comportamiento son limitadas, una variable de contexto adecuada en el factory es suficiente para tomar esta decision.
- Una clase mock para testear no es excusa para implementar strategy.
- Si las clases que se están creando se conocen (porque están en el mismo paquete), no es necesario usar DI.
- Solo debemos hacer DI por constructor cuando realmente tenemos una estrategia y la define el cliente de nuestra api.
- Cuando _por las dudas_ generalizamos y hacemos DI, estamos escribiendo código extra innecesario.
- Cuando queremos mockear para unit test, es preferible soluciones hacky.

### Cuales son los problemas de la DI cuando se usa mal:

Aclarando que la inyección de dependencias por Factory Method es una buena practica, y recomendable, los vicios de implementarla por Constructor en cualquier lado cuando no es necesario serian:

- Sobrecargamos los factories y/o métodos con instancias innecesariamente.
- Expone los detalles internos de dependencias.
- Generamos confusión al dejar abierta las puertas al polimorfismo , cuando en realidad no lo hay.
- Acoplamos código. Por ejemplo un controller no debería saber que instancia de un DAO utilizar un Servicio de negocio.
- Hacemos el código difícil de leer y por consiguiente de mantener.

### Cuando SI deberíamos usar DI por constructor

- Cuando tenemos una estrategia, o sea polimorfismo para resolver un problema y el cliente la define (por ejemplo un callback a subrutinas).
- Cuando estamos programando un modulo y la implementación del comportamiento se define fuera del modulo.
- Cuando programamos una librería y queremos ser user friendly para terceros que podrían necesitar algún tipo de implementación hacky.
- Cuando debemos usar dependencias de un servicio que es provisto por otro modulo.
- Cuando accedemos a datos fuera de nuestro modulo, como apis o bases de datos.

### Alternativas creacionales

Cuando tenemos DI por constructor, no necesariamente podríamos usar un Factory Method, existen varios patrones creacionales que podrían ser útiles también, como Builders, Object Pool, etc, lo importante es que esta creación, este asociada al objeto que se crea, y no en cualquier lado y a su vez, instanciada en el componente que la necesita.

### Alternativa funcional

Estamos aprendiendo Go porque queremos ser pragmáticos, la mejor forma de programar go es usando fundamentos funcionales, en los cuales la inyección de dependencias toma un rumbo diferente.

Mas adelante en este tutorial hay mas ejemplos.

## Nota

Esta es una serie de tutoriales sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
