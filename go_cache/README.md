[English version](README_en.md)

# Una forma adecuada de hacer cache

Cachear valores es algo sencillo de hacer, pero si no lo hacemos bien, puede ser un problema.

## Una estructura base de cache

Vamos a definir algo de código con algo de lógica para almacenar nuestro cache.

Tenemos un constructor que nos va a permitir cachear value, por el tiempo definido en retain.

```go
func Memoize(value interface{}, retain time.Duration) *Memo

```

Memo es una estructura que va a contener los datos de cache necesarios. Podríamos combinar con una interfaz en caso que queramos hacer algo diferente.

Si queremos saber el valor actual usamos

```go
func (m *Memo) Value() interface{}
```

Value nos devuelve nil si el valor ya no es valido, o bien un valor puntual

Sin embargo vamos a poder seguir conociendo el valor de cache, sin importar si es valido o no :

```go
func (m *Memo) Cached() interface{}
```

## Un uso muy simple

La estructura Memo es de solo lectura.

Supongamos que queremos hacer un cache global, de un proceso costoso, digamos una llamada a la red muy cara.

Un ejemplo de uso muy sencillo, **aunque incorrecto**, podría ser :

```go
var cache *memoize.Memo = nil

func wrongCache1(id string) *Profile {
	if cache == nil || cache.Value() == nil {
		value := fetchProfile("123")
		cache = memoize.Memoize(value, 10*time.Minute)
	}

	return cache.Value().(*Profile)
}
```

Primero comprobamos que el valor de cache es valido usando Value(), si no es valido, buscamos el valor adecuado y lo cacheamos.

Luego retornamos el valor con Cached() dado que Value(), para prevenir race conditions, donde Value() no ser valido en ese momento.

## Veamos los problemas

El codigo anterior puede ser muy util para una aplicacion que no posee mutiples threads, pero si ejecutamos el codigo anterior en un entorno concurrente, tendremos muchos problemas...

### Race condition en el if

Si analizamos el siguiente if, podremos ver que tenemos un primer problema :

```go
	if cache == nil || cache.Value() == nil {
```

en un entorno de multitreads la expresión del if se evaluá en 2 pasos, si bien el primer paso es validar cache == nil la segunda validación podria fallar, y es justamente porque cache podría ser nil, dado que entre la primera expresión y la segunda, cache podría convertirse en nil.

Por ejemplo, si quisiéramos invalidar el cache, lo mas lógico seria hacer :

```go
func invalidateCache() {
	cache = nil
}
```

y cache es una variable compartida (un singleton), porque asi debe funcionar este cache.

**Solución**

Hacemos una copia del valor original y trabajamos sobre esa copia, de esta forma, es thread safe.

```go
func wrongCache2(id string) *Profile {
	currCache := cache

	if currCache == nil || currCache.Value() == nil {
		value := fetchProfile("123")
		currCache = memoize.Memoize(value, 10*time.Minute)
		cache = currCache
	}

	return cache.Value().(*Profile)
}
```

### Race condition en el return

Nótese que también, en el ejemplo, que return es incorrecto, porque Value podría volverse invalido entre la primera lectura y el return.

**Solución**

Lo correcto es retornar, ya que Cached no tiene en cuenta la validez del cache, y nunca es nil.

```go
	return currCache.Cached().(*Profile)
```

### Race condition en la carga de datos

