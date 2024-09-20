[Version en Español](README_en.md)

# Router Design Pattern

This repository talks about the effective use of router pattern in REST framework.

## The Pattern

This framework comes with functional languages, and becomed addopted in modern rest frameworks like express, gin, etc.

El concept is easy, we define a route, and fill the handlers that needs to be executed :

```go
func init() {
	getRouter().GET("/ping", pingHandler)
}
```

In it's simples way it's a route and a handler.

## Chain Of Responsibility

But the concept goes far away from there, to understand it, we need to know better the Chain of Responsibility Pattern (CoR).

- It's an old functional programming concept
- Allows us to reuse code
- Split responsibilities in functions
- Allow users to define responsibilities in controllers in a simpler way
- Controllers can build complex responses one at the time

![Chain of Responsibility](./img/cor.png)

As we can see, the idea is to to Como podemos ver, la idea is to break down a big routine in small functions that does only one thing, each function can continue or break the flow.

## Back to Router Design Pattern

The pattern is base on CoR.

The are 2 main concepts, middlewares and handlers

In both we have full control over the request.

### Middlewares

Generally we call them handlers, and the are used by the router in a high level, all the routes are affected by them.

The are useful for :

- Security: Authentication and Authorization
- Loading resources like i18n
- Error handling
- CORS validations
- Block hacking techniques
- Handle transactional contexts
- Preloading of resources related
- Loggers and stats

And much more.

In the repository we can see the error handler in middlewares/errors.go

In the router we setup the global middleware:

```go
router = gin.Default()
router.Use(middlewares.ErrorHandler)
```

And this is the implementation :

```go
// ErrorHandler a middleware to handle errors
func ErrorHandler(c *gin.Context) {
	c.Next()

	handleErrorIfNeeded(c)
}

func handleErrorIfNeeded(c *gin.Context) {
	err := c.Errors.Last()
  ...
```

This handler calls first c.Next, that corresponds to the next call in the chain, once it's executed, it analyzes the response to see if the are errors to handle, if there is one, it responds the according error.

In other implementations we could block the Next execution if needed.

### Route Handlers

It's the same as the middleware, but applies to single routes.

They are useful to :

- Validate requests
- Validate authorization
- Validate body structs and parameters
- Preload related data

Generally they are used to do an early exit if something is wrong about the request.

With them we can cut the chain, forcing errors, in a single way.

And decouple the responsibilities of the controller in single functions, easy to read, test and maintain.

Also we can reuse functions, like parameter validations, due that handlers are functions that we can add or remove from router in any route.

As we see in the file get_hello_username.go

```go
// Internal configure ping/pong service
func init() {
	getRouter().GET(
		"/hello/:userName",
		validateUserName,
		sayHelloHandler,
	)
}

// validamos que el parámetro userName tenga al menos 5 caracteres
func validateUserName(c *gin.Context) {
	userName := c.Param("userName")

	if len(userName) < 5 {
		c.Error(errors.NewCustomError(400, "userName debe tener al menos 5 caracteres"))
		c.Abort()
		return
	}
}
```

function validateUserName is defined as route middleware, validates the url parameter to be correct, and in error case can abort the execution. This strategy is called early exit or guard clause.

Once done all validations, we can run the handler, ensuring that everything is correct.

```go
func sayHelloHandler(c *gin.Context) {
	userName := c.Param("userName")

	c.JSON(http.StatusOK, gin.H{
		"answer": service.SayHello(userName),
	})
}
```

Tests are simpler, because each function has one thing to test :

```go
func TestValidateUserName(t *testing.T) {
	response := test.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/hello/abc", nil)

	validateUserName(context)

	response.Assert(0, "")
	assert.Equal(t, context.Errors.Last().Error(), "userName debe tener al menos 5 caracteres")
}

```

### Middlewares lo preload data

Many times we ned some data preload in the context ,like logged in user profile, to be used later in other parts of the chain.

It's totally feasible, with the next warning: Preload data can only be accesses by the controller, we should never call a service with the context to fetch context info.

An example :

```go
func LoadCurrentUser(c *gin.Context) {
  token, err := c.GetHeader("Authorization")
	if err != nil {
		return
	}

	userProfile := userDao.FindUserByToken(token)

	c.Set("profile", userProfile)
}

func CurrentUserProfile(c *gin.Context) *users.Profile {
	if t, _ := c.Get("profile"); t != nil {
		return t.(*users.Profile)
	}
	return nil
}
```

In the previous same, the same middleware that preload the data LoadCurrentUser, is the one that will look int to the context to get the related data.

When we need to get the current user we call CurrentUserProfile.

## Source

[Routing Design Patterns](https://medium.com/@goldhand/routing-design-patterns-fed766ad35fa)

[Cadena de responsabilidad](https://es.wikipedia.org/wiki/Cadena_de_responsabilidad)

[Guard Clause](https://deviq.com/design-patterns/guard-clause)

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](https://github.com/nmarsollier/go_index/blob/main/README_en.md)
