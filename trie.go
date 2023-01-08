package goblin

import (
	"fmt"
	"net/http"
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
	ep := explodePath(path)
	for i := 0; i < len(ep); i++ {
		nextNode, ok := curNode.children[ep[i]]
		if ok {
			curNode = nextNode
		}
		// Create a new node.
		if !ok {
			curNode.children[ep[i]] = &node{
				label:    ep[i],
				actions:  make(map[string]*action),
				children: make(map[string]*node),
			}
			curNode = curNode.children[ep[i]]
		}
		// last loop.
		// If there is already registered data, overwrite it.
		if i == len(ep)-1 {
			curNode.label = ep[i]
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
// path is assumed to be formatted by the cleanPath function.
func (t *tree) Search(method string, path string) (*action, []Param, error) {
	// t.paramsPool is a pool for parameters.
	var params *[]Param
	var matchErr error

	curNode := t.node
	if path == "/" {
		if len(curNode.actions) == 0 {
			// no matching handler and middlewares was found.
			matchErr = ErrNotFound
		}
	}
	// NOTE: use the path string as is without using explodepath.
	// NOTE: The number of /'s in the path is the loop count.
	ep := explodePath(path)
	for i := 0; i < len(ep); i++ {
		nextNode, ok := curNode.children[ep[i]]
		if ok {
			curNode = nextNode
			continue
		}
		cc := curNode.children
		// leaf node
		if len(cc) == 0 {
			if curNode.label != ep[i] {
				// no matching path was found.
				matchErr = ErrNotFound
				break
			}
			break
		}
		isParamMatch := false
		// parameter matching
		for c := range cc {
			if c[0:1] == paramDelimiter {
				ptn := getPattern(c)
				reg, err := regC.Get(ptn)
				if err != nil {
					matchErr = ErrNotFound
					break
				}
				if reg.Match([]byte(ep[i])) {
					pn := getParamName(c)

					if params == nil {
						t.paramsPool.New = func() interface{} {
							// NOTE: It is better to set the maximum value of paramters to capacity.
							// NOTE: The parameter should be cached in a variable and set in the pool when returning.
							return &[]Param{}
						}
						params = t.getParams()
					}

					(*params) = append((*params), Param{
						key:   pn,
						value: ep[i],
					})
					t.putParams(params)

					curNode = cc[c]
					isParamMatch = true
					break
				}
				// no matching path was found.
				matchErr = ErrNotFound
			}
		}
		if !isParamMatch {
			// no matching path was found.
			matchErr = ErrNotFound
		}
	}

	if matchErr != nil {
		return nil, nil, matchErr
	}
	actions := curNode.actions[method]
	if actions == nil {
		// no matching handler and middlewares was found.
		matchErr = ErrMethodNotAllowed
		if matchErr != nil {
			return nil, nil, ErrMethodNotAllowed
		}
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

// explodePath converts a path to a slice split　by path delimiter.
// path expects a path processed by cleanPath.
func explodePath(path string) []string {
	splitFn := func(c rune) bool {
		return c == '/'
	}
	return strings.FieldsFunc(path, splitFn)
}
