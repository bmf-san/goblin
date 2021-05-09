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
	method      string
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

// GET sets a route for GET method.
func (r *Router) GET(path string) *Router {
	tmpRoute.method = http.MethodGet
	tmpRoute.path = path
	return r
}

// POST sets a route for POST method.
func (r *Router) POST(path string) *Router {
	tmpRoute.method = http.MethodPost
	tmpRoute.path = path
	return r
}

// PUT sets a route for PUT method.
func (r *Router) PUT(path string) *Router {
	tmpRoute.method = http.MethodPut
	tmpRoute.path = path
	return r
}

// PATCH sets a route for PATCH method.
func (r *Router) PATCH(path string) *Router {
	tmpRoute.method = http.MethodPatch
	tmpRoute.path = path
	return r
}

// DELETE sets a route for DELETE method.
func (r *Router) DELETE(path string) *Router {
	tmpRoute.method = http.MethodDelete
	tmpRoute.path = path
	return r
}

// OPTION sets a route for OPTION method.
func (r *Router) OPTION(path string) *Router {
	tmpRoute.method = http.MethodOptions
	tmpRoute.path = path
	return r
}

// Handler sets a handler.
func (r *Router) Handler(handler http.Handler) {
	tmpRoute.handler = handler
	r.Handle()
}

// Handle handles a route.
func (r *Router) Handle() {
	r.tree.Insert(tmpRoute.method, tmpRoute.path, tmpRoute.handler, tmpRoute.middlewares)
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

	if result.method != method {
		http.Error(w, fmt.Sprintf(`"Access %s: %s"`, path, err), http.StatusNotFound)
		return
	}

	if result.params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, result.params)
		req = req.WithContext(ctx)
	}

	h.ServeHTTP(w, req)
}
