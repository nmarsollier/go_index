[Versión en Español](https://github.com/nmarsollier/go_di_ioc/blob/main/README.md)

# DI and IoC in GO

This repository talks about alternatives of dependency injection, in Go. 

## Dependency Injection

_Spoiler alert: This is what we can change_

Is that IoC strategy that allows us to insert dependencies in a class, to be used internally by the class.

In the folder [ejemplo_tradicional](./ejemplo_tradicional/) we have code samples.

Most writers recommends to use dependency injection to split and decouple logic layers.

In Go the mos common strategy is to use Dependency Injection by constructor function.

The code looks like: 

```go
srv := service.NewService(dao.NewDao())
fmt.Println(srv.SayHello())
```

where Service is something like :

```go
// IHelloDao DAO interface to inject service
type IHelloDao interface {
	Hello() string
}

// HelloService this is the business service client
type HelloService struct {
	dao IHelloDao
}

// NewService Its the factory, and depends on a IHelloDao implementation
func NewService(dao IHelloDao) *HelloService {
	return &HelloService{dao}
}

// SayHello is the business method
func (s HelloService) SayHello() string {
	return s.dao.Hello()
}
```

It's implemented passing the service dependency instance in the client constructor.

According to the bibliography this pattern allows us to :

- Decouple code, so clients can be extensible
- Reduce code complexity.
- Clients are independent.
- Allows code reusability.
- It is more testable.
 
_And that is true but up at the point_

Because we don't fully decouple, by contrary, we couple much more, our code has to create new instances in some bottstrap method, and in wrong places. So we couple a all the code in wrong files. 

And the client and services are not decoupled, services need to implement an interface.

## Use of Factory Methods as IoC

Lets look how we can do it better.

The Inversion of Control strategy, basically means to have a framework that builds dependencies on demand when the code need it, and the dependency is obtained from the "context".

A service locator, is a pattern that basically is that framework, that knows dependencies, and brings dependencies when they are needed, but it has the same bootstrap problem, it couples all the services factories in a single place.

Checking the responsibility assignment patterns expressed in GRASP, one of the classic ways and correct of IoC use is to use service Factory Methods.

Just think in a factory method, as a framework that build dependencies, and depending on the context, it brings to the client the correct instance of a service where it's needed.

The correct why to use this pattern, is to write the build strategy (factory method) in the same service package, next to the service implementation, so the build strategy is clean, and in the correct place.

Using factory method, we can avoid dependency injections by client constructors, and leave the client to get the right service implementation.

Lets check the code in the main method: the client creation is not coupled anymore with the service, it's simple, and decouples the main method to create a dependency.

```go
	srv := service.NewService()

	fmt.Println(srv.SayHello())
```

That is aligned to the expert pattern.

Lets check the instance creation :

```go
func NewService() *HelloService {
	return &HelloService{
		dao.NewDao(),
	}
}
```

The service fetch a Dao implementation from the Dao factory, that is the artifact that knows how to build the dao, and that is out IoC strategy.
 
```go
func NewDao() *HelloDao {
	return new(HelloDao)
}
```

Inside this factory we can use many different creational patterns, like singleton, object pool, new instance, or whatever. 

Also there could be many factory functions, not a single one, solving different scenarios. 

To mock tests, we just create the struct.

```go
func TestSayHelo(t *testing.T) {
	// Mockeamos
	mockedDao := new(daoMock)

	s := HelloService{
		mockedDao,
	}
	assert.Equal(t, "Hello", s.SayHello())
}
```

Following the idea to not create interaces where do not existe the strategy pattern, the dao does not have any interface, it's just an structure, to it's easely mocked.

Pros:
- Allow us to encapsulate the code properly, creating the correct dependency in the place it's needed.
- Reduce complexity con constructors.
- Decouple constructors without bootstrap methods.
- We use expert pattern properly.
- We write strategies in the correct file..

## Lets talk about my fundamentals 

Indeed, dependency injection is a good practice, the problem is the way in witch it is done, many strategies exposed in books does not scale, because they end distribute knowledge in incorrect places. (see GRASP patterns)

### The Strategy pattern

One abuse of interfaces in many implementations is to create interfaces that does not do anything.

The strategy pattern is about implement different behaviors through the definition of an interface, so we can fo polymorphism.

The Strategy Pattern gives meaning to constructor Dependency Injection.

We should not use DI by constructor if we don't have strategy. If it really exist an interface and the developer can implement a different bejaviur, so it makes sense.

But if the only implementations are internal, of if there is only one implementation, so Factory Methods are better.

Why I'm expressing this ? Because It's pretty common to observe thigs like :

- Implement interfaces anyway, just to split layers.
- Implement interfaces when there is only one implementation.
- Use  interfaces to mock tests
- Or just because the books says

### What we should really consider is:

- We should not have an interface if there is no polymorphism.
- Neither when the options are limited.
- A test is not an excuse to create an interface or DI by constructor.
- Or when "just is case" we generalize, and we always do DI, we are overcoding the app.
  
### With are the problems when we do DI and it's not needed ?

Having in mind, that dependency injection by Factory Method is a good practice, the constructor way has the problems:

- Overweight the constructors.
- We generate a confusion to developers, opening a door to implement custom solution to something that has not.
- We couple code. For example a main.go method does not need to know that a service needs a dao.
- We do code hard to read, then hard to maintain.

### So when it's ok to do DI by constructor ?

- When we have an strategy (polymorphism) and the client can define the behavior. (Ex: A callback routine).
- When we are coding a module, and the implementation is done outside.
- When we are working in a library, and we want to be friendly with the users.

### Creational alternatives

When we have DI, not necessarily we need use a Factory Method, there are many crational patterns, the important thing is to asociate that creation in the service module definition.

## Resources

[Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection)

[GRASP](https://es.wikipedia.org/wiki/GRASP)

[Service locator pattern](https://en.wikipedia.org/wiki/Service_locator_pattern)

[Strategy (patrón de diseño)](https://es.wikipedia.org/wiki/Strategy_(patr%C3%B3n_de_dise%C3%B1o))

[Patrón de diseño](https://es.wikipedia.org/wiki/Patr%C3%B3n_de_dise%C3%B1o)

[YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)

[KISS](https://en.wikipedia.org/wiki/KISS_principle)

[InterfaceImplementationPair](https://martinfowler.com/bliki/InterfaceImplementationPair.html)

[Foo/FooImpl pairs - stop doing it!](http://wrschneider.github.io/2015/07/27/foo-fooimpl-pairs.html)

[When to Mock](https://blog.cleancoder.com/uncle-bob/2014/05/10/WhenToMock.html)

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](https://github.com/nmarsollier/go_index/blob/main/README_en.md)
