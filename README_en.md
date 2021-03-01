[Version en Español](https://github.com/nmarsollier/go_index/blob/main/README.md)

# A guide to code Go microservices

Go is a very particular language, it is and not many things, has a complex to understand community.

I you have already read official Go documentation, and others, but you don't know how to properly develop good code, this site can help you.

It talks about a series of coding patterns, and architecture notes applied to Go language, mostly to be used in a small systems and microservices.

All the times I talk about non popular solutions, but very effective, not because they are not popular they don't have fundaments, by contrary, I try to do a good balance between simple code, efficient and maintainable code.

These are advanced ideas, to solve specific problems in an specific context, so please do not generalize then and use them everywhere. 

---
A good architect is the one that designs architectures that solve problems so elegantly simple, that anyone can understand it. 

---

## Content Table

In all these topis I talk about a particular concept, and sometimes to focus in a concept, i could not respect other important concept, so please just evaluate the concept that I'm expressing in the topic. Also note that each repository contains code, so you can see samples.

I recommend to read them in the correct order.

[DI and IoC](https://github.com/nmarsollier/go_di_ioc/blob/main/README_en.md)

[A functional way](https://github.com/nmarsollier/go_functional/blob/main/README_en.md)

[REST Controllers in Go](https://github.com/nmarsollier/go_rest_controller/blob/main/README_en.md)

[Router Design Pattern](https://github.com/nmarsollier/go_router_design/blob/main/README_en.md)

[A little bit more declarative](https://github.com/nmarsollier/go_declarative/blob/main/README_en.md)

[Builder Pattern in Router](https://github.com/nmarsollier/go_router_builder/blob/main/README_en.md)

[Polymorphism in Functional Paradigm](https://github.com/nmarsollier/go_functional_polimorfism/blob/main/README_en.md)

[Directory Structure](https://github.com/nmarsollier/go_directories/blob/main/README_en.md)

## Why I'm writing these notes ?

With the necessity to scale in microservices environments, go has become very popular. But I see a huge vacuum about how to implement a single microservice.

We are living times where there is too much information, many writers that redo the wheel, making new extravagant concepts, that generally are wrong interpreted.

Clean Architecture, Domain Driven Design, guidelines, are wonderful, they try to put feet on the ground, but expose information without too much context, and with the medium noise, developers are implementing code that are mostly problems than solutions.

Microservices opens a new world, each single microservice si a systems that cover a single problem, and that is a huge benefic, because a microservice exposes an interface, and we have the freedom to implement internally the correct code to resolve the problem, as a black box.

But if we code without the microservice context in ming, and we use the same recipe used for monolyths, we ends crowded, implementing patterns that solves problems that we don´t have, we adopt hexagonal designs, we decouple too many layers, we prepare microservices to scale, with polyglot databases, and things like that, and many of those things are things that the microservice architecture solves by itself.

The most common mistake, is to use too much patterns thinking that more is better, the patterns are good when they solve a problem, without the problem, the pattern is wrong.

## Code samples

Academic Microservices samples using the notes of this guide : 

[imagego](https://github.com/nmarsollier/imagego).

[authgo](https://github.com/nmarsollier/authgo).
