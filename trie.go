package goblin

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)

// tree is a trie tree.
type tree struct {
	node       *node
	paramsPool sync.Pool
}

// node is a node of tree.
type node struct {
	label    string
	actions  map[string]*action // key is method
	children map[string]*node   // key is label of next nodes
}

// action is an action.
type action struct {
	middlewares middlewares
	handler     http.Handler
}

const (
	paramDelimiter    string = ":"
	leftPtnDelimiter  string = "["
	rightPtnDelimiter string = "]"
	ptnWildcard       string = "(.+)"
)

// newTree creates a new trie tree.
func newTree() *tree {
	return &tree{
		node: &node{
			label:    "/",
			actions:  make(map[string]*action),
			children: make(map[string]*node),
		},
	}
}

// Insert inserts a route definition to tree.
func (t *tree) Insert(methods []string, path string, handler http.Handler, mws middlewares) {
	path = cleanPath(path)
	curNode := t.node

	if path == "/" {
		curNode.label = path
		for i := 0; i < len(methods); i++ {
			curNode.actions[methods[i]] = &action{
				middlewares: mws,
				handler:     handler,
			}
		}
		return
	}

	path = removeTrailingSlash(path)

	cnt := strings.Count(path, "/")
	var l string

	for i := 0; i < cnt; i++ {
		// Delete the / at head of path. ex. /foo/bar → foo/bar
		if path[:1] == "/" {
			path = path[1:]
		}

		idx := strings.Index(path, "/")
		if idx > 0 {
			// ex. foo/bar/baz → foo
			l = path[:idx]
		}
		if idx == -1 {
			// ex. foo → foo
			l = path
		}

		nextNode, ok := curNode.children[l]
		if ok {
			curNode = nextNode
			if idx > 0 {
				l = path[:idx]
				// foo/bar/baz → /bar/baz
				path = path[idx:]
			}
		}
		// Create a new node.
		if !ok {
			curNode.children[l] = &node{
				label:    l,
				actions:  make(map[string]*action),
				children: make(map[string]*node),
			}
			curNode = curNode.children[l]
			if idx > 0 {
				l = path[:idx]
				// foo/bar/baz → /bar/baz
				path = path[idx:]
			}
		}
		// last loop.
		// If there is already registered data, overwrite it.
		if i == cnt-1 {
			curNode.label = l
			for j := 0; j < len(methods); j++ {
				curNode.actions[methods[j]] = &action{
					middlewares: mws,
					handler:     handler,
				}
			}
			break
		}
	}
}

// regCache represents the cache for a regular expression.
type regCache struct {
	s sync.Map
}

// getReg gets a compiled regexp from cache or create it.
func (rc *regCache) getReg(ptn string) (*regexp.Regexp, error) {
	v, ok := rc.s.Load(ptn)
	if ok {
		reg, ok := v.(*regexp.Regexp)
		if !ok {
			return nil, fmt.Errorf("the value of %q is wrong", ptn)
		}
		return reg, nil
	}
	reg, err := regexp.Compile(ptn)
	if err != nil {
		return nil, err
	}
	rc.s.Store(ptn, reg)
	return reg, nil
}

var regC = &regCache{}

// Search searches a path from a tree.
func (t *tree) Search(method string, path string) (*action, []Param, error) {
	path = cleanPath(path)
	curNode := t.node

	if path == "/" && curNode.label == "/" && curNode.actions[method] == nil {
		return nil, nil, ErrNotFound
	}

	path = removeTrailingSlash(path)

	cnt := strings.Count(path, "/")
	var l string
	var params *[]Param

	for i := 0; i < cnt; i++ {
		// Delete the / at head of path. ex. /foo/bar → foo/bar
		if path[:1] == "/" {
			path = path[1:]
		}

		idx := strings.Index(path, "/")
		if idx > 0 {
			// ex. foo/bar/baz → foo
			l = path[:idx]
		}
		if idx == -1 {
			// ex. foo → foo
			l = path
		}

		nextNode, ok := curNode.children[l]
		if ok {
			curNode = nextNode
			if idx > 0 {
				l = path[:idx]
				// foo/bar/baz → /bar/baz
				path = path[idx:]
			}
			continue
		}

		cc := curNode.children

		// leaf node
		if len(cc) == 0 {
			if curNode.label != l {
				// no matching path was found.
				return nil, nil, ErrNotFound
			}
			break
		}
		isParamMatch := false
		// parameter matching
		for c := range cc {
			if c[0:1] == paramDelimiter {
				ptn := getPattern(c)
				if ptn != "" {
					reg, err := regC.getReg(ptn)
					if err != nil {
						return nil, nil, ErrNotFound
					}
					if !reg.Match([]byte(l)) {
						return nil, nil, ErrNotFound
					}
				}

				pn := getParamName(c)

				if params == nil {
					t.paramsPool.New = func() interface{} {
						// NOTE: It is better to set the maximum value of paramters to capacity.
						return &[]Param{}
					}
					params = t.getParams()
				}

				(*params) = append((*params), Param{
					key:   pn,
					value: l,
				})
				t.putParams(params)

				curNode = cc[c]
				isParamMatch = true
				if idx > 0 {
					// ex. foo/bar/baz → /bar/baz
					path = path[idx:]
				}
			}
		}
		if !isParamMatch {
			// no matching path was found.
			return nil, nil, ErrNotFound
		}
	}

	actions := curNode.actions[method]
	if actions == nil {
		// no matching handler and middlewares was found.
		return nil, nil, ErrMethodNotAllowed
	}
	if params == nil {
		return actions, nil, nil
	}
	return actions, *params, nil
}

// getPattern gets a pattern from a label.
// ex.
// :id[^\d+$] → ^\d+$
// :id        → (.+)
func getPattern(label string) string {
	leftI := strings.Index(label, leftPtnDelimiter)
	rightI := strings.Index(label, rightPtnDelimiter)

	// if label doesn't have any pattern, return wild card pattern as default.
	if leftI == -1 || rightI == -1 {
		return ""
	}

	return label[leftI+1 : rightI]
}

// getParamName gets a parameter from a label.
// ex.
// :id[^\d+$] → id
// :id        → id
func getParamName(label string) string {
	leftI := strings.Index(label, paramDelimiter)
	rightI := func(l string) int {
		var n int

		for i := 0; i < len(l); i++ {
			n = i
			if l[i:i+1] == leftPtnDelimiter {
				n = i
				break
			}
			if i == len(l)-1 {
				n = i + 1
				break
			}
		}

		return n
	}(label)

	return label[leftI+1 : rightI]
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

// removeTrailingSlash removes trailing slash from path.
func removeTrailingSlash(path string) string {
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	return path
}
