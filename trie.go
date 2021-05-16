package goblin

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// Tree is a trie tree.
type Tree struct {
	node *Node
}

// Node is a node of tree.
type Node struct {
	label       string
	actions     map[string]http.Handler // key: method value: handler
	middlewares middlewares
	children    map[string]*Node
}

// Param is parameter.
type Param struct {
	key   string
	value string
}

// Params is parameters.
type Params []*Param

// Result is a search result.
type Result struct {
	handler     http.Handler
	params      Params
	middlewares middlewares
}

const (
	pathRoot          string = "/"
	pathDelimiter     string = "/"
	paramDelimiter    string = ":"
	leftPtnDelimiter  string = "["
	rightPtnDelimiter string = "]"
	ptnWildcard       string = "(.+)"
)

// NewTree creates a new trie tree.
func NewTree() *Tree {
	return &Tree{
		node: &Node{
			label:       pathRoot,
			actions:     make(map[string]http.Handler),
			middlewares: nil,
			children:    make(map[string]*Node),
		},
	}
}

// Insert inserts a route definition to tree.
func (t *Tree) Insert(methods []string, path string, handler http.Handler, mws middlewares) error {
	curNode := t.node

	// For root node.
	if path == pathRoot {
		curNode.label = path
		for _, method := range methods {
			curNode.actions[method] = handler
		}
		curNode.middlewares = mws
	}

	ep := explodePath(strings.Split(path, pathDelimiter))
	for i, l := range ep {
		nextNode, ok := curNode.children[l]
		if ok {
			curNode = nextNode
		}
		// Create a new node.
		if !ok {
			curNode.children[l] = &Node{
				label:       l,
				actions:     make(map[string]http.Handler),
				middlewares: nil,
				children:    make(map[string]*Node),
			}
			curNode = curNode.children[l]
		}
		// last loop.
		if i == len(ep)-1 {
			curNode.label = l
			for _, method := range methods {
				curNode.actions[method] = handler
			}
			curNode.middlewares = mws
		}
	}

	return nil
}

// regCache represents the cache for a regular expression.
type regCache struct {
	s sync.Map
}

// Get gets a compiled regexp from cache or create it.
func (rc *regCache) Get(ptn string) (*regexp.Regexp, error) {
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
func (t *Tree) Search(method string, path string) (*Result, error) {
	var params Params

	n := t.node

	label := explodePath(strings.Split(path, pathDelimiter))
	curNode := n

	for _, l := range label {
		if nextNode, ok := curNode.children[l]; ok {
			curNode = nextNode
		} else {
			cc := curNode.children
			for c := range cc {
				if string([]rune(c)[0]) == paramDelimiter {
					ptn := getPattern(c)

					reg, err := regC.Get(ptn)
					if err != nil {
						return nil, err
					}
					if reg.Match([]byte(l)) {
						param := getParamName(c)
						params = append(params, &Param{
							key:   param,
							value: l,
						})

						curNode = cc[c]
						break
					} else {
						return &Result{}, errors.New("param does not match")
					}
				}
			}
		}
	}

	handler := curNode.actions[method]
	if handler == nil {
		return &Result{}, errors.New("handler is not registered")
	}

	return &Result{
		handler:     curNode.actions[method],
		params:      params,
		middlewares: curNode.middlewares,
	}, nil
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
		return ptnWildcard
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
		r := []rune(l)

		var n int

		for i := 0; i < len(r); i++ {
			n = i
			if string(r[i]) == leftPtnDelimiter {
				n = i
				break
			} else if i == len(r)-1 {
				n = i + 1
				break
			}
		}

		return n
	}(label)

	return label[leftI+1 : rightI]
}

// explodePath removes an empty value in slice.
func explodePath(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
