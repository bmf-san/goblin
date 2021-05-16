package goblin

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewTree(t *testing.T) {
	actual := NewTree()
	expected := &Tree{
		node: &Node{
			label:       pathRoot,
			actions:     make(map[string]http.Handler),
			middlewares: nil,
			children:    make(map[string]*Node),
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

func TestInsert(t *testing.T) {
	tree := NewTree()
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	cases := []struct {
		method      string
		path        string
		handler     http.Handler
		middlewares middlewares
	}{
		{
			method:      http.MethodGet,
			path:        "/",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodPost,
			path:        "/",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodGet,
			path:        "/foo",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodPost,
			path:        "/foo",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodGet,
			path:        "/foo/bar",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodPost,
			path:        "/foo/bar",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodGet,
			path:        "/foo/bar/baz",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
		{
			method:      http.MethodPost,
			path:        "/foo/bar/baz",
			handler:     fooHandler,
			middlewares: []middleware{first, second, third},
		},
	}

	for _, c := range cases {
		if err := tree.Insert([]string{c.method}, c.path, c.handler, c.middlewares); err != nil {
			t.Errorf("err: %v\n", err)
		}
	}
}

// Item is a set of routing definition.
type Item struct {
	method string
	path   string
}

func TestSearchAllMethod(t *testing.T) {
	tree := NewTree()

	rootGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootGetHandler, []middleware{first})
	tree.Insert([]string{http.MethodPost}, "/", rootPostHandler, []middleware{first})
	tree.Insert([]string{http.MethodPut}, "/", rootPutHandler, []middleware{first})
	tree.Insert([]string{http.MethodPatch}, `/`, rootPatchHandler, []middleware{first})
	tree.Insert([]string{http.MethodDelete}, `/`, rootDeleteHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo", fooGetHandler, []middleware{first})
	tree.Insert([]string{http.MethodPost}, "/foo", fooPostHandler, []middleware{first})
	tree.Insert([]string{http.MethodPut}, "/foo", fooPutHandler, []middleware{first})
	tree.Insert([]string{http.MethodPatch}, `/foo`, fooPatchHandler, []middleware{first})
	tree.Insert([]string{http.MethodDelete}, `/foo`, fooDeleteHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootGetHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &Result{
				handler:     rootPostHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPut,
				path:   "/",
			},
			expected: &Result{
				handler:     rootPutHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPatch,
				path:   "/",
			},
			expected: &Result{
				handler:     rootPatchHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodDelete,
				path:   "/",
			},
			expected: &Result{
				handler:     rootDeleteHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooGetHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooPostHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPut,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooPutHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPatch,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooPatchHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodDelete,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooDeleteHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchPathCommon(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet, http.MethodPost}, "/foo", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, "/foo/bar", fooBarHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodDelete,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}
func TestSearchWithoutRoot(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/foo", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/bar", barHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Result{
				handler:     barHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchTrailingSlash(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo/", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/bar/", barHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo/bar/", fooBarHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar/",
			},
			expected: &Result{
				handler:     barHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Result{
				handler:     barHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchStaticPath(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/bar", barHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo/bar", fooBarHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Result{
				handler:     barHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchPathWithParams(t *testing.T) {
	tree := NewTree()

	fooIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameDateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/foo/:id`, fooIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id/:name`, fooIDNameHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id/:name/:date`, fooIDNameDateHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &Result{
				handler: fooIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &Result{
				handler: fooIDNameHandler,
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
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john/2020",
			},
			expected: &Result{
				handler: fooIDNameDateHandler,
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
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchPriority(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/", rootPriorityHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo", fooPriorityHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/:id", IDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/:id", IDPriorityHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootPriorityHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooPriorityHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: &Result{
				handler: IDPriorityHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)

		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}
	}
}

func TestSearchRegexp(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, "/", rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/:*[(.+)]`, rootWildCardHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo", fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id[^\d+$]`, fooIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id[^\d+$]/:name[^\D+$]`, fooIDNameHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, "/foo/bar", fooBarHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar/:id`, fooBarIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar/:id/:name`, fooBarIDNameHandler, []middleware{first})

	cases := []struct {
		hasError bool
		item     *Item
		expected *Result
	}{
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Result{
				handler:     rootHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &Result{
				handler:     nil,
				params:      Params{},
				middlewares: []middleware{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &Result{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "wildcard",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &Result{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "1234",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Result{
				handler:     fooHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Result{
				handler:     nil,
				params:      Params{},
				middlewares: []middleware{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodOptions,
				path:   "/bar",
			},
			expected: &Result{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "bar",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &Result{
				handler: fooIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/notnumber",
			},
			expected: &Result{
				handler:     nil,
				params:      Params{},
				middlewares: []middleware{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &Result{
				handler: fooIDNameHandler,
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
				middlewares: []middleware{first},
			},
		},
		{
			hasError: true,

			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/1",
			},
			expected: &Result{
				handler:     nil,
				params:      Params{},
				middlewares: []middleware{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Result{
				handler:     fooBarHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/1",
			},
			expected: &Result{
				handler: fooBarIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodPost,
				path:   "/foo/bar/1",
			},
			expected: &Result{
				handler:     nil,
				params:      Params{},
				middlewares: []middleware{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/1/john",
			},
			expected: &Result{
				handler: fooBarIDNameHandler,
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
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)

		if c.hasError {
			if err == nil {
				t.Errorf("expected err: %v actual: %v", err, actual)
			}

			if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
				t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
			}

			if len(actual.params) != len(c.expected.params) {
				t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
			}

			for i, param := range actual.params {
				if !reflect.DeepEqual(param, c.expected.params[i]) {
					t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
				}
			}

			if len(actual.middlewares) != len(c.expected.middlewares) {
				t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
			}

			for i, mws := range actual.middlewares {
				if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
					t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
				}
			}

			continue
		}

		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		if len(actual.params) != len(c.expected.params) {
			t.Errorf("actual: %v expected: %v\n", len(actual.params), len(c.expected.params))
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		if len(actual.middlewares) != len(c.expected.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.middlewares), len(c.expected.middlewares))
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
			}
		}

	}
}

func TestSearchWildCardRegexp(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodOptions}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/:*[(.+)]`, rootWildCardHandler, []middleware{first})

	cases := []struct {
		item     *Item
		expected *Result
	}{
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/",
			},
			expected: &Result{
				handler:     rootHandler,
				params:      Params{},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &Result{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "wildcard",
					},
				},
				middlewares: []middleware{first},
			},
		},
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &Result{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "1234",
					},
				},
				middlewares: []middleware{first},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if err != nil {
			t.Errorf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.handler) != reflect.ValueOf(c.expected.handler) {
			t.Errorf("actual:%v expected:%v", actual.handler, c.expected.handler)
		}

		for i, param := range actual.params {
			if !reflect.DeepEqual(param, c.expected.params[i]) {
				t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
			}
		}

		for i, mws := range actual.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.middlewares[i])
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

func TestGetParamName(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
	}{
		{
			actual:   getParamName(`:id[^\d+$]`),
			expected: "id",
		},
		{
			actual:   getParamName(`:id[`),
			expected: "id",
		},
		{
			actual:   getParamName(`:id]`),
			expected: "id]",
		},
		{
			actual:   getParamName(`:id`),
			expected: "id",
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}

func TestExplodePath(t *testing.T) {
	cases := []struct {
		actual   []string
		expected []string
	}{
		{
			actual:   explodePath(strings.Split("/", pathDelimiter)),
			expected: nil,
		},
		{
			actual:   explodePath(strings.Split("/foo", pathDelimiter)),
			expected: []string{"foo"},
		},
		{
			actual:   explodePath(strings.Split("/foo/bar", pathDelimiter)),
			expected: []string{"foo", "bar"},
		},
		{
			actual:   explodePath(strings.Split("/foo/bar/baz", pathDelimiter)),
			expected: []string{"foo", "bar", "baz"},
		},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
