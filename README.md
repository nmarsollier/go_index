# Una guía sobe GO en ambientes de microservicios

Go es un lenguaje muy particular, es y no es muchas cosas, tiene una comunidad muy compleja de entender.

Si ya leíste la documentación oficial de Go, y algún que otro tutorial, pero no sabes bien como desarrollar en Go, este sitio puede ayudarte.

Este es una serie de notas sobre patrones de programación y arquitectura aplicadas a Go, mayormente para ser usadas en microservicios y sistemas pequeños, donde planteo estrategias no tan convencionales de programación en Go.

Y digo no tan convencionales, porque son tips y soluciones no tan populares, pero muy efectivas, aunque no por eso carecen de fundamentos, sino mas bien todo lo contrario, intento proponer un balance ideal entre código sencillo, eficiente y mantenible.

Son ideas pensadas para resolver problemas específicos o en contextos puntuales, por lo tanto no generalizar al leer las notas. 

---
Un buen arquitecto es aquel que diseña arquitecturas sencillas que resuelven el problema elegantemente, que todos pueden entender. 

---

## Tabla de Contenidos

En general en cada apartado ilustro un concepto en particular, pudiendo no respetar otros que no se discuten en esos ejemplos. Recomiendo leerlos en orden, para obtener un contexto adecuado en los capítulos posteriores.

[DI y IoC](https://github.com/nmarsollier/go_di_ioc)

[Un enfoque mas funcional](https://github.com/nmarsollier/go_functional)

[REST Controllers en go](https://github.com/nmarsollier/go_rest_controller)

[Router Design Pattern](https://github.com/nmarsollier/go_router_design)

[Un poco mas declarativos](https://github.com/nmarsollier/go_declarative)

[Builder Pattern en Router](https://github.com/nmarsollier/go_router_builder)

[Polimorfismo con Funciones](https://github.com/nmarsollier/go_functional_polimorfism)

[Estructuras de Directorios](https://github.com/nmarsollier/go_directories)

*Estructuras escalables* (En preparación) 

## Que me motiva a escribir estas notas ?

Con la necesidad de escalar en forma de Microservicios, Go se ha vuelto popular. Sin embargo existe un vacío enorme sobre como implementar correctamente un microservicio.

Estamos viviendo momentos donde existe demasiada información, muchos autores reescriben la rueda, creando conceptos y soluciones extravagantes, y generalmente son mal interpretadas.

Lineamientos como Clean Architecture, Domain Driven Design, son geniales, intentan poner los pies en la tierra, pero exponen información sin un contexto simple de definir, y con tanto ruido en el medio los programadores terminan con mas problemas que soluciones.

Microservicios nos abre un mundo nuevo, cada microservicio es un sistema que ataca un problema puntual, esto nos beneficia en gran medida, porque nuestro microservicio expone una interfaz y nos da libertar de implementar internamente usando el código justo y necesario que resuelve el problema.

Si no tenemos en cuenta el contexto de un microservicio y programamos con la misma receta, terminamos abarrotados, implementando patrones que resuelven problemas que no tenemos, pensamos en diseños hexagonales, desacoplamos el código, encapsulamos negocio, preparamos nuestra app para que sea políglota, y muchas otras cosas "para que escale bien", que son precisamente los problemas que en microservicios se resuelven desde la arquitectura y no desde un microservicio puntual.

El error mas común, es que muchos consideran que mientras mas patrones pongamos en nuestro código, mejor se vuelve, las ideas son buenas cuando resuelven un problema, si no tenemos ese problema, la idea no sirve. 


## Mas

Un microservicio simple, pero que expresa algunas de las notas en este tutorial es [imagego](https://github.com/nmarsollier/imagego).

Puedes ver algunas notas sobre mi perfil en [el indice](https://github.com/nmarsollier/index).
