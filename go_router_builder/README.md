[English Version](README_en.md)

# Builder Pattern en Router

## Usar middlewares como Builder Pattern ?

En las notas previas vimos que podríamos usar el router para conformar una respuesta, si pensamos como funciona el patrón builder por ejemplo :

```go
dialog.NewBuilder().Title("hola").AcceptAction("Aceptar", "ok").Build()
```

Esto podemos expresarlo en algo como :

```go
	getRouter().GET(
		"/dialog",
		setTitle,
		setAcceptAction,
		build
	)
```

Este repositorio plantea ejemplos de como usar efectivamente el patrón Builder en un Router de Gin.

## Manejemos información compleja

Podríamos pensar en algún escenario donde la información se genere de formas complejas y nos convenga generar un respuesta en forma segregada utilizando el router como patrón Builder.

Supongamos un microservicio que pretende obtener información del usuario y del perfil, en una sola llamada.

Si bien, no es puntualmente el mismo caso de builder del ejemplo anterior, el uso del router es el mismo:

```go
router().GET(
	"/users/:id",
	validateUserName,
	fetchUser,
	fetchProfile,
	build,
)
```

Como siempre, primero validamos los datos del request, luego buscamos la información necesaria para dar nuestra respuesta : Buscamos el perfil, y el usuario.

Como ultimo paso en la función build armaremos la respuesta final.

Las funciones fetchUser y fetchProfile, son simples, llaman a un servicio del negocio, y ponen los datos en el contexto de gin.

```go
func fetchUser(c *gin.Context) {
	c.Set("user", user.FetchUser())
	c.Next()
}
```

NOTA: Mucho cuidado con el contexto del gin, solo debemos usarlo en el controller, seria incorrecto enviar información a los servicios usando este sistema.

Para armar nuestra respuesta, simplemente buscamos lo que tenemos en el contexto

```go
func build(c *gin.Context) {
	user := c.MustGet("user").(*user.User)
	profile := c.MustGet("profile").(*profile.Profile)

	c.JSON(http.StatusOK, gin.H{
		"login":  user.Login,
		"access": user.Access,
		"name":   profile.Name,
	})
}
```

## Ejecución en paralelo

Las funciones en los services tienen un delay de 1 segundo, para simular una conexión de red. Si vemos el log de respuesta, una llamada tiene una duración de 2 segundos aproximadamente.

Esto lo podemos resolver simplemente ejecutando los ruteos FETCH en forma paralela. En get_parallel_user_id.go vamos a hacer uso de las mismas funciones de ruteo, solo que llamaremos los Fetch en paralelo.

```go
router().GET(
	"/parallel/users/:id",
	validateUserName,
	inParallel(
		fetchUserInParallel,
		fetchProfileInParallel,
	),
	build,
)
```

Debemos crear una función que nos permita la ejecución en paralelo, inParallel esta escrita simple, pero efectiva, el tiempo de respuesta es de 1 segundo, por lo que las llamadas remotas se ejecutan en paralelo.

```go
func inParallel(handlers ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(handlers))

		for _, handler := range handlers {
			go func(handlerFunc gin.HandlerFunc) {
				defer waitGroup.Done()
				handlerFunc(c)
			}(handler)
		}

		waitGroup.Wait()

		c.Next()
	}
}
```

Es genérica, por lo que debería estar en algún paquete de utilidades reutilizables.

En las funciones fetch, lo único a tener en cuenta es que ahora los handlers no llamaran a Next, porque lo hace inParallel.

```go
func fetchUserInParallel(c *gin.Context) {
	c.Set("user", user.FetchUser())
}
```

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index/blob/main/README.md)
