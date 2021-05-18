package goblin

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewResult(t *testing.T) {
	actual := newResult()
	expected := &result{}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

func TestNewTree(t *testing.T) {
	actual := NewTree()
	expected := &tree{
		node: &node{
			label:    pathRoot,
			actions:  make(map[string]*action),
			children: make(map[string]*node),
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

// item is a set of routing definition.
type item struct {
	method string
	path   string
}

// caseOnlySuccess is a struct for testOnlySuccess.
type caseOnlySuccess struct {
	item     *item
	expected *result
}

// caseWithFailure is a struct for testWithFailure.
type caseWithFailure struct {
	hasError bool
	item     *item
	expected *result
}

// insertItem is a struct for insert method.
type insertItem struct {
	methods     []string
	path        string
	handler     http.Handler
	middlewares []middleware
}

// searchItem is a struct for search method.
type searchItem struct {
	method string
	path   string
}

func TestSearchFailure(t *testing.T) {
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	cases := []struct {
		insertItems []*insertItem
		searchItem  *searchItem
		expected    error
	}{
		{
			// no matching path was found.
			insertItems: []*insertItem{
				&insertItem{
					methods:     []string{http.MethodGet},
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: ErrNotFound,
		},
		{
			// no matching param was found.
			insertItems: []*insertItem{
				&insertItem{
					methods:     []string{http.MethodGet},
					path:        `/foo/:id[^\d+$]`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				method: http.MethodGet,
				path:   "/foo/name",
			},
			expected: ErrNotFound,
		},
		{
			// no matching param was found.
			insertItems: []*insertItem{
				&insertItem{
					methods:     []string{http.MethodGet},
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: ErrNotFound,
		},
		{
			// no matching handler and middlewares was found.
			insertItems: []*insertItem{
				&insertItem{
					methods:     []string{http.MethodGet},
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				method: http.MethodGet,
				path:   "/",
			},
			expected: ErrNotFound,
		},
		{
			// no matching handler and middlewares was found.
			insertItems: []*insertItem{
				&insertItem{
					methods:     []string{http.MethodGet},
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: ErrMethodNotAllowed,
		},
	}

	for _, c := range cases {
		tree := NewTree()
		for _, i := range c.insertItems {
			tree.Insert(i.methods, i.path, i.handler, i.middlewares)
		}
		actual, err := tree.Search(c.searchItem.method, c.searchItem.path)
		if actual != nil {
			t.Fatalf("actual: %v expected err: %v", actual, err)
		}

		if err != c.expected {
			t.Fatalf("err: %v expected: %v\n", err, c.expected)
		}
	}
}

func TestSearchOnlyRoot(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "//",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: nil,
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: nil,
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWithoutRoot(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/bar`, barHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &result{
				actions: &action{
					handler:     barHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: nil,
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchCommonPrefix(t *testing.T) {
	tree := NewTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		// TODO: /fooで登録、/fo一致しちゃう問題
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/",
		// 	},
		// 	expected: nil,
		// },
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/b",
		// 	},
		// 	expected: nil,
		// },
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/f",
		// 	},
		// 	expected: nil,
		// },
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/fo",
		// 	},
		// 	expected: nil,
		// },
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/fooo",
		// 	},
		// 	expected: nil,
		// },
		// {
		// 	hasError: true,
		// 	item: &item{
		// 		method: http.MethodGet,
		// 		path:   "/foo/bar",
		// 	},
		// 	expected: nil,
		// },
	}

	testWithFailure(t, tree, cases)
}

func TestSearchAllMethod(t *testing.T) {
	tree := NewTree()

	rootGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootOptionsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooOptionsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootGetHandler, []middleware{first})
	tree.Insert([]string{http.MethodPost}, `/`, rootPostHandler, []middleware{first})
	tree.Insert([]string{http.MethodPut}, `/`, rootPutHandler, []middleware{first})
	tree.Insert([]string{http.MethodPatch}, `/`, rootPatchHandler, []middleware{first})
	tree.Insert([]string{http.MethodDelete}, `/`, rootDeleteHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/`, rootOptionsHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo`, fooGetHandler, []middleware{first})
	tree.Insert([]string{http.MethodPost}, `/foo`, fooPostHandler, []middleware{first})
	tree.Insert([]string{http.MethodPut}, `/foo`, fooPutHandler, []middleware{first})
	tree.Insert([]string{http.MethodPatch}, `/foo`, fooPatchHandler, []middleware{first})
	tree.Insert([]string{http.MethodDelete}, `/foo`, fooDeleteHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/foo`, fooOptionsHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootGetHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPostHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPut,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPutHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPatchHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootDeleteHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootOptionsHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooGetHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooPostHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPut,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooPutHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooPatchHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooDeleteHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooOptionsHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPathCommonMultiMethods(t *testing.T) {
	tree := NewTree()

	rootGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPatchHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooGetPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarGetPostDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootGetHandler, []middleware{first})
	tree.Insert([]string{http.MethodPost}, `/`, rootPostHandler, []middleware{second})
	tree.Insert([]string{http.MethodDelete}, `/`, rootDeleteHandler, []middleware{third})
	tree.Insert([]string{http.MethodPatch}, `/`, rootPatchHandler, []middleware{first, second, third})
	tree.Insert([]string{http.MethodGet, http.MethodPost}, `/foo`, fooGetPostHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, `/foo/bar`, fooBarGetPostDeleteHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootGetHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPostHandler,
					middlewares: []middleware{second},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootDeleteHandler,
					middlewares: []middleware{third},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPatchHandler,
					middlewares: []middleware{first, second, third},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooGetPostHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooGetPostHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarGetPostDeleteHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarGetPostDeleteHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarGetPostDeleteHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchTrailingSlash(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/`, fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/bar/`, barHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar/`, fooBarHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "//",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar/",
			},
			expected: &result{
				actions: &action{
					handler:     barHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &result{
				actions: &action{
					handler:     barHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchStaticPath(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooBarHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/bar`, barHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar`, fooBarHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: &result{
				actions: &action{
					handler:     barHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/baz",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/baz",
			},
			expected: nil,
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/baz",
			},
			expected: nil,
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPathWithParams(t *testing.T) {
	tree := NewTree()

	idHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooIDNameDateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/:id`, idHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id`, fooIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id/:name`, fooIDNameHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id/:name/:date`, fooIDNameDateHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: &result{
				actions: &action{
					handler:     idHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &result{
				actions: &action{
					handler:     fooIDHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &result{
				actions: &action{
					handler:     fooIDNameHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
					&param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john/2020",
			},
			expected: &result{
				actions: &action{
					handler:     fooIDNameDateHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
					&param{
						key:   "name",
						value: "john",
					},
					&param{
						key:   "date",
						value: "2020",
					},
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPriority(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fooPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	IDPriorityHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/`, rootPriorityHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo`, fooPriorityHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/:id`, IDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/:id`, IDPriorityHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootPriorityHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooPriorityHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: &result{
				actions: &action{
					handler:     IDPriorityHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
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

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/:*[(.+)]`, rootWildCardHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id[^\d+$]`, fooIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/:id[^\d+$]/:name[^\D+$]`, fooIDNameHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar`, fooBarHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar/:id`, fooBarIDHandler, []middleware{first})
	tree.Insert([]string{http.MethodGet}, `/foo/bar/:id/:name`, fooBarIDNameHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &result{
				actions: &action{
					handler:     rootWildCardHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "*",
						value: "wildcard",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &result{
				actions: &action{
					handler:     rootWildCardHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "*",
						value: "1234",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expected: &result{
				actions: &action{
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/bar",
			},
			expected: &result{
				actions: &action{
					handler:     rootWildCardHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "*",
						value: "bar",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expected: &result{
				actions: &action{
					handler:     fooIDHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/notnumber",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expected: &result{
				actions: &action{
					handler:     fooIDNameHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
					&param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/1",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/1",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarIDHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodPost,
				path:   "/foo/bar/1",
			},
			expected: nil,
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/1/john",
			},
			expected: &result{
				actions: &action{
					handler:     fooBarIDNameHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "id",
						value: "1",
					},
					&param{
						key:   "name",
						value: "john",
					},
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWildCardRegexp(t *testing.T) {
	tree := NewTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodOptions}, `/`, rootHandler, []middleware{first})
	tree.Insert([]string{http.MethodOptions}, `/:*[(.+)]`, rootWildCardHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/",
			},
			expected: &result{
				actions: &action{
					handler:     rootHandler,
					middlewares: []middleware{first},
				},
				params: params{},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expected: &result{
				actions: &action{
					handler:     rootWildCardHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "*",
						value: "wildcard",
					},
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expected: &result{
				actions: &action{
					handler:     rootWildCardHandler,
					middlewares: []middleware{first},
				},
				params: params{
					&param{
						key:   "*",
						value: "1234",
					},
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234/foo",
			},
			expected: nil,
		},
	}

	testWithFailure(t, tree, cases)
}

func testWithFailure(t *testing.T, tree *tree, cases []caseWithFailure) {
	for _, c := range cases {
		actual, err := tree.Search(c.item.method, c.item.path)

		if c.hasError {
			if err == nil {
				t.Fatalf("actual: %v expected err: %v", actual, err)
			}

			if actual != c.expected {
				t.Errorf("actual:%v expected:%v", actual, c.expected)
			}

			continue
		}

		if err != nil {
			t.Fatalf("err: %v actual: %v expected: %v\n", err, actual, c.expected)
		}

		if reflect.ValueOf(actual.actions.handler) != reflect.ValueOf(c.expected.actions.handler) {
			t.Errorf("actual:%v expected:%v", actual.actions.handler, c.expected.actions.handler)
		}

		if len(actual.actions.middlewares) != len(c.expected.actions.middlewares) {
			t.Errorf("actual: %v expected: %v\n", len(actual.actions.middlewares), len(c.expected.actions.middlewares))
		}

		for i, mws := range actual.actions.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expected.actions.middlewares[i]) {
				t.Errorf("actual: %v expected: %v\n", mws, c.expected.actions.middlewares[i])
			}
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
			actual:   explodePath(""),
			expected: nil,
		},
		{
			actual:   explodePath("/"),
			expected: nil,
		},
		{
			actual:   explodePath("//"),
			expected: nil,
		},
		{
			actual:   explodePath("///"),
			expected: nil,
		},
		{
			actual:   explodePath("/foo"),
			expected: []string{"foo"},
		},
		{
			actual:   explodePath("/foo/bar"),
			expected: []string{"foo", "bar"},
		},
		{
			actual:   explodePath("/foo/bar/baz"),
			expected: []string{"foo", "bar", "baz"},
		},
		{
			actual:   explodePath("/foo/bar/baz/"),
			expected: []string{"foo", "bar", "baz"},
		},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.actual, c.expected) {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
