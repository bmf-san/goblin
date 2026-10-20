package goblin

import (
	"context"
	"errors"
	"net/http"
)

// Router represents the router which handles routing.
type Router struct {
	tree                    map[string]*tree
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	DefaultOPTIONSHandler   http.Handler
	globalMiddlewares       middlewares
	pending                 *route
}

// route represents the route which has data for a routing.
type route struct {
	methods     []string
	path        string
	handler     http.Handler
	middlewares middlewares
}

var (
	// Error for not found.
	ErrNotFound = errors.New("no matching route was found")
	// Error for method not allowed.
	ErrMethodNotAllowed = errors.New("methods is not allowed")
)

// NewRouter creates a new router.
func NewRouter() *Router {
	return &Router{
		tree:    map[string]*tree{},
		pending: &route{},
	}
}

func (r *Router) UseGlobal(mws ...middleware) {
	nm := NewMiddlewares(mws)
	r.globalMiddlewares = nm
}

// Use sets middlewares.
func (r *Router) Use(mws ...middleware) *Router {
	nm := NewMiddlewares(mws)
	r.pending.middlewares = nm
	return r
}

// Use sets methods.
func (r *Router) Methods(methods ...string) *Router {
	r.pending.methods = append(r.pending.methods, methods...)
	return r
}

// Handler sets a handler.
func (r Router) Handler(path string, handler http.Handler) {
	r.pending.handler = handler
	r.pending.path = path
	r.Handle()
}

// Handle handles a route.
func (r *Router) Handle() {
	for i := 0; i < len(r.pending.methods); i++ {
		_, ok := r.tree[r.pending.methods[i]]
		if !ok {
			r.tree[r.pending.methods[i]] = newTree()
		}
		r.tree[r.pending.methods[i]].Insert(r.pending.path, r.pending.handler, r.pending.middlewares)
	}
	*r.pending = route{}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	if method == http.MethodOptions {
		if r.DefaultOPTIONSHandler != nil {
			r.DefaultOPTIONSHandler.ServeHTTP(w, req)
			return
		}
	}

	t, ok := r.tree[method]
	if !ok {
		if r.MethodNotAllowedHandler == nil {
			methodNotAllowedHandler().ServeHTTP(w, req)
			return
		}
		r.MethodNotAllowedHandler.ServeHTTP(w, req)
		return
	}

	action, params, err := t.Search(req.URL.Path)
	if err == ErrNotFound {
		if r.NotFoundHandler == nil {
			http.NotFoundHandler().ServeHTTP(w, req)
			return
		}
		r.NotFoundHandler.ServeHTTP(w, req)
		return
	}

	h := action.handler
	// append globalMiddlewares to head of middlewares.
	mws := append(r.globalMiddlewares, action.middlewares...)
	if mws != nil {
		h = mws.then(h)
	}
	if params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, params)
		req = req.WithContext(ctx)
		// set path values for compatibility with net/http.
		// iterate in reverse so that the first value wins for duplicated names, as GetParam does.
		for i := len(params) - 1; i >= 0; i-- {
			req.SetPathValue(params[i].key, params[i].value)
		}
	}
	h.ServeHTTP(w, req)
}

// methodNotAllowedHandler is a default handler when status code is 405.
func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
