[Versión en español](README.md)

# After functional, we talk about dependency injection again

We have already seen in previous articles how to apply functional programming, why it is better to use this programming style in Go, now we review dependency injection with what we already know.

## A pattern that does not exist

In functional programming, dependency injection, per se, does not exist, structures are never created with the intention of accessing dependencies.

> Interfaces in Go are not for injecting dependencies but allow implementing the [strategy pattern](../go_di_ioc/README_en.md), however, the strategy to be used should not be passed as parameters, but rather each function should receive the necessary parameters to find the required dependencies.

When we write a function in the functional style, we basically respect the following:

- Functions should do only one thing (this is key)
- Clear and concise name that is self-explanatory
- In Go, short and easy-to-remember name
- 2 or 3 arguments maximum, as long as they are clear and concise
- Functions should only receive the data they need, no more, no less.
- If we have many arguments, pass a structure, it simplifies refactoring, maintenance, and the meaning of the parameters.
- Functions should have only one level of abstraction.
- Structures that are passed as parameters and returned should be immutable
- In general, functions that correspond to the same business layer and are related should be written close together in the code.

## So, how do we pass dependencies?

Passing dependencies is not a problem, each function receives the parameters necessary for its proper functioning, each function has direct access to the functions it depends on, so it is not necessary to send pointers in its parameters, for example, a service and a DAO, the service should have the logic to determine the corresponding DAO to call, this is called the expert information pattern.

More strategies are explained later in this guide.

## What data should a function receive?

Only what it needs, we should never pass data that the function does not need, or complex structures that are not used later, it is always preferable to receive the exact parameters and when they are many or confusing, define and pass those parameters through structures, so that it is clear that a function needs that and only that, no more, no less.

Functions should be seen as closed boxes from the outside world, they are boxes that need certain information because they respond to a very clear functionality, and that data they need is enough for the user of the function to understand what is needed and sometimes just by knowing the data we already know why.

> Always passing a structure as a function parameter is not good practice, we only define a structure as a parameter when the parameters are confusing to read, otherwise, it is better to pass individual parameters.

A very common mistake in HTTP services is to pass the context and have functions extract values from the context, functions should receive the context only to cancel goroutines, for example, but never to extract values from it.

The context is a bag of information that never makes it clear what requirements it must have to be valid, although we can use the context and we must use the context to put values, these values are restricted in their use within the controllers, when we call a service we extract those values and it should be called with the value that the service specifically needs.

## Similar Strategies

A more professional way to do this: In this project, I adopted a strategy that allows passing a variable parameter to functions that represents a "context", but not a Go context, rather a business services context that should be used. This context is generally empty, unless we want to provide different implementations to libraries (e.g., when we do unit testing).

Key concepts of this approach:

- Functions are responsible for creating the necessary services (we do not pass the services to use as a parameter).
- Functions are decoupled from the way services are created.
- Services have a constructor that receives the context (var arg) and based on the context determines the instance to use.
- Functions that need to use a service use the function from the previous point to access those services.

Now, each Service that can have more than one implementation is responsible for

[imagego](https://github.com/nmarsollier/imagego).

## References

[Pitfalls of context values and how to avoid or mitigate them in Go](https://www.calhoun.io/pitfalls-of-context-values-and-how-to-avoid-or-mitigate-them/)

[Twelve Go Best Practices](https://talks.golang.org/2013/bestpractices.slide#1)

## Note

This is a series of notes on simple programming patterns in Go.

[Table of Contents](../README_en.md)
