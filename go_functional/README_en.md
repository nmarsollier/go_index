[Versión en Español](README.md)

# A functional way to Go

On microservices environments, most of out code is private to the project, microservice encapsultaes the code, and communicate that logic using some interface. So , as a black box that it is, It does not has to be complex or scalable.

So for microservices proposes solutions for a single problem, where interfaces are all that the world knows about them.

### We code Go as Java...

A Go struct associated to code, and generalized with an interface is basically an OO object.

You will hear that Go is not Object Oriented languaje, but in the language definition these artifacts are explicitly OO programming.

Now, if instead of code like OO, what if we take a functional approach ?

### Go as functional paradigm ?

Following the single responsibility pattern o interface segregation principle, it should be pretty common to have service definitions with only one public function.

---

Services with mora than one responsibility are hard to follow, read and maintain, having classes and function with single implementations is key.

---

Lets talk about the next struct :

```go
// HelloService this is the business service
type HelloService struct {
	dao IHelloDao
}
```

This is basically encapsulates a pointer in an struct, that points to a function. This king of artifacts are an OO antipattern, very popularized in Java with EJB, when we didn't know how to split layers.

I'm not against struct, but it's only and pass trough over functions, and has no meaning.

What if we avoid all these things, and we just use a direct call to a function where structs does not exists to hold function pointers, but data ?

The main file, does not needs any instance, and we just call a function :

```go
func main() {
	fmt.Println(service.SayHello())
}
```

The DAO is very simple too, it is just a function.

```go
func Hello() string {
	return "Holiiiiis"
}
```

The service is a little bit more complex, but simpler that previous implementations :

```go
// It's a private pointer that allow us to mock it.
var sayHelloFunc func() string = dao.Hello

// SayHello is the business
func SayHello() string {
	return sayHelloFunc()
}
```

Due that the service is something that we need to test, to allow mock the DAO, we define the dao pointer, I don't know better hacks to do it in golang so far.

The test is :

```go
// Mock the dao
sayHelloFunc = func() string {
	return "Hello"
}

assert.Equal(t, "Hello", SayHello())
```

The strategy to use a pointer, comceptually is the same as using an interface, every function exposes an interface in its declaration.

Note that the mocking function is defined in the client package, not in the service that is important.

---

It is a little bit hacky, but if something can be hacky is the test.

It would be good to take the time to compare and note the quantity of code that we are not writhing just using this paradigm, also apps are more efficient.

---

### Pros

This concept of coding Go in a functional way, simplifies the implementation compared to OO, and in the most solutions the balance between simplicity and good design is very good.

Patterns in OO become obsolete, we need to read some functional techniques, but those are simplest than OO, and more intuitive.

---

Functional languages are extremely expressive. In a functional language one does not need design patterns because the language is likely so high level, you end up programming in concepts that eliminate design patterns all together.

---

Some pros:

- Testing is simpler
- The solutions are natural, we don't have to think in complex patterns
- Less concurrency problems
- Simpler and maintainable code
- We can be more declarative
- Less patterns to learn

## Cons

### Parallel testing

There is a problem if we want to execute parallel testing, because changing the function pointer to some mocked result can interfere with some other parallel running test.

Lucky that Go does not run tests in parallel by default, is something that we need to set up, and we can avoid strategically run tests in parallel if we have this problem.

Now if we need to run tests in parallel, and it's not an option, one strategy is to mock only one time, and return different responses accoding to some parameter in the mocked function, so all the tests use the same mock.

```go
fetchUserMock = func(id string) (User, error) {
   // here we check the id, and return different responses
	 // according the the id

```

By contrary, if the function does not get any parameter, we can use the runtime.Caller to check the caller function, and use that to return different result.

### Some things require state

The state handling that and object contains is important in some cases, like a Builder for example, or Memento pattern. Implementing things like that can be hard in FP.

## Personal opinion about OO programming

OO programming is very good, but it complexity is underestimated, when projects scales the refactoring and maintenance of then becomes harder, so in general we end having spaguetti code.

The book Domain Driven Design, express a ver correct way to express OO modeling and implementation, but it's very weird to see a good implementation, most of the time programmers can't follow the complexity.

Many developers understand that the concept of Clean Architecture and DDD is to encapsulate the business rules, and isolate them with interfaces and layers, but forgets the fundamental part, to code sustentable applications we need have a balance between complexity and pattern usage.

We can't build all the systems using the same magic receipt, every implementation requires brain energy to get good results, and that is hard to get, mostly for junior developers.

Another important thing is the continuous regactoring that POO requires, it is higher that functional, most of the time teams does not maintain apps.

Functional paradigm is simple, refactoring functions is simple, testing is also straight forward, layer separation can be easily done, and having in mind single function responsibilities, everything becomes much more easy to understand.

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](../README_en.md)