Supongamos ahora que mas de un proceso llama en forma concurrent a esta función, si fetchProfile fuera una operación de milisegundos, no seria un problema, pero ya que es una operación que puede llevar segundos en terminar, estos procesos podrían tomar el valor de cache invalido, todos juntos, y todos ellos llamarían a fetchProfile en forma concurrente, lo que podría provocar no solo llamadas múltiples, sino colapsar el servicio remoto, haciendo que las respuestas sean cada vez mas caras de evaluar [Ver](https://en.wikipedia.org/wiki/Cache_stampede).

```go
func WrongCache2(id string) *Profile {
	currCache := cache

	if currCache == nil || currCache.Value() == nil {
			value := fetchProfile("123")
```

Como podemos ver en los tests TestNonConcurrentWrongCache2 nos da como resultado

```
=== RUN   TestNonConcurrentWrongCache2
Fetching Profile... 0
```

Pero el caso concurrente no nos da el resultado esperado :

```
=== RUN   TestConcurrentWrongCache2
Fetching Profile... 1
Fetching Profile... 5
Fetching Profile... 2
Fetching Profile... 3
Fetching Profile... 4
Fetching Profile... 6
Fetching Profile... 9
Fetching Profile... 0
Fetching Profile... 8
```

**Solución**

Debemos bloquear el proceso para que solo un proceso acceda a esta rutina de loading.

Para hacer esto en forma ordenada, necesitamos cambiar un poco la estrategia, vamos a usar early exits :

```go
// Mutex nos va a permitir lockear procesos concurrentes
var mutex = &sync.Mutex{}

func wrongCache3(id string) *Profile {
	currCache := cache
	// Si el valor del cache actual es valido, retornamos ese valor
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

  // Si el valor del cache era invalido, hacemos lock,
	// todas las funciones concurrentes deben esperar a unlock.
	defer(mutex.Unlock())
	mutex.Lock()

  // Si múltiples procesos quedaron bloqueados por mutex, cuando se liberen
	// lo mas probable es que cache ya tenga un valor valido, en cuyo caso, nos fijamos
	// para evitar la llamada a la red
	currCache = cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	// Llamamos a la red y actualizamos el cache
	value := fetchProfile(id)
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}
```

Ahora el resultado es el esperado, se llama una sola vez, en procesos concurrentes.

```
=== RUN   TestConcurrentWrongCache3
Fetching Profile... 2
--- PASS: TestConcurrentWrongCache3 (0.00s)
```

## La solución

La función anterior sigue teniendo el nombre wrongCache3, porque si bien es correcta, se esta desperdiciando una oportunidad importante, si nosotros tenemos un cache expirado, pero valido, el proceso bloquea las llamadas, haciendo que múltiples llamadas queden bloqueadas, cuando es realidad tenemos un cache valido que podríamos usar, mientras se carga el dato actualizado.

Que es lo que se espera de este tipo de cache :

- Que si no tiene valor, solo la primera llamada haga un fetch remoto, las demas esperen ese resultado bloqueadas
- Que se haga una sola llamada
- Que si el cache esta obsoleto, pero es valido como respuesta, se devuelva el cache y el paralelo se obtengan nuevos valores, sin bloquear llamadas

Esta estrategia no solo es mas rapida, sino que es resilente con errores de carga de datos remoto.

**Solución**

Explico la lógica de la solución :

```go
var cache *memoize.Memo = nil
var mutex = &sync.Mutex{}
var loading int32 = 0

func fineFetchProfile(id string) *Profile {
	// Copiamos el cache en variable local, para evitar race conditions del if
	currCache := cache
	if currCache != nil && currCache.Value() != nil {
		// El cache es valido, por lo tanto retornamos el valor
		// Retornamos Cached() que nos evita un race condition de expiración con Value()
		return currCache.Cached().(*Profile)
	}

	// loading es un semaforo, si tiene valor 1 es porque se estan cargando datos
	// actualmente, si tiene valor 0 es porque no hay cargas en proceso
	loadData := atomic.CompareAndSwapInt32(&loading, 0, 1)

	// CompareAndSwapInt32 devuelve en forma atomica true, si pudo
	// asignar loading de 0 a 1 , conceptualmente comprobamos si no hay otro
	// proceso haciendo loading de datos
	if !loadData && currCache != nil {
		// si llegamos aca, es porque hay un proceso cargando datos en el cache
		// y el currCache es no nulo lo que significa que hay un cache que si bien
		// expiró todavía es un dato valido
		return currCache.Cached().(*Profile)
	}

	// Hasta ahora lo anterior no bloqueaba procesos concurrentes
	// a partir de aca solo un proceso puede estar en ejecución,
	// los demás esperan bloqueados a partir de lock()
	defer mutex.Unlock()
	mutex.Lock()

	// Luego de Unlock los procesos que estuvieron esperando, ya tienen
	// un cache valido, y no deben proseguir, simplemente usar el
	// cache actual, hay que validar esto.
	// El primer proceso "el que carga" no retorna.
	currCache = cache
	if loading == 0 && currCache != nil && currCache.Value() != nil {
		// Entra aca si es un proceso que estuvo esperando una carga
		// y luego del unlock, hay un valor valido para retornar
		return currCache.Cached().(*Profile)
	}

	// steamos que no estamos haciendo mas loading
	// esto no puede ser antes, porque hay returns en el medio
	// que harían que el defer seteara este valor a 0, y seria incorrecto
	defer func() { loading = 0 }()

	// Hacemos nuestro proceso costoso
	value := fetchProfile(id)

	// acutalizamos cache
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}
```

**Las Pruebas que validan este comportamiento**

El siguiente test es raro de leer, pero sirve para mostrar los resultados :

```go
func TestConcurrentFetchProfile(t *testing.T) {
	invalidateCache()

	// 2 Etapas, la primera vamos a llamar en forma concurrente 10
	// veces la función, al no tener valor inicializado, se espera
	// que la primera llamada haga el fetch y luego las posteriores
	// retornen el valor obtenido
	var waitGroup sync.WaitGroup
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := fineFetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 1 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	// Las llamadas anteriores van a funcionar igual mientras el cache sea valido
	// Ahora voy a simular que el cache expiro, pongo un valor con 1 segundo
	// de expire time, y espero un segundo para simular que expiró
	cache = memoize.Memoize(fetchProfile("Expired"), 1*time.Second)
	time.Sleep(2 * time.Second)

	// En las siguientes 10 llamadas, lo que espero es que
	// se retorne el valor cacheado, ya que es valido, sin embargo
	// Se haga fetch de un nuevo valor. por lo tanto el resultado
	// del fetch deberia ser el ultimo, las otras 9 deberia volver el cache
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := fineFetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 2 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	p := fineFetchProfile("Final")
	t.Logf("Value after changes = %s \n", p.Name)
}
```

Ahora como podemos ver en estos tests :

```
=== RUN   TestConcurrentFetchProfile
Fetching Profile... 1        // aca se hace fetch del 2

// Como no hay cache valido, se recupera el 1 los demas se bloquean
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 5 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 1 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 2 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 9 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 0 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 3 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 4 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 8 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 6 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 7 = Profile # 1

// Se esta haciendo fetch de 1
Fetching Profile... 1

// el cache es valido, se retorna el cache, y solo el 1 se retorna cuando
// se tiene el valor
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 2 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 6 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 4 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 0 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 7 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 9 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 3 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 5 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 8 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 1 = Profile # 1
TestConcurrentFetchProfile: service_test.go:69: Value after changes = Profile # 0
// El 1 se devuelve al final
--- PASS: TestConcurrentFetchProfile (3.01s)
```

## Haciendo una librería para generalizar el cache

SafeMemoize es una estructura definida en el archivo safe_memorize.go que nos permite reutilizar la logica anterior.

---

NOTA

En el ejemplo anterior nos perdemos otra oportunidad importante, y la pregunta es porque tengo que esperar para retornar la llamada al profile 1, porque no puedo retornar el valor de cache incluso para profile 1 y llamar el update en una goroutine ?

Este ejemplo final resuelve ese problema, pero te dejo a ti ver el código como lo resuelvo.

---

Tiene un solo método que deberíamos usar :

```go
// Value get cached value, fetching data if needed
func (m *SafeMemoize) Value(
	fetchFunc func() *Memo,
) interface{} {
	...
}
```

La única pieza de código que nos exige, es una función que recupera los datos remotos y retorna un nuevo Memo para actualizar, la mayoría de las veces esta función va a ser un closure como el siguiente :

```go
var profileMemoize = memoize.NewSafeMemoize()

// FetchProfile fetch the current profile
func FetchProfile(id string) *Profile {
	return profileMemoize.Value(
		func() *memoize.Memo {
			data := fetchProfile(id)
			return memoize.Memoize(data, 10*time.Minute)
		},
	).(*Profile)
}
```

Y eso es todo, ahora tenemos toda la lógica de cache generalizada y podemos reutilizarla donde sea necesario.

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index/blob/main/README.md)
