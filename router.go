package goblin

import (
	"context"
	"fmt"
	"net/http"
)

// Router represents the router which handles routing.
type Router struct {
	tree *Tree
}

// route represents the route which has data for a routing.
type route struct {
	methods     []string
	path        string
	handler     http.Handler
	middlewares middlewares
}

var tmpRoute = &route{}

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
	if err != nil {
		http.Error(w, fmt.Sprintf(`"Access %s: %s"`, path, err), http.StatusNotImplemented)
		return
	}

	h := result.handler

	if result.middlewares != nil {
		h = result.middlewares.then(result.handler)
	}

	if result.params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, result.params)
		req = req.WithContext(ctx)
	}

	h.ServeHTTP(w, req)
}
