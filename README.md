[English](https://github.com/bmf-san/goblin) [日本語](https://github.com/bmf-san/goblin/blob/master/README-ja.md)

# goblin
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GitHub release](https://img.shields.io/github/release/bmf-san/goblin.svg)](https://github.com/bmf-san/goblin/releases)
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/goblin)](https://goreportcard.com/report/github.com/bmf-san/goblin)
[![codecov](https://codecov.io/gh/bmf-san/goblin/branch/master/graph/badge.svg?token=ZLOLQKUD39)](https://codecov.io/gh/bmf-san/goblin)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/goblin.svg)](https://pkg.go.dev/github.com/bmf-san/goblin)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/goblin/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/goblin?badge)

A golang http router based on trie tree.

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/d654ddf2b81c2b4123684f93071af0cf559eb0b5.png" alt="goblin" title="goblin" width="250px">

This logo was created by [gopherize.me](https://gopherize.me/gopher/d654ddf2b81c2b4123684f93071af0cf559eb0b5).

# Table of contents
- [goblin](#goblin)
- [Table of contents](#table-of-contents)
- [Features](#features)
- [Install](#install)
- [Example](#example)
- [Usage](#usage)
  - [Method based routing](#method-based-routing)
  - [Named parameter routing](#named-parameter-routing)
  - [Regular expression based routing](#regular-expression-based-routing)
  - [Middleware](#middleware)
  - [Customizable error handlers](#customizable-error-handlers)
  - [Default OPTIONS handler](#default-options-handler)
- [Benchmark tests](#benchmark-tests)
- [Design](#design)
- [Contribution](#contribution)
- [Sponsor](#sponsor)
- [Stargazers](#stargazers)
- [Forkers](#forkers)
- [License](#license)
  - [Author](#author)

# Features
- Go1.20 >= 1.16
- Simple data structure based on trie tree
- Lightweight
  - Lines of codes: 2428
  - Package size: 140K
- No dependencies other than standard packages
- Compatible with net/http
- More advanced than net/http's [Servemux](https://pkg.go.dev/net/http#ServeMux)
  - Method based routing
  - Named parameter routing
  - Regular expression based routing
  - Middleware
  - Customizable error handlers
  - Default OPTIONS handler
- 0allocs
  - Achieve 0 allocations in static routing
  - About 3allocs for named routes
     - Heap allocation occurs when creating parameter slices and storing parameters in context

# Install
```sh
go get -u github.com/bmf-san/goblin
```

# Example
A sample implementation is available.

Please refer to [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Usage
## Method based routing
Routing can be defined based on any HTTP method.

The following HTTP methods are supported.
`GET/POST/PUT/PATCH/DELETE/OPTIONS`

```go
r := goblin.NewRouter()

r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/")
}))

r.Methods(http.MethodGet, http.MethodPost).Handler(`/methods`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        fmt.Fprintf(w, "GET")
    }
    if r.Method == http.MethodPost {
        fmt.Fprintf(w, "POST")
    }
}))

http.ListenAndServe(":9999", r)
```

## Named parameter routing
You can define routing with named parameters (`:paramName`).

```go
r := goblin.NewRouter()

r.Methods(http.MethodGet).Handler(`/foo/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))

r.Methods(http.MethodGet).Handler(`/foo/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    name := goblin.GetParam(r.Context(), "name")
    fmt.Fprintf(w, "/foo/%v", name)
}))

http.ListenAndServe(":9999", r)
```

## Regular expression based routing
By using regular expressions for named parameters (`:paramName[pattern]`), you can define routing using regular expressions.

```go
r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

## Middleware
Supports middleware to help pre-process requests and post-process responses.

Middleware can be defined for any routing.

Middleware can also be configured globally. If a middleware is configured globally, the middleware will be applied to all routing.

More than one middleware can be configured.

Middleware must be defined as a function that returns http.

```go
// Implement middleware as a function that returns http.Handl
func global(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "global: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "global: after\n")
	})
}

func first(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "first: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "first: after\n")
	})
}

func second(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "second: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "second: after\n")
	})
}

func third(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "third: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "third: after\n")
	})
}

r := goblin.NewRouter()

// Set middleware globally
r.UseGlobal(global)
r.Methods(http.MethodGet).Handler(`/globalmiddleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/globalmiddleware\n")
}))

// Use methods can be used to apply middleware
r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middleware\n")
}))

// Multiple middleware can be configured
r.Methods(http.MethodGet).Use(second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middlewares\n")
}))

http.ListenAndServe(":9999", r)
```

A request to `/globalmiddleware` gives the following results.

```
global: before
/globalmiddleware
global: after
```

A request to `/middleware` gives the following results.

```
global: before
first: before
middleware
first: after
global: after
```

A request to `/middlewares` gives the following results.

```
global: before
second: before
third: before
middlewares
third: after
second: after
global: after
```

## Customizable error handlers
You can define your own error handlers.

The following two types of error handlers can be defined

- NotFoundHandler
  - Handler that is executed when no result matching the routing is obtained
- MethodNotAllowedHandler
  - Handler that is executed when no matching method is found

```go
func customMethodNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotFound")
	})
}

func customMethodAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotAllowed")
	})
}

r := goblin.NewRouter()
r.NotFoundHandler = customMethodNotFound()
r.MethodNotAllowedHandler = customMethodAllowed()

http.ListenAndServe(":9999", r)
```

## Default OPTIONS handler
You can define a default handler that will be executed when a request is made with the OPTIONS method.

```go
func DefaultOPTIONSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
	})
}

r := goblin.NewRouter()
r.DefaultOPTIONSHandler = DefaultOPTIONSHandler()

http.ListenAndServe(":9999", r)
```

The default OPTIONS handler is useful, for example, in handling CORS OPTIONS requests (preflight requests).

# Benchmark tests
We have a command to run a goblin benchmark test.

Please refer to [Makefile](https://github.com/bmf-san/goblin/blob/master/Makefile).

Curious about benchmark comparison results with other HTTP Routers?

Please see here!
[bmf-san/go-router-benchmark](https://github.com/bmf-san/go-router-benchmark)

# Design
This section describes the internal data structure of goblin.

While [radix tree](https://en.wikipedia.org/wiki/Radix_tree) is often employed in performance-optimized HTTP Routers, goblin uses [trie tree](https://en.wikipedia.org/wiki/Trie).

Compared to radix trees, trie tree have a disadvantage in terms of performance due to inferior memory usage. However, the simplicity of the algorithm and ease of understanding are overwhelmingly in favor of the trie tree.

HTTP Router may seem like a simple application with a simple specification, but it is surprisingly complex. You can see this by looking at the test cases.
(If you have an idea for a better-looking test case implementation, please let us know.)

One advantage of using a simple algorithm is that it contributes to code maintainability. (This may sound like an excuse for the difficulty of implementing a radix tree... in fact, the difficulty of implementing an HTTP Router based on a radix tree frustrated me once...)

Using the source code of [_examples](https://github.com/bmf-san/goblin/blob/master/_examples) as an example, I will explain the internal data structure of goblin.

The routing definitions are represented in a table as follows.

| Method | Path | Handler | Middleware |
| -- | -- | -- | -- |
| GET | / | RootHandler | N/A |
| GET | /foo | FooHandler | CORS |
| POST | /foo | FooHandler | CORS |
| GET | /foo/bar | FooBarHandler | N/A |
| GET | /foo/bar/:name | FooBarNameHandler | N/A |
| POST | /foo/:name | FooNameHandler | N/A|
| GET | /baz | BazHandler | CORS |

In gobin, such routing is represented as the following tree structure.

```
legend：<HTTP Method>,[Node]

<GET>
    ├── [/]
    |
    ├── [/foo]
    |        |
    |        └── [/bar]
    |                 |
    |                 └── [/:name]
    |
    └── [/baz]

<POST>
    └── [/foo]
             |
             └── [/:name]
```

The tree is constructed for each HTTP method.

Each node has handler and middleware definitions as data.

In order to simplify the explanation, data such as named routing data and global middleware data are omitted here.

Various other data is held in the internally constructed tree.

If you want to know more, use the debugger to take a peek at the internal structure.

If you have any ideas for improvements, please let us know!


# Contribution
Issues and Pull Requests are always welcome.

We would be happy to receive your contributions.

Please review the following documents before making a contribution.

[CODE_OF_CONDUCT](https://github.com/bmf-san/goblin/blob/master/.github/CODE_OF_CONDUCT.md)
[CONTRIBUTING](https://github.com/bmf-san/goblin/blob/master/.github/CONTRIBUTING.md)

# Sponsor
If you like it, I would be happy to have you sponsor it!

[GitHub Sponsors - bmf-san](https://github.com/sponsors/bmf-san)

Or I would be happy to get a STAR.

It motivates me to keep up with ongoing maintenance :D

# Stargazers
[![Stargazers repo roster for @bmf-san/goblin](https://reporoster.com/stars/bmf-san/goblin)](https://github.com/bmf-san/goblin/stargazers)

# Forkers
[![Forkers repo roster for @bmf-san/goblin](https://reporoster.com/forks/bmf-san/goblin)](https://github.com/bmf-san/goblin/network/members)

# License
Based on the MIT License.

[LICENSE](https://github.com/bmf-san/goblin/blob/master/LICENSE)

## Author
[bmf-san](https://github.com/bmf-san)

- Email
  - bmf.infomation@gmail.com
- Blog
  - [bmf-tech.com](http://bmf-tech.com)
- Twitter
  - [bmf-san](https://twitter.com/bmf-san)


