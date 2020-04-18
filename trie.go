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
type Params []Param

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
		},
	}
}

// Insert insert a route definition to tree.
func (t *Tree) Insert(method string, path string, handler http.HandlerFunc) error {
	curNode := t.method[method]

	if path == "/" {
		if len(curNode.label) != 0 && curNode.handler == nil {
			return errors.New("Root node already exists")
		}

		curNode.label = path
		curNode.handler = handler

		return nil
	}

	for _, l := range deleteEmpty(strings.Split(path, "/")) {
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
func (t *Tree) Search(method string, path string) (http.HandlerFunc, *Params, error) {
	var params Params

	n := t.method[method]

	if len(n.label) == 0 && len(n.children) == 0 {
		return nil, nil, errors.New("tree is empty")
	}

	label := deleteEmpty(strings.Split(path, "/"))
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
			for c := range curNode.children {
				if string([]rune(c)[0]) == ":" {
					ptn := getPattern(c)

					// HACK: regexp is slow so initialize a pattern as a global variable.
					if regexp.MustCompile(ptn).Match([]byte(l)) {
						param := getParameter(c)
						params = append(params, Param{
							key:   param,
							value: l,
						})

						curNode = curNode.children[c]
						break
					} else {
						return nil, nil, errors.New("param does not match")
					}
				}
			}
		}
	}

	if curNode.handler == nil {
		return nil, nil, errors.New("handler is not registered")
	}

	return curNode.handler, &params, nil
}

// wildcard pattern.
const ptnWildcard = `(.+)`

// getPattern get a pattern from a label.
// ex.
// :id[^\d+$] → ^\d+$
// :id        → *
func getPattern(label string) string {
	startI := strings.Index(label, "[")
	endI := strings.Index(label, "]")

	// if label has not pattern, return wild card pattern as default.
	if startI == -1 || endI == -1 {
		return ptnWildcard
	}

	return label[startI+1 : endI]
}

// getParameter get a parameter from a label.
// ex.
// :id[^\d+$] → id
// :id        → id
func getParameter(label string) string {
	startI := strings.Index(label, ":")
	endI := func(l string) int {
		r := []rune(l)

		var n int

		for i := 0; i < len(r); i++ {
			n = i
			if string(r[i]) == "[" {
				n = i
				break
			} else if i == len(r)-1 {
				n = i + 1
				break
			}
		}

		return n
	}(label)

	return label[startI+1 : endI]
}
