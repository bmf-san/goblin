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

`GET/POST/PUT/PATCH/DELETE/OPTION`

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

# Examples
See [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Benchmark
## Environment
|          key          |                             value                             |
| --------------------- | ------------------------------------------------------------- |
| version               | [1.0.0](https://github.com/bmf-san/goblin/releases/tag/1.0.0) |
| Model Name            | MacBook Air                                                   |
| Model Identifier      | MacBookAir8,1                                                 |
| Processor Name        | Dual-Core Intel Core i5                                       |
| Processor Speed       | 1.6 GHz                                                       |
| Number of Processors  | 1                                                             |
| Total Number of Cores | 2                                                             |
| Memory                | 16 GB                                                         |

## Test targets
Run a total of 203 static routes of GithubAPI.

- [beego/mux](https://github.com/beego/mux)
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- [dimfeld/httptreemux](https://github.com/dimfeld/httptreemux)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [go-chi/chi](https://github.com/go-chi/chi)

## How to run
```sh
cd benchmark
go test -bench . -benchmem
```

## Results
Date: Mon May 10 22:42:43 JST 2021

```
GithubAPI Routes: 203
   goblin: 72184 Bytes
   beego-mux: 110408 Bytes
   HttpRouter: 37088 Bytes
   httptreemux: 78800 Bytes
   gin: 59824 Bytes
   chi: 71528 Bytes
goos: darwin
goarch: amd64
pkg: github.com/bmf-san/goblin/benchmark
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
BenchmarkGoblin-4                            837           1235383 ns/op         1066204 B/op       3455 allocs/op
BenchmarkBeegoMux-4                          837           1291392 ns/op         1147724 B/op       3475 allocs/op
BenchmarkHttpRouter-4                       1110           1089544 ns/op         1024059 B/op       2603 allocs/op
BenchmarkHttpTreeMux-4                       930           1083137 ns/op         1076140 B/op       3108 allocs/op
BenchmarkGin-4                               912           1102566 ns/op         1010552 B/op       2438 allocs/op
BenchmarkChi-4                               892           1571383 ns/op         1101479 B/op       3047 allocs/op
BenchmarkGoblinRequests-4                     42          37165428 ns/op          895791 B/op      11226 allocs/op
BenchmarkBeegoMuxRequests-4                   42          30871852 ns/op          977030 B/op      11246 allocs/op
BenchmarkHttpRouterRequests-4                 40          42418396 ns/op          853732 B/op      10384 allocs/op
BenchmarkHttpTreeMuxRequests-4                40          32165799 ns/op          904872 B/op      10872 allocs/op
BenchmarkHttpGinRequests-4                    39          26844252 ns/op          839744 B/op      10212 allocs/op
BenchmarkHttpChiRequests-4                    45          27507046 ns/op          930594 B/op      10813 allocs/op
PASS
ok      github.com/bmf-san/goblin/benchmark     23.699s
```

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
