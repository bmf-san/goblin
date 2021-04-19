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

func main() {
	r := goblin.NewRouter()
	r.Use(first, second)
	r.GET(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "middleware\n")
	}))
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
