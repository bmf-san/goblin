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
			http.MethodOptions: &Node{
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

func TestInsert(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	if err := tree.Insert(http.MethodGet, "/", fooHandler); err != nil {
		t.Errorf("err: %v\n", err)
	}

	if err := tree.Insert(http.MethodGet, "/foo", fooHandler); err != nil {
		t.Errorf("err: %v\n", err)
	}
}

// Item is a set of routing definition.
type Item struct {
	method string
	path   string
}

// Expected is a set of expected.
type Expected struct {
	handler http.HandlerFunc
	params  Params
}

func TestSearchAllMethod(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", fooHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodPost, "/", fooHandler)
	tree.Insert(http.MethodPost, "/foo", fooHandler)
	tree.Insert(http.MethodPut, "/", fooHandler)
	tree.Insert(http.MethodPut, "/foo", fooHandler)
	tree.Insert(http.MethodPatch, `/`, fooHandler)
	tree.Insert(http.MethodPatch, `/foo`, fooHandler)
	tree.Insert(http.MethodDelete, `/`, fooHandler)
	tree.Insert(http.MethodDelete, `/foo`, fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPut,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPut,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPatch,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodPatch,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodDelete,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodDelete,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestSearchWithoutRoot(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/foo/", fooHandler)
	tree.Insert(http.MethodGet, "/bar/", fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestSearchTrailingSlash(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", fooHandler)
	tree.Insert(http.MethodGet, "/foo/", fooHandler)
	tree.Insert(http.MethodGet, "/foo/bar/", fooHandler)
	tree.Insert(http.MethodGet, "/bar/", fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestSearchStatic(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", fooHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, "/foo/bar", fooHandler)
	tree.Insert(http.MethodGet, "/bar", fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestSearchParam(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, `/foo/:id`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id/:name`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id/:name/:date`, fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: Expected{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: Expected{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
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
				path:   "/foo/1/john/2020",
			},
			expected: Expected{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
					&Param{
						key:   "name",
						value: "john",
					},
					&Param{
						key:   "date",
						value: "2020",
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestSearchRegexp(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/:name[^\w+$]`, fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: Expected{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: Expected{
				handler: fooHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
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
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}
	}
}

func TestNegativeSearchRegexp(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/:name[^\w+$]`, fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/one",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/1/1",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/one/john",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
				t.Errorf("error:%v actual handler:%v actual params:%v expected:%v", err, actual.handler, actual.params, c.expected)
			}

			for i, param := range actual.params {
				if !reflect.DeepEqual(param, c.expected.params[i]) {
					t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
				}
			}
		}
	}
}

func TestPositiveAndNegativeSearchRandom(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", fooHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, "/foo/bar", fooHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id/:name`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]`, fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/:name`, fooHandler)

	cases := []struct {
		item     Item
		expected Expected
	}{
		{
			item: Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/123",
			},
			expected: Expected{
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
			expected: Expected{
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
			expected: Expected{
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
				path:   "/foo/123/john",
			},
			expected: Expected{
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
				path:   "/bar",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar/foo",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/bar/foo/123",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/id/name/any",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/bar/id/name/any/more",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/name",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/name/any",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
		{
			item: Item{
				method: http.MethodGet,
				path:   "/foo/123/name/any",
			},
			expected: Expected{
				handler: nil,
				params:  Params{},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
				t.Errorf("error:%v actual handler:%v actual params:%v expected:%v", err, actual.handler, actual.params, c.expected)
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

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
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
