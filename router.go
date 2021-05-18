package goblin

import (
	"context"
	"errors"
	"net/http"
)

// Router represents the router which handles routing.
type Router struct {
	tree *tree
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
	if err != nil {
		status := handleErr(err)
		w.WriteHeader(status)
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

func handleErr(err error) int {
	var status int
	switch err {
	case ErrMethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case ErrNotFound:
		status = http.StatusNotFound
	}
	return status
}
