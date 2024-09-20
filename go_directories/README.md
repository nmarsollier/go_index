[English version](README_en.md)

# Estructuras de Directorios

En ninguno de los proyectos anteriores he realizado una estructura de directorios razonable, sino mas bien adaptada al ejemplo que quería mostrar.

Es curioso que no hablamos mucho sobre la estructura de directorios, solo cuando lo hacemos muy mal, y nos damos cuenta que no encontramos nada, nos ponemos a reorganizar código.

Es estas notas propongo una opción para estructurar los directorios, que si bien no es la única forma de ordenarnos bien las cosas, al menos es una forma ordenada.

## Programación en Capas

Todos los diseños y arquitectura populares proponen una división por capas, dependiendo de lo que leamos hay diferencias, pero hay al menos 3 capas que resultan características y a la vez útiles :

- Presentación: Son los que comunican con el mundo exterior
- Negocio: Donde encapsulamos toda la lógica de negocios.
- Infraestructura: Frameworks y librerías en general

Desde mi punto de vista para la gran mayoría de los sistemas, estas 3 capas son suficientes.

Entonces la división principal, que seria un buen punto de partida es dividir estas capas de esta forma.

```bash
├── controllers
├── model
└── utils

```

### Capa de Presentación (controller)

Expresamos todos los servicios que nuestra aplicación brinda al exterior, todas las interfaces con el mundo exterior.

Lo mas importante de los controllers es que se encargan de todo lo relacionado a la comunicación con los clientes, pero no toman ninguna decision de negocio, para eso llaman siempre a la capa de servicios.

Si nuestra aplicación solamente expone un protocolo, por ejemplo rest, podríamos poner todo lo relacionado a rest directamente en este directorio, no estaría mal.

Si por el contrario implementamos mas de un tipo de controller, digamos rabbit, react y protocolo buffers, podríamos organizarnos asi :

```bash
├── controllers
│   ├── rest
│   ├── rabbit
│   └── protobuff
```

En este ejemplo puntual, de este repo, si bien tengo un solo protocolo, es complejo, y prefiero separar conceptos reutilizables de rutas no reutilizables.

```bash
├── controllers
│   ├── middlewares
│   └── router
```

Donde middleware son los middlewares del ruteador y router son las rutas que mi sistema proporciona en rest.

Que es lo que me anima a mi a separar en mas directorios, bueno la cantidad y diversidad de archivos que estoy poniendo dentro del paquete. Si abriendo el directorio tengo alguna duda de lo que significa lo que hay adentro, es porque debería ordenarse mejor.

### Capa de Negocio (model)

Dentro de esta carpeta tendríamos toda la lógica de negocio.

Lo mas importante de estas carpetas de negocio, es que las mismas definan una interfaz de métodos públicos claros y que puedan utilizarse sin dudas desde cualquier otro model o controller, deben desconocer por completo los detalles de la comunicación con el cliente.

Si estamos programando un microservicio que realiza una sola cosa puntual, y nada mas, por ejemplo un servidor que almacena imágenes, y que solo podemos agregar y consultar imágenes, podríamos poner todo lo relacionado al negocio en esta misma carpeta.

Si por el contrario estamos modelando mas de un agregado en nuestro negocio, conviene separar esos conceptos en subcarpetas.

```bash
├── model
│   ├── profile
│   └── user
```

En este ejemplo, suponemos que tenemos 2 agregados, profile y user. Bueno es saludable tener 2 carpetas, cada una encapsulara um concepto particular.

Lo que incluimos en cada carpeta dependerá un poco del mismo criterio anterior, si todo es simple y se entiende, podríamos tener los archivos: service.go, dao.go, api.go, etc.

Muchas veces un agregado es muy complejo, podría manejar decenas de casos de uso, por lo que nos conviene separar un poco mas el archivo service.go, en estos casos conviene escribir cada caso de uso en un archivo separado, y hasta quizás incluir subcarpetas para daos, apis y estructuras de datos, por ejemplo.

### Capa de Infraestructura (utils)

Esta es la carpeta de infraestructura, todas las utilidades que sean reutilizables por las otras capas, se ponen en esta carpeta.

Siempre conviene separar conceptos en esta carpeta y la separación depende del criterio del desarrollador.

Es muy común que a lo largo del tiempo se reestructure conforme el sistema evoluciona, por lo tanto hay que buscar una estructura que tenga sentido en este momento y no preocuparse mucho por los detalles de separación.

## Otras capas

### Capa de comunicación

Muchas veces tenemos capas que son de entrada y salida, por ejemplo, cuando nuestros servicios utilizan rabbit, o protocol buffers y la comunicación es bidireccional.

En estos escenarios, existe un gris muy grande entre que es infraestructura, que es controller y que es una api.

Para estos casos especiales nos conviene separar desde la raíz estos términos para que queden bien claras las intenciones.

```bash
├── rest
│   ├── middlewares
│   └── routes
├── rabbit
│   ├── consummers
│   └── publishers
├── grpc
│   ├── client
│   └── server
├── model
└── utils

```

### Daos y Apis

En general no necesitan una división particular, en general son capas que se modelan internamente en la capa de negocio, es una subcapa digamos. Está en cada uno de nosotros, y dependiendo de la complejidad de la lógica de negocio, si queremos o no tener esta capa.

Cuando una api o acceso a base de datos es compartido por varios módulos, lo mejor es modelarlo y encapsularlo como un paquete mas del negocio.

Ejemplo: Lectura de feature flags remotas, que se necesiten chequear desde distintos módulos.

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index/blob/main/README.md)
