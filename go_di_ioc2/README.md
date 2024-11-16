<!-- cSpell:language es -->

[English version](README_en.md)

# Luego de funcional, volvemos a hablar de inyección de dependencias

Ya vimos en artículos anteriores, como aplicar programación funcional, porque es mejor usar ese estilo de programación en go, ahora revisamos inyección de dependencias con lo que ya sabemos.

## Un patrón que no existe

En programación funcional, la inyección de dependencias, per se, no existe, las estructuras no son creadas nunca con la intención acceder a dependencias.

> Las interfaces en go no son para inyectar dependencias sino que permiten implementar el [patrón strategy](../go_di_ioc/README.md), sin embargo la estrategia a utilizar no debe ser pasada por parámetros, sino mas bien cada función debe recibir los parámetros necesarios para encontrar las dependencias requeridas.

Cuando escribimos una función en el estilo funcional, básicamente respetamos lo siguiente :

- Las funciones deben realizar solo una cosa (es la clave)
- Nombre claro y conciso que se autoexplique
- En go, nombre corto y sencillo de recordar
- 2 o 3 argumentos máximo, siempre que sean claros y concisos
- Las funciones solo deben recibir los datos que necesitan, ni mas ni menos.
- Si tenemos muchos argumentos se pasa una estructura, simplifica refactor, mantenimiento y significado de los parámetros.
- Las funciones deben tener un solo nivel de abstracción.
- Las estructuras que se pasan por parámetros y que se retornan, deben ser inmutables
- En general las funciones que corresponden a la misma capa de negocio y están relacionadas, deben escribirse cerca en el código.

## Ahora bien, como pasamos dependencias ?

El paso de dependencias no es un problema, cada función recibe los parámetros que sean necesarios para su correcto funcionamiento, cada función tiene acceso directo a las funciones de las cuales depende, por lo que no es necesario que enviemos punteros en sus parámetros, por ejemplo un service y un dao, el service debe tener la lógica que determine el dao correspondiente a llamar, esto se le llama patrón experto de información.

Mas adelante en esta guía se explican mas estrategias.

## Que datos debe recibir una función ?

Solo los que necesita, nunca debemos pasar datos que la función no necesita, o estructuras complejas que luego no se usan, siempre es preferible recibir los parámetros justos y cuando son muchos o confusos, definir y pasar por estructuras esos parámetros, de forma tal que quede claro que una función necesita eso y solo eso, ni mas, ni menos.

Las funciones deben verse como cajas cerradas desde el mundo exterior, son cajas que necesitan cierta información, porque responden a cierta funcionalidad bien clara, y esa data que necesitan es suficiente para que el que usa la función entienda que se necesita y a veces con solo conocer los datos ya sabemos el porque.

> Pasar siempre una estructura como parámetro de una función no es buena practica, solo definimos una estructura como parámetro cuando los parámetros son confusos de leer, caso contrario conviene pasar parámetros individuales.

Un error muy común en servicios http es pasar el contexto y que las funciones extraigan valores del contexto, las funciones deben recibir el contexto solo para cancelar goroutines, por ejemplo, pero nunca para sacar valores del mismo.

El contexto go es una bolsa de información que nunca nos deja claro que requisitos debe tener para ser valido, si bien podemos usar el contexto y debemos usar el contexto para poner valores, estos valores quedan restringidas en su uso dentro de los controladores, cuando llamamos a un service extraemos esos valores y se debe llamar con el valor que el service necesita puntualmente.

## Ejemplo

En este proyecto adopte una estrategia que permite pasar un parámetro variable a las funciones que representa un "contexto", pero no un contexto go, sino mas bien un contexto funcional de servicios de negocio que se deben usar, este contexto generalmente esta vacío, salvo que queramos proporcionar diferentes implementaciones a librerías (Ej: Cuando hacemos unit test, o bien cuando queremos pasar alguna instancia particular para diferentes casos de uso)

Ejemplos puntuales:

- Manejo de la conexión a la base de datos: Podemos mockear una conexión o bien el mismo constructor db.Get instancia la base de datos real.
- El manejo de logs, los logs dependen de un contexto particular, ya sea un request gin o un proceso de mensajes Rabbit. En particular el correlation_id tiene diferentes valores según sea necesario y las funciones que deban realizar logs deben tener ese contexto particular mientras duren las operaciones relacionadas.

