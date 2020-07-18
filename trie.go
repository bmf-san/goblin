package goblin

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Tree is a trie tree.
type Tree struct {
	method map[string]*Node
}

// Node is a node of tree.
type Node struct {
	label    string
	handler  http.HandlerFunc
	children map[string]*Node
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
	handler http.HandlerFunc
	params  Params
}

const (
	pathDelimiter     = "/"
	paramDelimiter    = ":"
	leftPtnDelimiter  = "["
	rightPtnDelimiter = "]"
	ptnWildcard       = "(.+)"
)

// NewTree is create a new trie tree.
func NewTree() *Tree {
	return &Tree{
		method: map[string]*Node{
			http.MethodGet: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPost: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPut: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodPatch: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodDelete: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
			http.MethodOptions: &Node{
				label:    "",
				handler:  nil,
				children: make(map[string]*Node),
			},
		},
	}
}

// Insert insert a route definition to tree.
func (t *Tree) Insert(method string, path string, handler http.HandlerFunc) error {
	curNode := t.method[method]

	if path == pathDelimiter {
		if len(curNode.label) != 0 && curNode.handler == nil {
			return errors.New("Root node already exists")
		}

		curNode.label = path
		curNode.handler = handler

		return nil
	}

	for _, l := range deleteEmpty(strings.Split(path, pathDelimiter)) {
		if nextNode, ok := curNode.children[l]; ok {
			curNode = nextNode
		} else {
			curNode.children[l] = &Node{
				label:    l,
				handler:  handler,
				children: make(map[string]*Node),
			}

			curNode = curNode.children[l]
		}
	}

	return nil
}

// Search search a path from a tree.
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
			// 3 /foo/:id[^\w+$]
			// priority is 1, 2, 3
			if len(curNode.children) == 0 {
				return &Result{}, errors.New("handler is not regsitered")
			}

			count := 0
			for c := range curNode.children {
				if string([]rune(c)[0]) == paramDelimiter {
					ptn := getPattern(c)

					// HACK: regexp is slow so initialize a pattern as a global variable.
					if regexp.MustCompile(ptn).Match([]byte(l)) {
						param := getParameter(c)
						params = append(params, &Param{
							key:   param,
							value: l,
						})

						curNode = curNode.children[c]
						count++
						break
					} else {
						return &Result{}, errors.New("param does not match")
					}
				}

				count++

				// If no match is found until the last loop.
				if count == len(curNode.children) {
					return &Result{}, errors.New("handler is not regsitered")
				}
			}
		}
	}

	if curNode.handler == nil {
		return &Result{}, errors.New("handler is not registered")
	}

	return &Result{
		handler: curNode.handler,
		params:  params,
	}, nil
}

// getPattern get a pattern from a label.
// ex.
// :id[^\d+$] → ^\d+$
// :id        → *
func getPattern(label string) string {
	leftI := strings.Index(label, leftPtnDelimiter)
	rightI := strings.Index(label, rightPtnDelimiter)

	// if label has not pattern, return wild card pattern as default.
	if leftI == -1 || rightI == -1 {
		return ptnWildcard
	}

	return label[leftI+1 : rightI]
}

// getParameter get a parameter from a label.
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
