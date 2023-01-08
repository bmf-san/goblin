package goblin

import (
	"context"
	"errors"
	"net/http"
	"path"
	"strings"
)

// Router represents the router which handles routing.
type Router struct {
	tree                    *tree
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	DefaultOPTIONSHandler   http.Handler
	globalMiddlewares       middlewares
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

	// NOTE: want to separate this from the error when the parameter is not found.
	// Error for not found.
	ErrNotFound = errors.New("no matching route was found")
	// Error for method not allowed.
	ErrMethodNotAllowed = errors.New("methods is not allowed")
)

// NewRouter creates a new router.
func NewRouter() Router {
	return Router{
		tree: newTree(),
	}
}

func (r *Router) UseGlobal(mws ...middleware) {
	nm := NewMiddlewares(mws)
	r.globalMiddlewares = nm
}

// Use sets middlewares.
func (r Router) Use(mws ...middleware) Router {
	nm := NewMiddlewares(mws)
	tmpRoute.middlewares = nm
	return r
}

func (r Router) Methods(methods ...string) Router {
	tmpRoute.methods = append(tmpRoute.methods, methods...)
	return r
}

// Handler sets a handler.
func (r Router) Handler(path string, handler http.Handler) {
	tmpRoute.handler = handler
	tmpRoute.path = path
	r.Handle()
}

// Handle handles a route.
func (r Router) Handle() {
	r.tree.Insert(tmpRoute.methods, tmpRoute.path, tmpRoute.handler, tmpRoute.middlewares)
	tmpRoute = &route{}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	if method == http.MethodOptions {
		if r.DefaultOPTIONSHandler != nil {
			r.DefaultOPTIONSHandler.ServeHTTP(w, req)
			return
		}
	}
	path := cleanPath(req.URL.Path)
	action, params, err := r.tree.Search(method, path)
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

	h := action.handler
	// append globalMiddlewares to head of middlewares.
	mws := append(r.globalMiddlewares, action.middlewares...)
	if mws != nil {
		h = mws.then(action.handler)
	}
	if params != nil {
		ctx := context.WithValue(req.Context(), ParamsKey, params)
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

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// This method borrowed from from net/http package.
// see https://cs.opensource.google/go/go/+/master:src/net/http/server.go;l=2310;bpv=1;bpt=1
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}
