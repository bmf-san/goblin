package goblin

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewTree(t *testing.T) {
	actual := NewTree()
	expected := &Tree{
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

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

func TestSearch(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", fooHandler)
	tree.Insert(http.MethodGet, "/foo/", fooHandler)
	tree.Insert(http.MethodGet, "/foo/bar/", fooHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id/:name`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/:name`, fooHandler)

	tree.Insert(http.MethodPost, "/", fooHandler)
	tree.Insert(http.MethodPost, "/foo/", fooHandler)
	tree.Insert(http.MethodPost, "/foo/bar/", fooHandler)
	tree.Insert(http.MethodPost, `/foo/bar/:id`, fooHandler)
	tree.Insert(http.MethodPost, `/foo/bar/:id/:name`, fooHandler)
	tree.Insert(http.MethodPost, `/foo/:id[^\d+$]/`, fooHandler)
	tree.Insert(http.MethodPost, `/foo/:id[^\d+$]/:name`, fooHandler)

	type Item struct {
		method string
		path   string
	}

	cases := []struct {
		item     Item
		expected *Result
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/123/",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/123",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/123/john",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
					&Param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/123/john/",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
					&Param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/bar/",
			},
			expected: &Result{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/123/",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/bar/123",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/bar/123/john",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
					&Param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo/123/john/",
			},
			expected: &Result{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "123",
					},
					&Param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)

		if err != nil {
			t.Errorf("handler:%v params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual: %v expected: %v\n", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestGetPattern(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
	}{
		{
			actual:   getPattern(`:id[^\d+$]`),
			expected: `^\d+$`,
		},
		{
			actual:   getPattern(`:id[`),
			expected: ptnWildcard,
		},
		{
			actual:   getPattern(`:id]`),
			expected: ptnWildcard,
		},
		{
			actual:   getPattern(`:id`),
			expected: ptnWildcard,
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}

func TestGetParameter(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
	}{
		{
			actual:   getParameter(`:id[^\d+$]`),
			expected: "id",
		},
		{
			actual:   getParameter(`:id[`),
			expected: "id",
		},
		{
			actual:   getParameter(`:id]`),
			expected: "id]",
		},
		{
			actual:   getParameter(`:id`),
			expected: "id",
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
