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
	method map[string]*Node
}

// Node is a node of tree.
type Node struct {
	label       string
	handler     http.Handler
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
	pathDelimiter     = "/"
	paramDelimiter    = ":"
	leftPtnDelimiter  = "["
	rightPtnDelimiter = "]"
	ptnWildcard       = "(.+)"
)

// NewTree creates a new trie tree.
func NewTree() *Tree {
	return &Tree{
		method: map[string]*Node{
			http.MethodGet: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPost: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPut: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPatch: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodDelete: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodOptions: {
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
		},
	}
}

// Insert inserts a route definition to tree.
func (t *Tree) Insert(method string, path string, handler http.Handler, mws middlewares) error {
	curNode := t.method[method]

	if path == pathDelimiter {
		if len(curNode.label) != 0 && curNode.handler == nil {
			return errors.New("Root node already exists")
		}

		curNode.label = path
		curNode.handler = handler
		curNode.middlewares = mws

		return nil
	}

	for _, l := range deleteEmpty(strings.Split(path, pathDelimiter)) {
		if nextNode, ok := curNode.children[l]; ok {
			curNode = nextNode
		} else {
			curNode.children[l] = &Node{
				label:       l,
				handler:     handler,
				middlewares: mws,
				children:    make(map[string]*Node),
			}

			curNode = curNode.children[l]
		}
	}

	return nil
}

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

	n := t.method[method]

	if len(n.label) == 0 && len(n.children) == 0 {
		return nil, errors.New("tree is empty")
	}

	label := deleteEmpty(strings.Split(path, pathDelimiter))
	curNode := n

	for _, l := range label {
		if nextNode, ok := curNode.children[l]; ok {
			curNode = nextNode
		} else {
			// pattern matching priority depends on an order of routing definition
			// ex.
			// 1 /foo/:id
			// 2 /foo/:id[^\d+$]
			// 3 /foo/:id[^\D+$]
			// priority is 1, 2, 3
			if len(curNode.children) == 0 {
				return &Result{}, errors.New("handler is not registered")
			}

			count := 0
			cc := curNode.children
			for c := range cc {
				if string([]rune(c)[0]) == paramDelimiter {
					ptn := getPattern(c)

					reg, err := regC.Get(ptn)
					if err != nil {
						return nil, err
					}
					if reg.Match([]byte(l)) {
						param := getParameter(c)
						params = append(params, &Param{
							key:   param,
							value: l,
						})

						curNode = cc[c]
						count++
						break
					} else {
						return &Result{}, errors.New("param does not match")
					}
				}

				count++

				// If no match is found until the last loop.
				if count == len(cc) {
					return &Result{}, errors.New("handler is not registered")
				}
			}
		}
	}

	if curNode.handler == nil {
		return &Result{}, errors.New("handler is not registered")
	}

	return &Result{
		handler:     curNode.handler,
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

// getParameter gets a parameter from a label.
// ex.
// :id[^\d+$] → id
// :id        → id
func getParameter(label string) string {
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

// deleteEmpty removes an empty value in slice.
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
