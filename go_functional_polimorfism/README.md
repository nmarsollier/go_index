<!-- cSpell:language es -->

[English version](README_en.md))

# Polimorfismo con Funciones

Uno de los aspectos mas decisivos a la hora de crear interfaces en Go es el uso de polimorfismo.

## Una estrategia sencilla

Digamos que tenemos un servicio rest que retorna datos del perfil de un usuario, pero que dependiendo de quien lo llame (mobile o web), la imagen tiene una estrategia diferente de carga de acuerdo al origen.

Vamos una estrategia que me encanta como funciona :

```go
var GetImage = map[string]func(id string) string{
	"mobile": func(id string) string {
		return fmt.Sprintf("resource_%s", id)
	},
	"web": func(id string) string {
		return fmt.Sprintf("http://www.images.blablabla/images/%s", id)
	},
}
```

Si, es un simple mapa con funciones para cada caso, en lenguajes como Kotlin podríamos definir un Enum, pero para los fines de Go un map funciona perfecto.

Es un comportamiento muy simple, una función, pero si hacemos volar un poco la imaginación, podemos darnos cuenta que podría ser una estructura lo que retornemos, com punteros a funciones sobrecargadas según sea el caso, y funcionaria muy bien.

Este código se usa como:

```go
profile.GetImage[device](id)
```

Sencillo y elegante.

Incluso nos proporciona una función de validación de parámetros que es muy directa, y nos permite agregar nuevos devices sin mayores costos :

```go
// IsValidDevice checks if device is valid
func IsValidDevice(device string) bool {
	_, ok := GetImage[device]
	return ok
}
```

### Desventajas

- Bueno, la llamada a la función no es algo típico de ver
- Funciona solo para una función, pero si tenemos en cuenta el principio de single responsability, puedo vivir con ello.
- Si la sobrecarga del modulo es muy grande, bueno ya no conviene mucho esta estrategia.
- Mockear los tests puede ser un poquito mas hacky.

## Pero veamos que nos ahorramos, por si no lo notamos

### Si lo hubiéramos hecho con funciones al estilo clásico funcional

```go
func GetImage(device, id string) string{
	switch device {
	case "mobile":
		return fmt.Sprintf("resource_%s", id)
	case "web":
		return fmt.Sprintf("http://www.images.blablabla/images/%s", id)
	}
}
```

No parece gran cosa, pero la complejidad ciclomática de esa función es mayor a la del ejemplo anterior.

Agregar nuevos valores es un poco mas complejo.

Y Deberíamos programar una función para validar que no es tan directa.

### Si lo hubiéramos hecho con un estilo orientado a objetos

Esta sin lugar a dudas es la opción que elegiríamos como primera elección, cuando nos enfrentamos a polimorfismos.

Tendríamos que programar :

- Una interfaz que exponga el método
- 2 estructuras que implementen la interfaz
- Una función factory que dependiendo del device devuelva la implementación correcta, que tendría un switch
- Una función de validación igual de fea que en el caso anterior

Ya no me dan ganas ni de programar el ejemplo, no quiero ni contar las lineas de código.

Incluso un enfoque funcional nos permite definir sobrecargas de comportamientos para funciones puntuales, cuando muchas veces nos quejamos que Go no maneja Herencias, con la estrategia original, solo sobrecargamos comportamiento donde debe ser.

Lo mas probable que la próxima vez que enfrentes un caso de polimorfismo te preguntes, ¿ porque tengo que programar todo esto ?

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
