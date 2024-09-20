[English version](https://github.com/nmarsollier/go_di_ioc/blob/main/README_en.md)

# DI e IoC en GO

Este repositorio plantea alternativas de manejo de dependencias, a la programación tradicional de un proyecto Go. 

## Mal uso de Inyección de Dependencias

Es esa estrategia de IoC que nos permite proveer dependencias en una estructura de datos para que sean usadas internamente.

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

Y esa es una buena practica, pero cuando las dependencias son realmente implementaciones que se ajusten a casos donde si se debe usar inyección de dependencias (ejemplo: Strategy pattern), aclarando que esta estrategia no está mal, el problema es cuando abusamos de ella en casos donde no es correcto usarse, en nuestros sistemas, la gran mayoria de dependencias no necesitan usar DI.

Un uso como el anterior no desacopla realmente, todo lo contrario, terminamos acoplando mucho mas, nuestro código debe definir métodos bootstraps en lugares donde no deberían estar, acoplando todo el negocio en un archivo main.go por ejemplo. 

Ademas de acoplar, exponemos los detalles internos de nuestras implementaciones, algo que deberia estar encapsulado.

Podemos revisar las raices de programacion orientada a objetos y darnos cuenta que incluso iría en contra del principio de [Information Expert](https://en.wikipedia.org/wiki/GRASP_(object-oriented_design)#Information_expert), que nos dice que las dependencias deben ser instanciadas en el lugar donde se tiene la informacion, de modo tal que por ejemplo crear las dependencias de un servicio en la capa de controladores, genera acoplamiento y nuestro código ya no seria cohesivo.

La realidad, es que se usa porque nos da la posibilidad de mockear esas dependencias para realizar buenos testings. Sin embargo existen muchas alternativas que podemos usar sin llegar a usar DI.

## Uso de Factory Methods como IoC

Veamos como podemos mejorar la situación anterior.

Inversión de control básicamente significa tener un framework, que cuando necesite un recurso se lo pida a él, y el recurso se obtiene del contexto.

Un service locator, es básicamente un framework que conoce nuestras dependencias, y las inyecta en donde sea necesario, pero tiene el mismo problema que el bootstrap anterior, acopla todos los servicios en un solo lugar.

Si partimos los patrones generales de asignación de responsabilidades GRASP, una de las formas clásicas y adecuadas de construir objetos es el uso de Factory Methods.

Pensemos en ese factory method, como parte de un framework de inyección de dependencias, que dependiendo del contexto (parametros en el factory u ortro metodo) nos va a retornar la instancia correcta del servicio que necesitemos. Lo que tiene de adecuado este patrón, es que la estrategia de creación, se escribe junto a los servicios, por lo que queda mucho mas claro el funcionamiento del mismo.

Esta estrategia nos permite evitar inyectar las dependencias en los constructores y delegar la instanciación a funciones factory que estan en el lugar adecuado.

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
Tambien nos permite definir un contexto pasado por parametros para que nos devuelva la instancia ideal.

> Si bien podriamos pensar que esta estrategia es un antipatron (Dependency Freak), no lo es, porque la realidad es que HelloDao no es algo externo a nuestra aplicacion, y no deberia usarse inyeccion de dependencias para instanciarlo. En el articulo (Introduction to Dependecy Injection)[https://kariera.future-processing.pl/blog/introduction-to-dependency-injection/], podemos notar aclaraciones puntuales sobre cuando debemos y cuando no debemos usar inyección de dependencias.

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

> Si algo puede ser hacky son los tests, el codigo productivo debe ser legible y facil de mantener

Ventajas:
- Permite encapsular el código de forma correcta, definiendo los servicios que se necesitan en el lugar donde se usan.
- No expone los detalles internos de implementacion.
- Permite reducir complejidad de constructores.
- Nos evita tener que tener todos los constructores acoplados en un bootstrap.
- Podemos utilizar el patrón experto de forma más clara y concisa.
- Escribimos las estrategias de factories en el archivo adecuado.

## Ahora veamos los fundamentos 

En realidad, hacer inyección de dependencias es una practica simpatica, el problema es la forma en que se hace, se exponen muchas veces estrategias en los libros y los ejemplos son simples, y funcionan para ese ejemplo, pero no escalan, porque terminan repartiendo responsabilidades incorrectamente y exponiendo inadecuadamente los internals de las implementaciones. (GRASP patterns)

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
- Si las clases que se estan instanciando se conocen (porque estan en el mismo paquete), no es necesario usar DI.
- Solo debemos hacer DI por constructor cuando realmente tenemos una estrategia y la define el cliente de nuestra api.
- Cuando *por las dudas* generalizamos y hacemos DI, estamos escribiendo código extra innecesario.
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

### Alternativas creacionales

Cuando tenemos DI por constructor, no necesariamente podríamos usar un Factory Method, existen varios patrones creacionales que podrían ser útiles también, como Builders, Object Pool, etc, lo importante es que esta creación, este asociada al objeto que se crea, y no en cualquier lado y a su vez, instanciada en el componente que la necesita.

## Fuentes

[Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection)

[Introduction to Dependency Injection](https://kariera.future-processing.pl/blog/introduction-to-dependency-injection/)

[GRASP](https://es.wikipedia.org/wiki/GRASP)

[Service locator pattern](https://en.wikipedia.org/wiki/Service_locator_pattern)

[Strategy (patrón de diseño)](https://es.wikipedia.org/wiki/Strategy_(patr%C3%B3n_de_dise%C3%B1o))

[Patrón de diseño](https://es.wikipedia.org/wiki/Patr%C3%B3n_de_dise%C3%B1o)

[YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)

[KISS](https://en.wikipedia.org/wiki/KISS_principle)

[InterfaceImplementationPair](https://martinfowler.com/bliki/InterfaceImplementationPair.html)

[Foo/FooImpl pairs - stop doing it!](http://wrschneider.github.io/2015/07/27/foo-fooimpl-pairs.html)

[When to Mock](https://blog.cleancoder.com/uncle-bob/2014/05/10/WhenToMock.html)

[A theory of modern Golang](https://news.ycombinator.com/item?id=14521894)

## Nota

Esta es una serie de tutoriales sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index/blob/main/README.md)
