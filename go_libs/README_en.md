[Version in español](README.md)

# Shared Libraries

When we implement a specific microservice, this implementation should typically solve a specific business problem, and its interface with the outside world is defined through REST, message broker, or something similar. Its internal definition is not very important in the eyes of the users, and we can afford to be very flexible as we have seen in this guide.

But when we implement shared libraries, such as:

- Factories for remote services
- Error logging
- Access APIs (message brokers or databases)
- Internationalized access
- Access to images and shared resources
- Traceability
- General utilities
- Testing utilities (HTTP or mock databases)

These libraries aim to avoid code repetition in microservices and help development teams centralize decisions and control common logic.

However, they are not developed like the code we have been working on in this guide; instead, we must make them more robust so they can be more flexible and adapt to new implementations.

## Interfaces and Structures

Structures in Go define related data in a specific context.

For example:

```go
type Profile struct {
  ID    string
  Login string
  Name  string
  Web   string
}
```

Profile represents a user's profile.

Interfaces allow us to define behaviors decoupled from implementations. For example, if we define:

```go
type Profile interface {
  ID() string
  Login() string
  Name() string
  Web() string
}

func FetchProfile(id string) Profile {
   ...
}
```

We are defining a series of behaviors that must be fulfilled in the implementations, and there is a function FetchProfile that returns those behaviors.

The implementations of the interfaces are done with structures and methods associated with them:

```go
type profile struct {
  id    string
  login string
  name  string
  web   string
}

func (p *profile) ID() string {
  return p.id
}

func (p *profile) Login() string {
  return p.login
}

func (p *profile) Name() string {
  return p.name
}

func (p *profile) Web() string {
  return p.web
}
```

Then we can define the FetchProfile function as:

```go
func FetchProfile(id string) Profile {
  return &profile{
    id:    id,
    login: "nmarsollier",
    name:  "Nestor Marsollier",
    web:   "https://github.com/nmarsollier/profile",
  }
}
```

> The use of interfaces is what we could call polymorphism in Go, as we can implement different structures with different behaviors, but respecting a set of defined functions.

## And all this for what?

The above scheme is very important in the use of shared libraries, as it allows us to define an expected behavior, but not the details of internal implementation. This gives us some flexibility when adapting to changes and also when implementing customizations for each business.

We encapsulate and decouple, open possibilities for other implementations, and can achieve better testing by being able to better mock calls and results.

## Attention to Implicit Interfaces

In Go, interfaces are implicit, which means we do not need to be specific about which structure implements an interface. Simply by defining all the methods of the interface, the structure is already fulfilling it.

> We must be specific with function names in interfaces. If our function is called 'Serialize() string' in two different interfaces XMLSerializer and JSONSerializer, we could run into an error because for the compiler, if a structure implements Serialize, it is already implementing both interfaces.

## Where do we define the interfaces?

They are usually defined in libraries. Sometimes in microservice implementations, we need to be more flexible and require a definition of what they need to operate, which is basically a very reduced subset of what the library interface exposes, and we want to generalize so it can operate with other options.

For example, a database client or HTTP client might expose many methods, but we do not want to enable all the methods that a database client might expose to our DAOs, so we create an interface with a subset of what we need in our app and operate with that subset to avoid using methods we do not want to use from the original client.

> Some literature recommends defining the interface in each function, so that at each point where a function is defined, the necessary interfaces are defined. While this seems to give us flexibility, in reality, we end up having too many interfaces and a lot of code to maintain. I do not personally recommend it.

## Interface Rules

> Robustness principle: Be conservative with what you do, be liberal with what you accept

Our functions should be able to work correctly with any code structure that implements an interface; however, we should only use what our interface defines. We should never make interfaces that define more than what we use.

> Accept interfaces, return structs

Functions should be defined with parameters and returns in Interfaces; however, our code returns a specific implementation of a structure.

```go
func FetchProfile(id string) Profile {
  return &profile{
    id:    id,
    login: "nmarsollier",
    name:  "Nestor Marsollier",
    web:   "https://github.com/nmarsollier/profile",
  }
}
```

> The empty interface says nothing

We limit the use of interface{} to the minimum possible; we only use it when we cannot map the results of a third-party library or in a very specific implementation such as maps for cache.

## Note

This is a series of notes on simple programming patterns in Go.

[Table of Contents](../README_en.md)
