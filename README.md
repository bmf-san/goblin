# goblin
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)

A golang http router based on trie tree.

# Features
- Go 1.16
- Easy to use
- Lightweight
- Fully compatible with net/http
- No external dependencies
- Support named parameters with an optional regular expression.

# Install
```sh
go get -u github.com/bmf-san/goblin
```

# Usage
## Basic
goblin supports these http methods.

`GET/POST/PUT/PATCH/DELETE/OPTION`

You can define routing like this.

```go
r := goblin.NewRouter()

r.GET(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/")
}))

r.POST(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/")
}))
```

## Named parameters
You can use named parameters like this.

```go
r := goblin.NewRouter()

r.GET(`/foo/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))

r.POST(`/foo/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    name := goblin.GetParam(r.Context(), "name")
    fmt.Fprintf(w, "/foo/%v", name)
}))
```

## Named parameters with regular expression
You can also use named parameter with regular expression like this.

`:paramName[pattern]`

```go
r.GET(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

Since the default pattern is `(.+)`, if you don't define it, then `:id` is defined as `:id[(.+)]`.

## Note
A routing pattern matching priority depends on an order of routing definition.

```go
r := goblin.NewRouter()

r.GET(`/foo/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id`)
}))
r.GET(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\d+$]`)
}))
r.GET(`/foo/:id[^\D+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\D+$]`)
}))
```

# Examples
```go
package main

import (
	"fmt"
	"net/http"

	goblin "github.com/bmf-san/goblin"
)

func main() {
	r := goblin.NewRouter()

	r.GET(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.GET(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.GET(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.GET(`/foo/bar/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.GET(`/foo/bar/:id/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.GET(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.GET(`/foo/:id[^\d+$]/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))

	http.ListenAndServe(":9999", r)
}
```

If you want to try it, you can use an [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Benchmark
## Environment
goblin: [1.0.0](https://github.com/bmf-san/goblin/releases/tag/1.0.0)
Golang version: 1.14
Model Name: MacBook Air
Model Identifier: MacBookAir8,1
Processor Name: Dual-Core Intel Core i5
Processor Speed: 1.6 GHz
Number of Processors: 1
Total Number of Cores: 2
Memory: 16 GB

## Test targets
Run a total of 203 routes of GithubAPI.

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
Date: Mon Jul 27 23:27:31 JST 2020

```
GithubAPI Routes: 203
   goblin: 62768 Bytes
   beego-mux: 108224 Bytes
   HttpRouter: 37096 Bytes
   httptreemux: 78800 Bytes
   gin: 59128 Bytes
   chi: 71528 Bytes
goos: darwin
goarch: amd64
pkg: github.com/bmf-san/goblin/benchmark
BenchmarkGoblin-4                           1035            999689 ns/op         1056674 B/op       3455 allocs/op
BenchmarkBeegoMux-4                         1431            823894 ns/op         1142024 B/op       3475 allocs/op
BenchmarkHttpRouter-4                       1533            702788 ns/op         1021037 B/op       2603 allocs/op
BenchmarkHttpTreeMux-4                      1510            790050 ns/op         1073112 B/op       3108 allocs/op
BenchmarkGin-4                              1674            739079 ns/op         1007579 B/op       2438 allocs/op
BenchmarkChi-4                              1452            868848 ns/op         1095208 B/op       3047 allocs/op
BenchmarkGoblinRequests-4                     57          20326844 ns/op          883953 B/op      11220 allocs/op
BenchmarkBeegoMuxRequests-4                   50          23407305 ns/op          969482 B/op      11241 allocs/op
BenchmarkHttpRouterRequests-4                 51          24262251 ns/op          848098 B/op      10369 allocs/op
BenchmarkHttpTreeMuxRequests-4                50          20983605 ns/op          900222 B/op      10872 allocs/op
BenchmarkHttpGinRequests-4                    48          21747666 ns/op          834644 B/op      10202 allocs/op
BenchmarkHttpChiRequests-4                    55          21271806 ns/op          922561 B/op      10813 allocs/op
PASS
ok      github.com/bmf-san/goblin/benchmark     24.316sts-4                    45          22363389 ns/op          921963 B/op      10811 allocs/op
```

# Router design
Router accepts requests and dispatches handlers.

![architecture](https://user-images.githubusercontent.com/13291041/79637830-30750980-81bd-11ea-8815-93f7cd104499.png)

goblin based on trie tree structure.

![trie](https://user-images.githubusercontent.com/13291041/79637833-34089080-81bd-11ea-8af4-f0f3f7c2fedc.png)

# Contribution
We are always accepting issues, pull requests, and other requests and questions.
We look forward to your contribution！

# License
This project is licensed under the terms of the MIT license.

## Author

bmf - A Web Developer in Japan.

-   [@bmf-san](https://twitter.com/bmf_san)
-   [bmf-tech](http://bmf-tech.com/)
