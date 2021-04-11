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
	hasError bool
	handler  http.HandlerFunc
	params   Params
}

func TestSearchAllMethod(t *testing.T) {
	tree := NewTree()

	rootGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", rootGetHandler)
	tree.Insert(http.MethodGet, "/foo", fooGetHandler)
	tree.Insert(http.MethodPost, "/", rootPostHandler)
	tree.Insert(http.MethodPost, "/foo", fooPostHandler)
	tree.Insert(http.MethodPut, "/", rootPutHandler)
	tree.Insert(http.MethodPut, "/foo", fooPutHandler)
	tree.Insert(http.MethodPatch, `/`, rootPatchHandler)
	tree.Insert(http.MethodPatch, `/foo`, fooPatchHandler)
	tree.Insert(http.MethodDelete, `/`, rootDeleteHandler)
	tree.Insert(http.MethodDelete, `/foo`, fooDeleteHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootGetHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooGetHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &Expected{
				handler: rootPostHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooPostHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPut,
				path:   "/",
			},
			expected: &Expected{
				handler: rootPutHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPut,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooPutHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPatch,
				path:   "/",
			},
			expected: &Expected{
				handler: rootPatchHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodPatch,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooPatchHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodDelete,
				path:   "/",
			},
			expected: &Expected{
				handler: rootDeleteHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodDelete,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooDeleteHandler,
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
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, "/bar", barHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Expected{
				handler: barHandler,
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

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", rootHandler)
	tree.Insert(http.MethodGet, "/foo/", fooHandler)
	tree.Insert(http.MethodGet, "/bar/", barHandler)
	tree.Insert(http.MethodGet, "/foo/bar/", fooBarHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar/",
			},
			expected: &Expected{
				handler: barHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Expected{
				handler: barHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},
			expected: &Expected{
				handler: fooBarHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Expected{
				handler: fooBarHandler,
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

func TestSearchStaticPath(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", rootHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, "/bar", barHandler)
	tree.Insert(http.MethodGet, "/foo/bar", fooBarHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Expected{
				handler: barHandler,
				params:  Params{},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Expected{
				handler: fooBarHandler,
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

func TestSearchPathWithParams(t *testing.T) {
	tree := NewTree()

	fooIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameDateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, `/foo/:id`, fooIDHandler)
	tree.Insert(http.MethodGet, `/foo/:id/:name`, fooIDNameHandler)
	tree.Insert(http.MethodGet, `/foo/:id/:name/:date`, fooIDNameDateHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &Expected{
				handler: fooIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &Expected{
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
			},
		},
		{
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john/2020",
			},
			expected: &Expected{
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

func TestSearchPriority(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", rootHandler)
	tree.Insert(http.MethodGet, "/", rootPriorityHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, "/foo", fooPriorityHandler)
	tree.Insert(http.MethodGet, "/:id", IDHandler)
	tree.Insert(http.MethodGet, "/:id", IDPriorityHandler)

	cases := []struct {
		hasError bool
		item     *Item
		expected *Expected
	}{
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootHandler,
				params:  Params{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootPriorityHandler,
				params:  Params{},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooPriorityHandler,
				params:  Params{},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &Expected{
				handler: IDHandler,
				params:  Params{},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: &Expected{
				handler: IDPriorityHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if c.hasError {
			if reflect.ValueOf(actual.handler) == reflect.ValueOf(c.expected.handler) {
				t.Errorf("actual handler:%v actual params:%v expected:%v", actual.handler, actual.params, c.expected)
			}

			for i, param := range actual.params {
				if !reflect.DeepEqual(param, c.expected.params[i]) {
					t.Errorf("actual: %v expected: %v\n", param, c.expected.params[i])
				}
			}

			return
		}

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

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodGet, "/", rootHandler)
	tree.Insert(http.MethodOptions, `/:*[(.+)]`, rootWildCardHandler)
	tree.Insert(http.MethodGet, "/foo", fooHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]`, fooIDHandler)
	tree.Insert(http.MethodGet, `/foo/:id[^\d+$]/:name[^\D+$]`, fooIDNameHandler)
	tree.Insert(http.MethodGet, "/foo/bar", fooBarHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id`, fooBarIDHandler)
	tree.Insert(http.MethodGet, `/foo/bar/:id/:name`, fooBarIDNameHandler)

	cases := []struct {
		hasError bool
		item     *Item
		expected *Expected
	}{
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &Expected{
				handler: rootHandler,
				params:  Params{},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &Expected{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "wildcard",
					},
				},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &Expected{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "1234",
					},
				},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &Expected{
				handler: fooHandler,
				params:  Params{},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &Expected{
				handler: fooIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/notnumber",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &Expected{
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
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/1/1",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &Expected{
				handler: fooBarHandler,
				params:  Params{},
			},
		},
		{
			hasError: true,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/foo",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/1",
			},
			expected: &Expected{
				handler: fooBarIDHandler,
				params: Params{
					&Param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: false,
			item: &Item{
				method: http.MethodGet,
				path:   "/foo/bar/1/john",
			},
			expected: &Expected{
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
			},
		},
	}

	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)
		if c.hasError {
			if err == nil {
				t.Errorf("err: expected err actual: %v", actual)
			}
			return
		}

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

func TestSearchCORS(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(http.MethodOptions, `/`, rootHandler)
	tree.Insert(http.MethodOptions, `/:*[(.+)]`, rootWildCardHandler)

	cases := []struct {
		item     *Item
		expected *Expected
	}{
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/",
			},
			expected: &Expected{
				handler: rootHandler,
				params:  Params{nil},
			},
		},
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &Expected{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "wildcard",
					},
				},
			},
		},
		{
			item: &Item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &Expected{
				handler: rootWildCardHandler,
				params: Params{
					&Param{
						key:   "*",
						value: "1234",
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