Conceptos clave de este enfoque:

- Las funciones son las que se encargan de crear los servicios necesarios (no les pasamos los servicios a usar por parámetro).
- Las funciones están desacopladas de la forma en la que se crean los servicios.
- Los servicios tienen un constructor que recibe el contexto (var arg) y en base al contexto determina la instancia a usar.
- Las funciones que necesiten usar un servicio usan la función del punto anterior para acceder a esas funciones.

Ahora bien, cada Servicio que puede tener mas de una implementación es el encargado de

[imagego](https://github.com/nmarsollier/imagego).

**Construcción de dependencias**

El contexto se define como un parámetro variable, que cuando llamamos a los constructores se recibe como parámetro, cada componente que necesitemos "inyectar" como dependencia funcional, debe tener un constructor que recibe el contexto funcional, si el contexto ya provee una dependencia se usa esa dependencia, caso contrario el constructor retorna la adecuada.

Esta es una forma muy elegante de desacoplar las instancias y delegar la creación de las instancias correctamente a los componentes que conocen como crearlos.

```go
func Get(ctx ...interface{}) RedisClient {
  // Si el contexto proporciona una instancia usamos esa instancia sino retornamos la instancia de producción
	for _, o := range ctx {
		if client, ok := o.(RedisClient); ok {
			return client
		}
	}

	once.Do(func() {
		instance = redis.NewClient(&redis.Options{
			Addr:     env.Get().RedisURL,
			Password: "",
			DB:       0,
		})
	})
	return instance
}

```

**Inicialización**

El contexto se inicializa cuando comienzan las operaciones en un controller adecuadamente y se pasa a todas las funciones que sea necesario.

En este caso la función se define en el contexto de un servidor gin como :

```go

// Gets the context for external services
func GinCtx(c *gin.Context) []interface{} {
	var ctx []interface{}
	ctx = append(ctx, ginLogger(c))
	return ctx
}
```

En este caso puntual el contexto se inicializa con un logger que es usado para hacer un seguimiento del correlation_id. Esta instancia de logger analiza el request en busca de un header que requiera una trazabilidad de correlation_id si lo encuentra usa ese sino crea un logger nuevo con un correlation_id nuevo.

Todas la llamadas siguientes ya tendrán una instancia de logger a usar.

```go
func initPostImage() {
	server.Router().POST(
		"/v1/image",
		server.ValidateAuthentication,
		saveImage,
	)
}

func saveImage(c *gin.Context) {
	bodyImage, err := getBodyImage(c)
	if err != nil {
		c.Error(err)
		return
	}

  // Obtenemos el contexto y luego se lo pasamos a las funciones que lo necesiten
	ctx := server.GinCtx(c)
	id, err := image.Insert(image.New(bodyImage), ctx...)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, NewImageResponse{ID: id})
}
```

Como podemos ver, de forma correcta image.Insert no necesita saber que logger se usa ni que instancia de Redis se usa, simplemente se le provee un contexto, en el caso del logger inicializa gin, pero en el caso de Redis es el mismo factory de redis el que lo inicializa a demanda.

### Test

Con una pequeña modificación a GinCtx, podriamos mockear en GinCtx un contexto completo para tests

```go

// Gets the context for external services
func GinCtx(c *gin.Context) []interface{} {
	var ctx []interface{}
	// mock_ctx solo es para mocks en testing
	if mocks, ok := c.Get("mock_ctx"); ok {
		return mocks.([]interface{})
	}

	ctx = append(ctx, ginLogger(c))

	return ctx
}
```

Luego simplemente antes de llamar al test, debemos poner en el contexto gin los mocks necesarios.

```go
  // Inicialización previas de Gin

	// Mocks Redis
	ctrl := gomock.NewController(t)
	redisMock := redisx.NewMockRedisClient(ctrl)
	redisMock.EXPECT().Get(gomock.Any()).Return(redis.NewStringResult("", errs.NotFound)).Times(1)

  mock_ctx := []interface{}{
    redisMock,
    log.NewTestLogger()
  }

  engine.Use(func(c *gin.Context) {
    c.Set("mock_ctx", mock_ctx)
    c.Next()
  })

  // Llamamos al request que corresponda en gin.
```

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
