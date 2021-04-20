package goblin

import (
	"net/http"
)

// middleware represents the singular of middleware.
type middleware func(http.Handler) http.Handler

// middlewares represents the plural of middleware.
type middlewares []middleware

func NewMiddlewares(mws ...middleware) middlewares {
	return append([]middleware(nil), mws...)
}

// then executes middlewares.
func (mws middlewares) then(h http.Handler) http.Handler {
	for i := range mws {
		h = mws[len(mws)-1-i](h)
	}
	return h
}
