<!-- cSpell:language es -->

[English Version](README_en.md)

# REST Controllers en go

Este repositorio plantea una forma simple y efectiva de organizar nuestros servicios REST en un entorno de microservicios.

## RMR - Resource-Method-Representation

RMR es una variante de MVC.

Con RMR nuestros controllers se organizan en orden a las solicitudes Http
Cada archivo contiene una sola definición de una entrada a un controller
Estructuramos los directorios del controller en base a eso
Los nombres del archivo hacen referencia a la entrada http

Por consiguiente tenemos la estructura :

```go
 - controllers
       get_hello_username.go
       get_ping.go
       rest.go
```

Donde rest.go contiene la inicialización del framework, pero no las rutas.

```go
func Start() {
	getRouter().Run(":8080")
}

var router *gin.Engine = nil

func getRouter() *gin.Engine {
	if router == nil {
		router = gin.Default()
	}

	return router
}
```

Tenemos un archivo por método Rest (get_resource, put_resource, get_resources_id, etc) Cada uno maneja una sola función Rest.

En estos archivos definimos todo lo que tenga que ver con esa ruta desde la definición de la misma hasta la implementación del controller.

```go
// Internal configure ping/pong service
func init() {
	getRouter().GET("/ping", pingHandler)
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"answer": "pong",
	})
}
```

Podemos notar que toda la implementación de la ruta queda contenida en un solo archivo.

Este mismo concepto se puede adaptar y utilizar con cualquier protocolo como GRPC o sistemas de Mensajerías.

### Ventajas

- Simplifica el código desde su estructura, separa conceptos claramente
- Orienta nuestras apps y negocio específicamente a hacer algo puntual (una responsabilidad), desde el controller
- Orienta el código a responsabilidades simples
- Simplifica la lectura y búsqueda del controller
- Simplifica el testing
- Encapsula correctamente cada controller
- Desacopla las inicializacion de rutas, haciendo la definición de la misma algo sustentable y mantenible
- Permite una lectura clara de los middlewares de cada controller

## Fundamentos

### MVC en Microservicios

Porque microservicios es algo diferente a un monolito ?

- Un microservicio en general maneja un solo aspecto puntual de todo el modelo de negocio
- Posee una interfaz mucho mas simple
- En general, la arquitectura de microservicios define como se deben comunicar los microservicios, y muchas veces se une el concepto de View y Controller, ya que nos enfocamos muy pocos protocolos de E/S.
- Si existen diferentes View, en general se manejan con diferentes microservicios Api Gateways.
- No hay tanta segregación a implementar.

### El Controller

El en enfoque clásico MVC, el Controller y View poseen las siguientes funciones:

- Interpreta un Request
- Valida los datos de entrada
- Adaptar el request a una solicitud del negocio
- Llama al negocio
- Adapta la respuesta del negocio de acuerdo al cliente
- Maneja errores

Y en general en los microservicios, conviene usar un framework especifico que resuelva cómodamente estos aspectos.

### Notas finales sobre la definición del protocolo REST

No voy a entrar en detalles de como debería ser un protocolo REST, sino mas bien una simple introducción. En general podemos organizar nuestros repositorios de dos formas :

#### Resource Centric

Es donde el centro de la información es un recurso en particular. Hay mucha información, se le llama [RestFul](https://en.wikipedia.org/wiki/Representational_state_transfer).

Hay mucha info de ésto.

#### Use Case Centric

El formato RestFul no es util cuando tenemos múltiples casos de uso sobre un recurso, por ejemplo una factura, podemos necesitar hacer diversas consultas y diversas acciones sobre la misma. Por lo que es conveniente definir el protocolo en base a casos de uso.

De esto si que no hay un estándar. Existen varias formas de hacerlo, una forma simple y efectiva, es :

- Utilizamos GET para consultar (queries)
- Utilizamos POST para cambiar el estados (commands)
- Agregamos un /sustantivo al final de los métodos GET para saber que recurso queremos
- Agregamos un /verbo al final de un POST para saber que caso de uso queremos ejecutar

Ejemplo :

GET http://go_rest_controller.com/facturas/resumen

GET http://go_rest_controller.com/facturas/total_consolidado

GET http://go_rest_controller.com/facturas/:id/recibo

GET http://go_rest_controller.com/facturas/:id/detalles

GET http://go_rest_controller.com/facturas/:id/total

POST http://go_rest_controller.com/facturas/:id/enviar_correo

POST http://go_rest_controller.com/facturas/:id/pagar

POST http://go_rest_controller.com/facturas/:id/cancelar

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
