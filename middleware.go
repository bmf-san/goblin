package goblin

import (
	"net/http"
)

// middleware represents the singular of middleware.
type middleware func(http.Handler) http.Handler

// Middlewares represents the plural of middleware.
type middlewares []middleware

// then executes middlewares.
func (mws middlewares) then(h http.Handler) http.Handler {
	for i := range mws {
		h = mws[len(mws)-1-i](h)
	}
	return h
}
