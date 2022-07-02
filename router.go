package goblin

import (
	"context"
	"errors"
	"net/http"
)

// Router represents the router which handles routing.
type Router struct {
	tree                    *tree
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
}

// route represents the route which has data for a routing.
type route struct {
	methods     []string
	path        string
	handler     http.Handler
	middlewares middlewares
}

var (
	tmpRoute = &route{}

	// Error for not found.
	ErrNotFound = errors.New("no matching route was found")
	// Error for method not allowed.
	ErrMethodNotAllowed = errors.New("methods is not allowed")
)

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		tree: NewTree(),
	}
}

// Use sets middlewares.
func (r *Router) Use(mws ...middleware) *Router {
	nm := NewMiddlewares(mws)
	tmpRoute.middlewares = nm
	return r
}

func (r *Router) Methods(methods ...string) *Router {
	tmpRoute.methods = append(tmpRoute.methods, methods...)
	return r
}

// Handler sets a handler.
func (r *Router) Handler(path string, handler http.Handler) {
	tmpRoute.handler = handler
	tmpRoute.path = path
	r.Handle()
}

// Handle handles a route.
func (r *Router) Handle() {
	r.tree.Insert(tmpRoute.methods, tmpRoute.path, tmpRoute.handler, tmpRoute.middlewares)
	tmpRoute = &route{}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	result, err := r.tree.Search(method, path)
	if err == ErrNotFound {
		if r.NotFoundHandler == nil {
			http.NotFoundHandler().ServeHTTP(w, req)
			return
		}
		r.NotFoundHandler.ServeHTTP(w, req)
		return
	}

	if err == ErrMethodNotAllowed {
		if r.MethodNotAllowedHandler == nil {
			methodNotAllowedHandler().ServeHTTP(w, req)
			return
		}
		r.MethodNotAllowedHandler.ServeHTTP(w, req)
		return
	}

	h := result.actions.handler
	if result.actions.middlewares != nil {
		h = result.actions.middlewares.then(result.actions.handler)
	}
	if result.params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, result.params)
		req = req.WithContext(ctx)
	}
	h.ServeHTTP(w, req)
}

// methodNotAllowedHandler is a default handler when status code is 405.
func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
