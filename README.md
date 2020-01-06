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

r.GET(`/foo/:id/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v/", id)
}))

r.POST(`/foo/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    name := goblin.GetParam(r.Context(), "name")
    fmt.Fprintf(w, "/foo/%v/", name)
}))
```

## Named parameters with regular expression
You can also use named parameter with regular expression like this.

`[name:pattern]`

```go
r.GET(`/foo/:id[^\d+$]/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v/", id)
}))
```

A default pattern is wildcard.

`(.+)`

## Note
A routing pattern matching priority depends on an order of routing definition.

```go
r := goblin.NewRouter()

r.GET(`/foo/:id/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id/`)
}))
r.GET(`/foo/:id[^\d+$]/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\d+$]/`)
}))
r.GET(`/foo/:id[^\w+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `/foo/:id[^\w+$]/`)
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
	r.GET(`/foo/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/")
	}))
	r.GET(`/foo/bar/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar/")
	}))
	r.GET(`/foo/bar/:id/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v/", id)
	}))
	r.GET(`/foo/bar/:id/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v/", id, name)
	}))
	r.GET(`/foo/:id[^\d+$]/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v/", id)
	}))
	r.GET(`/foo/:id[^\d+$]/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v/", id, name)
	}))

	http.ListenAndServe(":8000", r)
}
```

If you want to try it, you can use an [_examples](https://github.com/bmf-san/goblin/blob/master/_examples).

# Router design
A role of an url router.

![A role of an url router](https://user-images.githubusercontent.com/13291041/70861219-30929d80-1f6e-11ea-8e86-114e8ba0942b.png "A role of an url router")

goblin based on trie tree structure.

![Based on trie tree](https://user-images.githubusercontent.com/13291041/70862745-7148e180-1f83-11ea-85d3-2cd8fb4db0d3.png "Based on trie tree")

# Contribution
Please take a look at our [CONTRIBUTING.md](https://github.com/bmf-san/goblin/blob/master/CONTRIBUTING.md) file.

# License
This project is licensed under the terms of the MIT license.

## Author

bmf - A Web Developer in Japan.

-   [@bmf-san](https://twitter.com/bmf_san)
-   [bmf-tech](http://bmf-tech.com/)
