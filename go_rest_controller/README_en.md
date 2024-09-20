[Versión en Español](README.md)

# REST Controllers in go

This notes are about a single and effective way to organize the rest services in a microservices environment.

## RMR - Resource-Method-Representation

RMR is an MVC variant.

With RMR the controllers are organized around Http request definitions.

Each file contains a single route definition, only one entry and all the code that the controllers needs to send response.

File names represent HTTP Rest entries very clear :

```go
 - controllers
       get_hello_username.go
       get_ping.go
       rest.go
```

Where rest.go contains the framework initialization only.

```go
func Start() {
	getRouter().Run(":8080")
}

var router *gin.Engine = nil

func getRouter() *gin.Engine {
	if router == nil {
		router = gin.Default()
	}

	return router
}
```

We have a single file for each rest entry (get_resource, put_resource, get_resources_id, etc) each one handles the single request.

```go
// Internal configure ping/pong service
func init() {
	getRouter().GET("/ping", pingHandler)
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"answer": "pong",
	})
}
```

We can note that all the implementation is contained in the same file.

The same concept can be done with other frameworks like GRPC or any message framework.

### Pros

- It simplifies ode structure from controllers
- Orient our apps and business around clean code file separation (one responsability at the time), from the controller
- It's easier to read and find
- Test are simpler
- Encapsulates the controller properly, not having route entries separated from implementations
- Decouples route initializations, getting better maintainable code

## Fundaments

### MVC in Microservicios

Why MS are different tha monoliths ?

- A microservice in general handles only one aspect of our business
- The interfaces are simpler, and limited
- A single microservice defines most of the times has a single communication way, the View and Controller definitions are mixed and we don't focus too much in that layer separation.
- If there are different Views most of the time those are different Gateways.

### The Controller

The classic MVC, Views and controllers has to do:

- Interpret Requests
- Validate the entry data
- Adapt the request to the business
- Call the model
- Adapt the response to the client
- Handle errors

And in general there is a framework that support us to do all these things.

### Final notes on REST

I'm will not detail REST protocols, just will give a concept separation about ways to organize routes, we can be :

#### Resource Centric

Where the key is the resource, it is called [RestFul](https://en.wikipedia.org/wiki/Representational_state_transfer).

#### Or Use Case Centric

When the RestFul verbs are not enough to represent the use cases, over a resource, it's is better to adopt this path definition, for example a receipt, re can pay, decline, add articles, print, and so. So it's important to specify a protocol based on use cases.

There is no standard, and there are several ways, one that I like is :

- We use GET to ask for data (queries)
- We use POST to change the server state (commands)
- We add a substantive, at the end of GET to know the resource to get
- We add a verb at the end to know with case use run on a resource update (POST)

Examples :

GET http://go_rest_controller.com/order/stats

GET http://go_rest_controller.com/order/totals

GET http://go_rest_controller.com/order/:id/receipt

GET http://go_rest_controller.com/order/:id/details

GET http://go_rest_controller.com/order/:id/totals

POST http://go_rest_controller.com/order/:id/send_email

POST http://go_rest_controller.com/order/:id/paid

POST http://go_rest_controller.com/order/:id/cancel

## Resources

[Introducing the RMR Web Architecture](https://www.peej.co.uk/articles/rmr-architecture.html)

[Chapter 4. The Resource-Oriented Architecture](https://www.oreilly.com/library/view/restful-web-services/9780596529260/ch04.html)

[Modelo-vista-controlador](https://es.wikipedia.org/wiki/Modelo%E2%80%93vista%E2%80%93controlador)

[REST Resource Naming Guide](https://restfulapi.net/resource-naming/)

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](../README_en.md)
