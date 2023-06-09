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
	maxParams  int
}

// node is a node of tree.
type node struct {
	label    string
	action   *action // key is method
	children []*node // key is label of next nodes
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
			action:   &action{},
			children: []*node{},
		},
	}
}

// Param is a parameter.
type Param struct {
	key   string
	value string
}

// getParams gets parameters.
func (t *tree) getParams() *[]Param {
	ps, _ := t.paramsPool.Get().(*[]Param)
	*ps = (*ps)[:0] // reset slice
	return ps
}

// putParams puts parameters.
func (t *tree) putParams(ps *[]Param) {
	if ps != nil {
		t.paramsPool.Put(ps)
	}
}

func (n *node) getChild(label string) *node {
	for i := 0; i < len(n.children); i++ {
		if n.children[i].label == label {
			return n.children[i]
		}
	}
	return nil
}

// Insert inserts a route definition to tree.
func (t *tree) Insert(path string, handler http.Handler, mws middlewares) {
	path = cleanPath(path)
	curNode := t.node

	if path == "/" {
		curNode.label = path
		curNode.action = &action{
			middlewares: mws,
			handler:     handler,
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
		l = path
		if idx > 0 {
			// ex. foo/bar/baz → foo
			l = path[:idx]
		}

		nextNode := curNode.getChild(l)
		if nextNode != nil {
			curNode = nextNode
			if idx > 0 {
				l = path[:idx]
				// foo/bar/baz → /bar/baz
				path = path[idx:]
			}
		} else {
			// Create a new node.
			child := &node{
				label:    l,
				action:   &action{},
				children: []*node{},
			}
			curNode.children = append(curNode.children, child)
			curNode = child
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
			curNode.action = &action{
				middlewares: mws,
				handler:     handler,
			}
			break
		}
	}
	if t.maxParams < cnt {
		t.maxParams = cnt
	}
	if t.paramsPool.New == nil && t.maxParams > 0 {
		t.paramsPool.New = func() interface{} {
			p := make([]Param, 0, t.maxParams)
			return &p
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
func (t *tree) Search(path string) (*action, []Param, error) {
	path = cleanPath(path)
	curNode := t.node

	if path == "/" && curNode.action == nil {
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
		// ex. foo → foo
		l = path
		if idx > 0 {
			// ex. foo/bar/baz → foo
			l = path[:idx]
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

		nextNode := curNode.getChild(l)
		if nextNode != nil {
			curNode = nextNode
			if idx > 0 {
				// foo/bar/baz → /bar/baz
				path = path[idx:]
			}
			continue
		}

		isParamMatch := false
		// parameter matching
		for _, c := range cc {
			if c.label[0:1] == paramDelimiter {
				ptn := getPattern(c.label)
				if ptn != "" {
					reg, err := regC.getReg(ptn)
					if err != nil {
						return nil, nil, ErrNotFound
					}
					if !reg.Match([]byte(l)) {
						return nil, nil, ErrNotFound
					}
				}

				pn := getParamName(c.label)

				if params == nil {
					params = t.getParams()
				}
				lp := len(*params)
				*params = (*params)[:lp+1]
				(*params)[lp] = Param{
					key:   pn,
					value: l,
				}
				t.putParams(params)

				curNode = c
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

	action := curNode.action
	if action.handler == nil {
		// no matching handler and middlewares was found.
		return nil, nil, ErrNotFound
	}
	if params == nil {
		return action, nil, nil
	}
	return action, *params, nil
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
