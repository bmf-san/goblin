package goblin

import (
	"context"
	"fmt"
	"net/http"
)

// Router is a represents the router handling HTTP.
type Router struct {
	tree *Tree
}

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		tree: NewTree(),
	}
}

// GET sets a route for GET method.
func (r *Router) GET(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodGet, path, handler, mws)
}

// POST sets a route for POST method.
func (r *Router) POST(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodPost, path, handler, mws)
}

// PUT sets a route for PUT method.
func (r *Router) PUT(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodPut, path, handler, mws)
}

// PATCH sets a route for PATCH method.
func (r *Router) PATCH(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodPatch, path, handler, mws)
}

// DELETE sets a route for DELETE method.
func (r *Router) DELETE(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodDelete, path, handler, mws)
}

// OPTION sets a route for OPTION method.
func (r *Router) OPTION(path string, handler http.Handler, mws middlewares) {
	r.Handle(http.MethodOptions, path, handler, mws)
}

// Handle handles a route.
func (r *Router) Handle(method string, path string, handler http.Handler, mws middlewares) {
	r.tree.Insert(method, path, handler, mws)
}

type key int

const (
	// ParamsKey is the key in a request context.
	ParamsKey key = iota
)

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

	if result.params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, result.params)
		req = req.WithContext(ctx)
	}

	h := result.handler

	if result.middlewares != nil {
		h = result.middlewares.then(result.handler)
	}

	h.ServeHTTP(w, req)
}

// GetParam gets parameters from request.
func GetParam(ctx context.Context, name string) string {
	params, _ := ctx.Value(ParamsKey).(Params)

	for i := range params {
		if params[i].key == name {
			return params[i].value
		}
	}

	return ""
}
