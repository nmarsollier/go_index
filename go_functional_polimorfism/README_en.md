[Version en Español](README.md))

# Polymorphism in Functional Paradigm

Polymorphism is the decisive point to create an interface. I will try to express different functional strategies than interfaces to avoid OO style.

## A simple strategy

Lets say that we have a rest service that returns user profile data, but depending on the calling device (mobile or web). the profile picture takes different strategy and then different result values.

A working strategy that works pretty well, and personally I like is :

```go
var GetImage = map[string]func(id string) string{
	"mobile": func(id string) string {
		return fmt.Sprintf("resource_%s", id)
	},
	"web": func(id string) string {
		return fmt.Sprintf("http://www.images.blablabla/images/%s", id)
	},
}
```

Yes, it's a simple map with functions, in languages like Kotlin we have enums, that are really powerful, but in Go, this strategy works perfectly.

That code is used like:

```go
profile.GetImage[device](id)
```

Simple and elegant.

This behavior is simple, a function, but we could try, for example, map a struct with function pointers, we could be doing more than a single action, and would work the same.

Indeed is so good, that allows us to write very direct validation function for available devices, without neuronal cost :

```go
// IsValidDevice checks if device is valid
func IsValidDevice(device string) bool {
	_, ok := GetImage[device]
	return ok
}
```

### Cons

- The function call is not so natural
- Works very well for a function, but is we follow the Single Responsibility principle, we can live with t hat.
- If the module overload big, we should avoid this strategy.
- Mocking tests is hacky, but tests are hacky by definition.

## Lets look at the code that we don't write

### Doing the last example just with a function

```go
func GetImage(device, id string) string{
	switch device {
	case "mobile":
		return fmt.Sprintf("resource_%s", id)
	case "web":
		return fmt.Sprintf("http://www.images.blablabla/images/%s", id)
	}
}
```

It's not a big deal, it's still good, but the cyclomatic complexity of the code is bigger.

Adding new values requires more work.

And we need to provide a function to validate the available devices, that is not so direct.

### What if we do this sample in an OO way

This, without questions, is the first choice most of the time, when we deal with polymorphism.

We should code :

- An interface, that exposes the method
- 2 structs with different implementations
- A Factory function to get the correct implementation, based on the device, using the same switch that the last example.
- A custom validation function, to check device names.

And just writhing the list, I lost the interest to write the code, also i don't want to count the lines of codes needed to do it.

Most of the times, we point that Go does not handles inheritance, and for good or bad, using the map functional polymorphism approach we are doing the strategy in a single module function, like high level module inheritance.

It's highly likely that the next time that you need to write a polimorphism for just a funcion, you will ask yourself, why i'm doing interfaces for this ?

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](../README_en.md)
