package main

import (
	"fmt"
	"net/http"

	goblin "github.com/bmf-san/goblin"
)

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

func main() {
	r := goblin.NewRouter()
	mw := goblin.NewMiddlewares(first)
	mws := goblin.NewMiddlewares(second, third)
	r.GET(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "middleware\n")
	}), mw)
	r.GET(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "middlewares\n")
	}), mws)
	r.GET(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}), nil)
	r.GET(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}), nil)
	r.GET(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}), nil)
	r.GET(`/foo/bar/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}), nil)
	r.GET(`/foo/bar/:id/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}), nil)
	r.GET(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}), nil)
	r.GET(`/foo/:id[^\d+$]/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}), nil)
	http.ListenAndServe(":9999", r)
}
