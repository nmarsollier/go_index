[English version](https://github.com/nmarsollier/go_index/blob/main/README_en.md)

# Una guía sobe GO en ambientes de microservicios

Este es una serie de notas sobre patrones de programación y arquitectura aplicadas a Go, mayormente para ser usadas en microservicios y sistemas pequeños, donde planteo tips y soluciones no tan populares, muy efectivas, que no por eso carecen de fundamentos, sino mas bien todo lo contrario, intento proponer un balance ideal entre código sencillo, eficiente y mantenible.

---

Un buen arquitecto es aquel que diseña arquitecturas sencillas que resuelven el problema tan elegantemente, que todos pueden entender.

---

## Tabla de Contenidos

En general en cada apartado ilustro un concepto en particular, pudiendo no respetar otros que no se discuten en esos ejemplos. Recomiendo leerlos en orden, para obtener un contexto adecuado en los capítulos posteriores.

[DI e IoC](https://github.com/nmarsollier/go_di_ioc/blob/main/README.md)

> De donde sale la inyección de dependencias ? - Porque no sirve en Go? - Estrategias mas simples ?

[Un enfoque mas funcional](https://github.com/nmarsollier/go_functional/blob/main/README.md)

> Go no es Orientado a Objetos - Beneficios de prog. funcional - Estrategias de mock sin interfaces innecesarias

[Repasemos DI e IoC funcional](https://github.com/nmarsollier/go_di_ioc2/blob/main/README.md)

> En programación funcional, la inyección de dependencias no existe como en OO.

[REST Controllers en go](https://github.com/nmarsollier/go_rest_controller/blob/main/README.md)

> El MVC simple y bien explicado - Organizar correctamente nuestros controllers - Organizar el código desde el controller - REST en servicios de negocio

[Router Design Pattern](https://github.com/nmarsollier/go_router_design/blob/main/README.md)

> Que es el Router Pattern - Porque digo que no lo aprovechamos totalmente - Estrategias

[Un poco mas declarativos](https://github.com/nmarsollier/go_declarative/blob/main/README.md)

> Que es ser declarativo y como le podemos sacar ventaja en simplificar el código

[Builder Pattern en Router](https://github.com/nmarsollier/go_router_builder/blob/main/README.md)

> Como podemos aprovechar el Router mucho mas usándolo como Builder - Procesamiento en paralelo muy simple

[Polimorfismo con Funciones](https://github.com/nmarsollier/go_functional_polimorfism/blob/main/README.md)

> Evitando interfaces innecesarias - Estrategias para hacer polimorfismo con funciones

[Estructuras de Directorios](https://github.com/nmarsollier/go_directories/blob/main/README.md)

> Dividir en capas el código - Como organizar el código inteligentemente para que encontremos cada cosa en su lugar.

[Una forma adecuada de hacer cache](https://github.com/nmarsollier/go_cache/blob/main/README.md)

> Hacer cache no es nada sencillo, se explican problemas comunes y se da una solución de cache para un caso muy puntual, cachear una respuesta remota.

## Que me motiva a escribir estas notas ?

Go es un lenguaje muy particular, es y no es muchas cosas, tiene una comunidad muy compleja de entender.

Si ya leíste la documentación oficial de Go, y algún que otro tutorial, pero no sabes bien como desarrollar en Go, este sitio puede ayudarte.

En arquitecturas de Microservicios, Go se ha vuelto popular. Sin embargo existe un vacío enorme sobre como implementar correctamente un microservicio.

Estamos viviendo momentos donde existe demasiada información, muchos autores reescriben la rueda, creando conceptos y soluciones extravagantes, y generalmente son mal interpretadas.

Lineamientos como Clean Architecture, Domain Driven Design, son geniales, intentan poner los pies en la tierra, pero exponen información sin un contexto simple de definir, y con tanto ruido en el medio los programadores terminan con mas problemas que soluciones.

Microservicios nos abre un mundo nuevo, cada microservicio es un sistema que ataca un problema puntual, esto nos beneficia en gran medida, porque nuestro microservicio expone una interfaz y nos da libertar de implementar internamente usando el código justo y necesario que resuelve el problema.

Si no tenemos en cuenta el contexto de un microservicio y programamos con la misma receta, terminamos abarrotados, implementando patrones que resuelven problemas que no tenemos, pensamos en diseños hexagonales, desacoplamos el código, encapsulamos negocio, preparamos nuestra app para que sea políglota, y muchas otras cosas "para que escale bien", que son precisamente los problemas que en microservicios se resuelven desde la arquitectura y no desde un microservicio puntual.

El error mas común, es que muchos consideran que mientras mas patrones pongamos en nuestro código, mejor se vuelve, las ideas son buenas cuando resuelven un problema, si no tenemos ese problema, la idea no sirve.

## Mas

Ejemplos de microservicios en Go que siguen las notas en este documento :

[imagego](https://github.com/nmarsollier/imagego).

[authgo](https://github.com/nmarsollier/authgo).

Puedes ver algunas notas sobre mi perfil en [el indice](https://github.com/nmarsollier/index).

```

```
