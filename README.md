# goblin
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)

A golang http router based on trie tree.

# Features
- Go 1.13
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

`GET/POST/PUT/PATCH/DELETE`

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

`[name:pattern]`

```go
r.GET(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

A default pattern is wildcard.

`(.+)`

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
r.GET(`/foo/:id[^\w+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\w+$]`)
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
go version: 1.14

```go
go test -bench=.
```

Run a total of 203 routes of GithubAPI.

Tested routes:
- [beego/mux](https://github.com/beego/mux)
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- [dimfeld/httptreemux](https://github.com/dimfeld/httptreemux)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [go-chi/chi](https://github.com/go-chi/chi)

Memory Consumption:
```
goblin: 62896 Bytes
beego-mux: 107328 Bytes
HttpRouter: 37096 Bytes
httptreemux: 78896 Bytes
gin: 59128 Bytes
chi: 71528 Bytes
```

Benchmark Results:
```
BenchmarkGoblin-4                           1363            793539 ns/op         1056676 B/op       3455 allocs/op
BenchmarkBeegoMux-4                         1420            972325 ns/op         1142023 B/op       3475 allocs/op
BenchmarkHttpRouter-4                       1393            878506 ns/op         1021039 B/op       2604 allocs/op
BenchmarkHttpTreeMux-4                      1459            832753 ns/op         1073111 B/op       3108 allocs/op
BenchmarkGin-4                              1033            974177 ns/op         1014084 B/op       2642 allocs/op
BenchmarkChi-4                              1239            846166 ns/op         1095217 B/op       3047 allocs/op
BenchmarkGoblinRequests-4                     60          19327539 ns/op          884059 B/op      11221 allocs/op
BenchmarkBeegoMuxRequests-4                   58          19789097 ns/op          969311 B/op      11241 allocs/op
BenchmarkHttpRouterRequests-4                 58          21749265 ns/op          848012 B/op      10370 allocs/op
BenchmarkHttpTreeMuxRequests-4                57          24634215 ns/op          900326 B/op      10874 allocs/op
BenchmarkHttpGinRequests-4                    45          24686299 ns/op          840777 B/op      10405 allocs/op
BenchmarkHttpChiRequests-4                    45          22363389 ns/op          921963 B/op      10811 allocs/op
```

# Router design
Router accepts requests and dispatches handlers.

![architecture](https://user-images.githubusercontent.com/13291041/79637830-30750980-81bd-11ea-8815-93f7cd104499.png)

goblin based on trie tree structure.

![trie](https://user-images.githubusercontent.com/13291041/79637833-34089080-81bd-11ea-8af4-f0f3f7c2fedc.png)

# Contribution
Please take a look at our [CONTRIBUTING.md](https://github.com/bmf-san/goblin/blob/master/CONTRIBUTING.md) file.

# License
This project is licensed under the terms of the MIT license.

## Author

bmf - A Web Developer in Japan.

-   [@bmf-san](https://twitter.com/bmf_san)
-   [bmf-tech](http://bmf-tech.com/)
