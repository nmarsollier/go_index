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

## Estrategias similares

Una forma mas profesional de hacer esto: En este proyecto adopte una estrategia que permite pasar un parámetro variable a las funciones que representa un "contexto", pero no un contexto go, sino mas bien un contexto de servicios de negocio que se deben usar, este contexto generalmente esta vacío, salvo que queramos proporcionar diferentes implementaciones a librerías (Ej: Cuando hacemos unit test)

Conceptos clave de este enfoque:

- Las funciones son las que se encargan de crear los servicios necesarios (no les pasamos los servicios a usar por parámetro).
- Las funciones están desacopladas de la forma en la que se crean los servicios.
- Los servicios tienen un constructor que recibe el contexto (var arg) y en base al contexto determina la instancia a usar.
- Las funciones que necesiten usar un servicio usan la función del punto anterior para acceder a esas funciones.

Ahora bien, cada Servicio que puede tener mas de una implementación es el encargado de

[imagego](https://github.com/nmarsollier/imagego).

## Referencias

[Pitfalls of context values and how to avoid or mitigate them in Go](https://www.calhoun.io/pitfalls-of-context-values-and-how-to-avoid-or-mitigate-them/)

[Twelve Go Best Practices](https://talks.golang.org/2013/bestpractices.slide#1)

## Nota

Esta es una serie de notas sobre patrones simples de programación en GO.

[Tabla de Contenidos](../README.md)
