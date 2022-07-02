package main

import (
	"fmt"
	"net/http"

	goblin "github.com/bmf-san/goblin"
)

func customMethodNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotFound")
	})
}

func customeMethodNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotAllowed")
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

func main() {
	r := goblin.NewRouter()
	r.NotFoundHandler = customMethodNotFound()
	r.MethodNotAllowedHandler = customeMethodNotAllowed()

	r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/middleware\n")
		fmt.Fprintf(w, "middleware\n")
	}))
	r.Methods(http.MethodGet).Use(second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/middlewares\n")
		fmt.Fprintf(w, "middlewares\n")
	}))
	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/\n")
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodGet).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/foo\n")
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodGet).Handler(`/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/bar\n")
		fmt.Fprintf(w, "/bar")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/foo/bar\n")
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprint(w, "/foo/bar/:id\n")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar/:id/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprint(w, "/foo/bar/:id/:name\n")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		fmt.Fprint(w, "/foo/:id[^\\d+$]\n")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := goblin.GetParam(r.Context(), "id")
		name := goblin.GetParam(r.Context(), "name")
		fmt.Fprint(w, "/foo/:id[^\\d+$]/:name\n")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))

	http.ListenAndServe(":9999", r)
}
