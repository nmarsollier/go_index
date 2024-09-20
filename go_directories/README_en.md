[Versión en Español](README_en.md)

# Directory Structure

None of the previous notes talks about directory structures, they use any structure just to justify the code in the sample.

It's curious that we don't talk too much about directory structure, just when we are doing it so wrong that we need to reorder all the things.

On this notes I'm proposing a simple way to create directory structures, that are not the only way to do it properly, but works fine.

---

Note that this structures are only valid on small apps or microservices.

---

## Layered programming

All the popular design and architectural patterns uses the layer code separation, some has more or less, but they are at least 3 ones that are common :

- Presentation: They are channels to communicate to the out world
- Business: Where we do the magic
- Application: Frameworks and supporting libraries

From my point of view, for most system these 3 layers are enough.

So thr first directory three level, that is a good starting point can be :

```bash
├── controllers
├── model
└── utils

```

### Presentation layer (controller)

We express all the use cases that our app serves to the out world.

The most important thing is that controllers does everything related with the communication with clients, but none business decision.

If our app exposes only one communication way, like rest, we could write all the code in tis directory, may be not bad at all.

If by the contrary, we implement more than one communication protocol, like rabbit, or protocol buggers, we can make subdirectories :

```bash
├── controllers
│   ├── rest
│   ├── rabbit
│   └── protobuff
```

In this repo sample, I have only one protocol, it's complex, so I choose to split middlewares and routes :

```bash
├── controllers
│   ├── middlewares
│   └── router
```

What should enforce us to split a directory ? , well the amount of files and the purpose on the module. If I open the directory and I have any doubts about what a file does, then maybe we should organize the code better.

### Business layer (model)

Inside this one, we have all business rules.

These business rules must define a public interface that must be clean, and other modules or layers can use them without doubt about what they do.

If we are coding microservices, maybe there is only one thing that out model does, maybe we can put all the business files in the same directory.

If we are modeling many aggregates, we must split them in several directories:

```bash
├── model
│   ├── profile
│   └── user
```

In the same, we have 2 aggregates, profile and user, it's healty to put them in 2 directories.

What we include in each folder depends on the same criteria, if everything is simple, we can have the files : service.go, dao.go, api.go, etc.

Many times the aggregate is complex, it could contain several use case functions, so it's better to put each one, in a different file, so maybe we could create subfolders for dao, schemas, and api.

### Application layer (utils)

All the utilities that we need, goes here.

It's always a good thing to split this in subfolders, but how depends on dev criteria.

This folder will be reestructure in the future for sure, we we must be open to refactor this libraries.

## Other layers

### Communication layer

There are scenarios where communication layer are for input and outputs, like when we consume a send events to rabbit or protocol buffers then the communication is bidirectional.

In those cases, there is a big gray about what is application and what is a controller or api.

It's easy to handle those scenarios from the root level, with clean intentions :

```bash
├── rest
│   ├── middlewares
│   └── routes
├── rabbit
│   ├── consummers
│   └── publishers
├── grpc
│   ├── client
│   └── server
├── model
└── utils

```

### DAO and Apis

In general does not require any particular division, they goes perfectly in the model, or internally in the aggregate package.

When an api or database connection is shared, it's better to model and encapsulate it in it's own package.

Example: feature flags remote reading, that need to be checked in different modules.

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](https://github.com/nmarsollier/go_index/blob/main/README_en.md)
