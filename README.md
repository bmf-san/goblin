# goblin
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GitHub release](https://img.shields.io/github/release/bmf-san/goblin.svg)](https://github.com/bmf-san/goblin/releases)
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/goblin)](https://goreportcard.com/report/github.com/bmf-san/goblin)
[![codecov](https://codecov.io/gh/bmf-san/goblin/branch/master/graph/badge.svg?token=ZLOLQKUD39)](https://codecov.io/gh/bmf-san/goblin)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/goblin.svg)](https://pkg.go.dev/github.com/bmf-san/goblin)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/goblin/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/goblin?badge)

A golang http router based on radix tree.

# Features
- Support Go1.19 >= 1.16
- Easy to use
- Lightweight
- Fully compatible with net/http
- No external dependencies
- Support custom error handler
- Support method-based routing
- Support variables in URL paths
- Support regexp route patterns
- Support middlewares

# Install
```sh
go get -u github.com/bmf-san/goblin
```

# Usage
## Method-based routing
Goblin supports method-based routing.

`GET/POST/PUT/PATCH/DELETE/OPTIONS`

You can define routing as follows.

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

## Variables in URL paths
goblin supports variabled in URL paths.

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

If you use the named parameters without regular expression as in the above case, it is internally interpreted as a wildcard (`(.+)`) regular expression.

So `:id` is substantially defined as `:id[(.+)]` internaly.

## Regexp route patterns
goblin support regexp route patterns.

`:paramName[pattern]`

```go
r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

## Matching priority
A routing pattern matching priority depends on an order of routing definition.

The one defined earlier takes precedence over the one defined later.

```go
r := goblin.NewRouter()

r.Methods(http.MethodGet).Handler(`/foo/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id`)
}))
r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\d+$]`)
}))
r.Methods(http.MethodGet).Handler(`/foo/:id[^\D+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\D+$]`)
}))

http.ListenAndServe(":9999", r)
```

In the above case, when accessing `/foo/1`, it matches the routing defined first.

So it doesn't match the 2nd and 3rd defined routings.

## Custom error handler.
goblin supports custom error handler.

You can be able to set your customized error handler.

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

r := goblin.NewRouter()

r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middleware\n")
}))
r.Methods(http.MethodGet).Use(second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

### Handling CORS Requests by using middleware
This is an example of using middleware.

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

r.Methods(http.MethodGet).Use(first).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "CORS")
}))
r.Methods(http.MethodOptions).Use(CORS).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    return
}))
```

# Examples
See [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Wiki
See [Wiki](https://github.com/bmf-san/goblin/wiki).

# Benchmark tests
Interested in a comparison with other HTTP routers?
Please take a look here.

[bmf-san/go-router-benchmark](https://github.com/bmf-san/go-router-benchmark)

# Contribution
We are always accepting issues, pull requests, and other requests and questions.

We look forward to your contribution！

# License
This project is licensed under the terms of the MIT license.

## Author

bmf - A Web Developer in Japan.

-   [@bmf-san](https://twitter.com/bmf_san)
-   [bmf-tech](http://bmf-tech.com/)
