[Versión en Español](README.md)

# A little bit more declarative

Go is an imperative language so far, in this note I will trey to give some tips to write code in a more declarative way.

## What meas more declarative ?

Go is not a declarative language, so, it is a good question.

To be declarative, is to avoid write code in a procedural way, like :

```go
// imperative
func Suma(start, end int) int {
	suma := 0
	for i := start; i < end; i++ {
		suma += i
	}
	return suma
}

// declarative
suma := Suma(1, 100)
```

The basic idea, is to build a set of generic functions (imperative) that are used as a high level declarative way. So business rules are expressed in a declarative format.

Functional languages are in general declarative. But the concept is so powerful that many OO languages has adopted tools to write code more declarative like Typescript, Kotlin o Swift.

### Procedural programming disappears ?

Not in GO, but we can build libraries and tools that let us express in a declarative way in high level control flow.

## What is the advantage ?

There are many of them, at business code writhing time, the declarative way allows the code to be more expressive, we define what we need and not how, it's easy to read and maintain.

## Example: A function to not repeat

One DRY concept, to avoid code repetition is to build a function that does that part of the code :

```go
func sayHelloHandler(c *gin.Context) {
	userName := c.Param("userName")

	c.JSON(http.StatusOK, gin.H{
		"answer": service.SayHello(userName),
	})
}
```

```go
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"answer": "pong",
	})
}
```

If we write a library like :

```go
func SendJSONAnswer(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"answer": data,
	})
}
```

That allow us to be more declarative, and to not repeat at the time, the new code is :

```go
func sayHelloHandler(c *gin.Context) {
	gu.SendJSONAnswer(c, service.SayHello(c.Param("userName")))
}
```

```go
func pingHandler(c *gin.Context) {
	gu.SendJSONAnswer(c, "pong")
}
```

## The function name is the key

An important aspect that we look for, is to write function names like the natural language.

For example, if instead of call the previous example sayHelloHandler we call it sayHello, it gets more natural:

```go
func init() {
	getRouter().GET(
		"/hello/:userName",
		validateUserName,
		sayHello,
	)
}
```

And we are being more declarative.

## A builder

A builder is a typical use scenario to be declarative.

```go
dialog.NewBuilder().Title("Hola Mundo").AcceptAction("Aceptar", "ok").Build()
```

I think that we like the builder pattern, because it is declarative.

The implementation could be :

```go
package dialog

import (
	"encoding/json"
)

type dialogAction struct {
	Label  string `json:"label,omitempty"`
	Action string `json:"action,omitempty"`
}

type dialog struct {
	Title  *string       `json:"title,omitempty"`
	Accept *dialogAction `json:"accept,omitempty"`
}

type DialogBuilder struct {
	dialog
}

func NewBuilder() *DialogBuilder {
	return &DialogBuilder{dialog{}}
}

func (d *DialogBuilder) Title(value string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title:  &value,
		Accept: d.dialog.Accept,
	}}
}

func (d *DialogBuilder) AcceptAction(label, action string) *DialogBuilder {
	return &DialogBuilder{dialog{
		Title: d.dialog.Title,
		Accept: &dialogAction{
			Label:  label,
			Action: action,
		},
	}}
}

func (d *DialogBuilder) Build() string {
	result, _ := json.Marshal(d.dialog)
	return string(result)
}
```

We could have used a single function, with parameters, but sometimes this pattern is more expressive and simple.

## Chaining behaviors

Could we write all the app declaratively ?

Not sure if worth in Go, but surely many times we will see code that could become simple if we are declarative.

Lets check this imperative function.

```go
func Shorten(name string) string {
	values := strings.Split(name, " ")

	result := ""

	for _, v := range values {
		if len(v) > 0 {
			result += strings.ToUpper(string(v[0]))
		}
	}

	return result
}
```

Have you discover what it does ? What if we define the same thing in a declarative way:

```go
func Shorten(name string) string {
	return fromString(name).
		split(" ").
		mapNotEmpty(func(s string) string {
			return strings.ToUpper(string(s[0]))
		}).
		joinToString()
}
```

Now after read the code we can see what it does, because :

- from an string
- it split the words in an array
- and maps the array in another one accordin to the criteria : take the uppercase first letter of each element
- then combines the array in an string

It simply converts "one little thing" in "OLT".

Now to make that code works, we need these libraries :

```go
type shorten struct{ string }
type shortenSlice struct{ slice []string }

func fromString(value string) shorten {
	return shorten{value}
}

func (s shorten) split(separator string) shortenSlice {
	return shortenSlice{strings.Split(s.string, separator)}
}

func (values shortenSlice) mapNotEmpty(f func(string) string) shortenSlice {
	var result []string
	for _, v := range values.slice {
		if len(v) > 0 {
			result = append(result, f(v))
		}
	}
	return shortenSlice{result}
}

func (values shortenSlice) joinToString() (result string) {
	for _, v := range values.slice {
		result += v
	}
	return result
}
```

A lot of code, but we should consider that the shorten and shortenSlice structs are generics libraries, and we can reuse them, we are focusing in the business rule that is the Shorten function.

It's also good to start from some point like library https://github.com/jucardi/go-streams, that give us some declarative functions to do basic things.

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](../README_en.md)
