package goblin

import (
	"fmt"
	"net/http"
)

// middleware represents the singular of middleware.
type middleware func(http.Handler) http.Handler

// middlewares represents the plural of middleware.
type middlewares []middleware

func NewMiddlewares(mws middlewares) middlewares {
	return append([]middleware(nil), mws...)
}

// then executes middlewares.
func (mws middlewares) then(h http.Handler) http.Handler {
	for i := range mws {
		h = mws[len(mws)-1-i](h)
	}
	return h
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
