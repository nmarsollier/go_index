<!-- cSpell:language es -->

[English Version](README_en.md)

# Un poco mas declarativos

Si bien a la fecha, Go es un lenguaje netamente imperativo, en este repositorio planteo algunos tips, para escribir código un poco mas declarativo.

## Que seria ser declarativo ?

Go no es declarativo, por consiguiente, es una excelente pregunta.

Ser declarativos es evitar escribir orientándonos a los procedimientos lógicos que nos dan un resultado, ejemplo :

```go
// Paradigma imperativo
func Suma(start, end int) int {
	suma := 0
	for i := start; i < end; i++ {
		suma += i
	}
	return suma
}

// Paradigma declarativo
suma := Suma(1, 100)
```

La idea, entonces es, generar una seria de funciones genéricas (imperativas) que nos abstraigan un poco de la parte procedural, para poder escribir nuestra lógica de negocios en forma declarativa, de forma tal que resulte mucho mas fácil de leer.

Los lenguajes funcionales en general son declarativos, y no imperativos.

Muchos lenguajes OO modernos también nos dan muchas herramientas para expresar el código en forma declarativa, como Typescript, Kotlin o Swift.

### La programación procedural desaparece ?

En go no, porque no es declarativo, deberemos programar funciones procedurales en algún lugar, pero la idea es que sean escritas a modo de tools, o librerías, que sean simples, que respeten el paradigma funcional, puntuales y que se usen a alto nivel.

## Y cual es la ventaja ?

Existen grandes ventajas, a la hora de definir lógica de negocios, la programación declarativa nos permite ser mas expresivos, definir lo que queremos obtener en un lenguaje mas humano y por consiguiente mas fácil de leer y de mantener.

## Ejemplo: Una función para no repetir

Una de las ideas principales, es intentar generar funciones para esas cosas que se repiten en todos lados, por ejemplo, del código de los tutoriales anteriores :

```go
func sayHelloHandler(c *gin.Context) {
	userName := c.Param("userName")

	c.JSON(http.StatusOK, gin.H{
		"answer": service.SayHello(userName),
	})
}
```

```go
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"answer": "pong",
	})
}
```

Si escribimos una librería como :

```go
func SendJSONAnswer(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"answer": data,
	})
}
```

Nos permite ser mas declarativo, a su vez no repetir. Vemos como queda el código nuevo:

```go
func sayHelloHandler(c *gin.Context) {
	gu.SendJSONAnswer(c, service.SayHello(c.Param("userName")))
}
```

```go
func pingHandler(c *gin.Context) {
	gu.SendJSONAnswer(c, "pong")
}
```

## El nombre de las funciones es la clave

Un aspecto interesante que buscamos cuando programamos en forma declarativa, es intentar que el código se parezca lo mas posible a una conversación en lenguaje natural.

Si bien esta forma de pensar no esta 100% alineada con la practica de escribir lo menos posible de Go, tampoco debemos considerar que Go nos sugiere escribir de menos a la hora de ser expresivos, hay que buscar el balance justo.

Por ejemplo, siguiendo el ejemplo anterior, si en vez de llamarle sayHelloHandler al método le llamamos sayHello, queda mas natural:

```go
func init() {
	getRouter().GET(
		"/hello/:userName",
		validateUserName,
		sayHello,
	)
}
```

Y somos mas declarativos.

## Un builder

Un builder es el caso típico de uso de programación declarativa

```go
dialog.NewBuilder().Title("Hola Mundo").AcceptAction("Aceptar", "ok").Build()
```

Pienso que lo tiene de interesante el patrón builder, es que nos gusta porque nos permite ser declarativos.

Una implementación de builder aceptable, podría ser :

```go
package dialog

import (
	"encoding/json"
)

type dialogAction struct {
	Label  string `json:"label,omitempty"`
	Action string `json:"action,omitempty"`
}

type dialog struct {
	Title  *string       `json:"title,omitempty"`
	Accept *dialogAction `json:"accept,omitempty"`
}

type DialogBuilder struct {
	dialog
}

func NewBuilder() *DialogBuilder {
	return &DialogBuilder{dialog{}}
}

func (d *DialogBuilder) Title(value string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title:  &value,
		Accept: d.dialog.Accept,
	}}
}

func (d *DialogBuilder) AcceptAction(label, action string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title: d.dialog.Title,
		Accept: &dialogAction{
			Label:  label,
			Action: action,
		},
	}}
}

func (d *DialogBuilder) Build() string {
	result, _ := json.Marshal(d.dialog)
	return string(result)
}
```

Podemos evitarnos la creación del builder, simplemente con una función, pero este tipo de patrones es mas declarativo y mas simple de leer.

## Encadenando Comportamientos

Podemos escribir el resto de las apps en forma mas declarativa ?
No se si justifica para un lenguaje como Go, pero seguramente alguna vez nos vamos a cruzar con código que es mas sencillo desarrollarlo con la siguiente estrategia...

Veamos un ejemplo de una función en forma imperativa

```go
func Shorten(name string) string {
	values := strings.Split(name, " ")

	result := ""

	for _, v := range values {
		if len(v) > 0 {
			result += strings.ToUpper(string(v[0]))
		}
	}

	return result
}
```

Ya se dieron cuenta que hace ? Que pasa si lo escribimos en forma declarativa :

```go
func Shorten(name string) string {
	return fromString(name).
		split(" ").
		mapNotEmpty(func(s string) string {
			return strings.ToUpper(string(s[0]))
		}).
		joinToString()
}
```

Podemos notar luego de leer el código declarativo, que es mucho mas simple entender que hace, porque tal cual se lee...

- tomamos un string
- separamos las palabras en un arreglo segun los espacios en blanco
- y concertimos en ese array en otro segun el criterio: a cada elemento le tomamos el primer char y lo convertimos a mayusculas
- luego volvemos a concatenar el array en un string

simplemente convierte algo como "uno dos tres" en "UDT".

Ahora para que el código anterior funcione, necesitamos estas herramientas :

```go
type shorten struct{ string }
type shortenSlice struct{ slice []string }

func fromString(value string) shorten {
	return shorten{value}
}

func (s shorten) split(separator string) shortenSlice {
	return shortenSlice{strings.Split(s.string, separator)}
}

func (values shortenSlice) mapNotEmpty(f func(string) string) shortenSlice {
	var result []string
	for _, v := range values.slice {
		if len(v) > 0 {
			result = append(result, f(v))
		}
	}
	return shortenSlice{result}
}

func (values shortenSlice) joinToString() (result string) {
	for _, v := range values.slice {
		result += v
	}
	return result
}
```

La cantidad de código a escribir parece mucha, pero consideremos que shorten y shortenSlice deberían ser librerías genéricas, y reutilizables, nuestro enfoque a la hora de leer que hace se limita a la función Shorten.

Siempre es conveniente comenzar con una buena base, por ejemplo la librería https://github.com/jucardi/go-streams, nos proporciona formas declarativas básicas para manejar listas y luego le agregaremos nuestras funciones mas concretas.

Creo que es importante ir revisando estas estrategias, en el futuro de Go nos han prometido Geneircs, y un manejo mas declarativo en el lenguaje, veremos.

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
