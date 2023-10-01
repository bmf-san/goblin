package goblin_test

import (
	"net/http"

	"github.com/bmf-san/goblin"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func RootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func FooHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func FooBarHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func FooBarNameHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func FooNameHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func BazHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func ExampleListenAndServe() {
	r := goblin.NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, RootHandler())
	r.Methods(http.MethodGet, http.MethodPost).Use(CORS).Handler(`/foo`, FooHandler())
	r.Methods(http.MethodGet).Handler(`/foo/bar`, FooBarHandler())
	r.Methods(http.MethodGet).Handler(`/foo/bar/:name`, FooBarNameHandler())
	r.Methods(http.MethodPost).Use(CORS).Handler(`/foo/:name`, FooNameHandler())
	r.Methods(http.MethodGet).Handler(`/baz`, BazHandler())

	http.ListenAndServe(":9999", r)
}
