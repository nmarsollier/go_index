[Versión en Español](README.md)

# Builder Pattern in Routers

## Can we use router middlewares as builder pattern ?

On previous notes we saw that we can use many handlers to prepare a result, if we think how a builder works :

```go
dialog.NewBuilder().Title("hola").AcceptAction("Aceptar", "ok").Build()
```

That can be expressed in the route as :

```go
	getRouter().GET(
		"/dialog",
		setTitle,
		setAcceptAction,
		build
	)
```

## Handling complex data

We could think in an scenario where the data is generated in a complex way, and then it could be convenient to segregate the response in many handlers.

Lets suppose a route that wants to mix user and profile information, in the same response.

It's not exactly the same case as previous but the solution is the same:

```go
router().GET(
	"/users/:id",
	validateUserName,
	fetchUser,
	fetchProfile,
	build,
)
```

As always, we check parameters at first, then we get necessary information about the user and profile, using the context to hold answers.

Finally we build the response.

Functions fetchUser y fetchProfile, are simple, the call the model, and put the result in the context.

```go
func fetchUser(c *gin.Context) {
	c.Set("user", user.FetchUser())
	c.Next()
}
```

---

NOTE

Caution adding data to go context, we should use them only in controllers, we cannot call services with that context to get that cached values, it would be incorrect because gin could not be the only one controller around.

---

To build the response we get data from context, and put al the data in the response.

```go
func build(c *gin.Context) {
	user := c.MustGet("user").(*user.User)
	profile := c.MustGet("profile").(*profile.Profile)

	c.JSON(http.StatusOK, gin.H{
		"login":  user.Login,
		"access": user.Access,
		"name":   profile.Name,
	})
}
```

## Parallel execution

Model functions contains 1 second delay to simulate along call, If we check the response time, it's around 2 seconds.

We can call both calls in parallel using this single trick:

```go
router().GET(
	"/parallel/users/:id",
	validateUserName,
	inParallel(
		fetchUserInParallel,
		fetchProfileInParallel,
	),
	build,
)
```

we need to code a function to run in parallel, inParallel it is a simple function but effective, response time now is around 1 second.

```go
func inParallel(handlers ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(handlers))

		for _, handler := range handlers {
			go func(handlerFunc gin.HandlerFunc) {
				defer waitGroup.Done()
				handlerFunc(c)
			}(handler)
		}

		waitGroup.Wait()

		c.Next()
	}
}
```

The only thing to notice, is that fetch functions does not call Next, because inParallel does that at the end.

```go
func fetchUserInParallel(c *gin.Context) {
	c.Set("user", user.FetchUser())
}
```

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](../README_en.md)
