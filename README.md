# goblin
[![GitHub release](https://img.shields.io/github/release/bmf-san/goblin.svg)](https://github.com/bmf-san/goblin/releases)
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/goblin)](https://goreportcard.com/report/github.com/bmf-san/goblin)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/goblin.svg)](https://pkg.go.dev/github.com/bmf-san/goblin)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/goblin/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/goblin?badge)

A golang http router based on trie tree.

# Features
- Go 1.16
- Easy to use
- Lightweight
- Fully compatible with net/http
- No external dependencies
- Support named parameters with an optional regular expression
- Support middlewares

# Install
```sh
go get -u github.com/bmf-san/goblin
```

# Usage
## Basic
Goblin supports these http methods.

`GET/POST/PUT/PATCH/DELETE/OPTIONS`

You can define routing as follows.

```go
r := goblin.NewRouter()

r.GET(`/`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/")
})

http.ListenAndServe(":9999", r)
```

## Matching priority
A routing pattern matching priority depends on an order of routing definition.

The one defined earlier takes precedence over the one defined later.

```go
r := goblin.NewRouter()

r.GET(`/foo/:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id`)
}))
r.GET(`/foo/:id[^\d+$]`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\d+$]`)
}))
r.GET(`/foo/:id[^\D+$]`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\D+$]`)
}))

http.ListenAndServe(":9999", r)
```

In the above case, when accessing `/foo/1`, it matches the routing defined first.

So it doesn't match the 2nd and 3rd defined routings.


## Named parameters
goblin supports named parameters as follows.

```go
r := goblin.NewRouter()

r.GET(`/foo/:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))

r.POST(`/foo/:name`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    name := goblin.GetParam(r.Context(), "name")
    fmt.Fprintf(w, "/foo/%v", name)
}))

http.ListenAndServe(":9999", r)
```

If you use the named parameters without regular expression as in the above case, it is internally interpreted as a wildcard (`(.+)`) regular expression.

So `:id` is substantially defined as `:id[(.+)]` internaly.

## Named parameters with regular expression
You can also use named parameter with regular expression as follows.

`:paramName[pattern]`

```go
r.GET(`/foo/:id[^\d+$]`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

## Middlewares
goblin supports middlewares.

You can be able to set one or more middlewares.

There is no problem even if you do not set the middleware.

```go
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

```

```go
r := goblin.NewRouter()

r.GET(`/middleware`).Use(first).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middleware\n")
}))
r.GET(`/middlewares`).Use(second, third).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middlewares\n")
}))

http.ListenAndServe(":9999", r)
```

In the above case, accessing `/middleware` will produce ouput similar to the following:

```
first: before
middleware
first: after
```

Accessing `/middlewares` will produce ouput similar to the following:
```
second: before
third: before
middlewares
third: after
second: after
```

## Handling CORS Requests by using middleware
```go
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Pagination-Count, Pagination-Pagecount, Pagination-Page, Pagination-Limit")

		next.ServeHTTP(w, r)
	})
}
```

```go
r.GET(`/`).Use(first).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "CORS")
}))
r.OPTIONS(`/`).Use(CORS).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    return
}))
```

# Examples
See [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Router design
This is a rough sketch of what the router process.

As you know, HTTP routers connect paths and handlers based on a protocol called HTTP.

<img src="https://user-images.githubusercontent.com/13291041/117673347-af344e80-b1e5-11eb-9b89-1ec616cd74a2.png" width="600"><br>

Goblin uses a data structure based on the trie tree.

Goblin internally creates a tree like this when the followin routing is provided:

<img src="https://user-images.githubusercontent.com/13291041/117675621-b3616b80-b1e7-11eb-9c64-8542a0f9c7c2.png" width="600"><br>

<img src="https://user-images.githubusercontent.com/13291041/117675761-d4c25780-b1e7-11eb-9ec7-e78ac0ce142b.png" width="800"><br>

It seems that routers that are more conscious of memory efficiency use a prefix tree (patricia trie), but goblin uses a simple trie tree.

If you want to take a closer look at the tree structure, use a debugger to look at the data structure.

If you find a bug, I would be grateful if you could contribute!

# Contribution
We are always accepting issues, pull requests, and other requests and questions.

We look forward to your contribution！

# License
This project is licensed under the terms of the MIT license.

## Author

bmf - A Web Developer in Japan.

-   [@bmf-san](https://twitter.com/bmf_san)
-   [bmf-tech](http://bmf-tech.com/)
