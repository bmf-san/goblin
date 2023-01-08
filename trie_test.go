package goblin

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewTree(t *testing.T) {
	actual := newTree()
	expected := &tree{
		node: &node{
			label:    "/",
			actions:  make(map[string]*action),
			children: make(map[string]*node),
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

// item is a set of routing definition.
type item struct {
	method string
	path   string
}

// caseWithFailure is a struct for testWithFailure.
type caseWithFailure struct {
	hasError       bool
	item           *item
	expectedAction *action
	expectedParams []Param
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
		insertItems []insertItem
		searchItem  *searchItem
		expected    error
	}{
		{
			// no matching path was found.
			insertItems: []insertItem{
				{
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
			// no matching Param was found.
			insertItems: []insertItem{
				{
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
			// no matching Param was found.
			insertItems: []insertItem{
				{
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
			insertItems: []insertItem{
				{
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
			insertItems: []insertItem{
				{
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
		tree := newTree()
		for _, i := range c.insertItems {
			tree.Insert(i.methods, i.path, i.handler, i.middlewares)
		}
		actualAction, actualParams, err := tree.Search(c.searchItem.method, c.searchItem.path)
		if actualAction != nil || actualParams != nil {
			t.Fatalf("actualAction: %v actualParams: %v expected err: %v", actualAction, actualParams, err)
		}

		if err != c.expected {
			t.Fatalf("err: %v expected: %v\n", err, c.expected)
		}
	}
}

func TestSearchOnlyRoot(t *testing.T) {
	tree := newTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/`, rootHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "//",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWithoutRoot(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchCommonPrefix(t *testing.T) {
	tree := newTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert([]string{http.MethodGet}, `/foo`, fooHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/b",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/f",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/fo",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/fooo",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchAllMethod(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootGetHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootPostHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPut,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootPutHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootPatchHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootDeleteHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootOptionsHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooGetHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooPostHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPut,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooPutHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooPatchHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooDeleteHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooOptionsHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPathCommonMultiMethods(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootGetHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootPostHandler,
				middlewares: []middleware{second},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootDeleteHandler,
				middlewares: []middleware{third},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPatch,
				path:   "/",
			},
			expectedAction: &action{
				handler:     rootPatchHandler,
				middlewares: []middleware{first, second, third},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooGetPostHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooGetPostHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarGetPostDeleteHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodPost,
				path:   "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarGetPostDeleteHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodDelete,
				path:   "/foo/bar",
			},

			expectedAction: &action{
				handler:     fooBarGetPostDeleteHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchTrailingSlash(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "//",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},

			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar/",
			},

			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},

			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/",
			},

			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},

			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchStaticPath(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/baz",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/baz",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/baz",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPathWithParams(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     idHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expectedAction: &action{
				handler:     fooIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expectedAction: &action{
				handler:     fooIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
				{
					key:   "name",
					value: "john",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john/2020",
			},
			expectedAction: &action{
				handler:     fooIDNameDateHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
				{
					key:   "name",
					value: "john",
				},
				{
					key:   "date",
					value: "2020",
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchPriority(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/1",
			},
			expectedAction: &action{
				handler:     IDPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchRegexp(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodPost,
				path:   "/",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "*",
					value: "wildcard",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "*",
					value: "1234",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/bar",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/bar",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "*",
					value: "bar",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1",
			},
			expectedAction: &action{
				handler:     fooIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/notnumber",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/john",
			},
			expectedAction: &action{
				handler:     fooIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
				{
					key:   "name",
					value: "john",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/1/1",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/1",
			},
			expectedAction: &action{
				handler:     fooBarIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodPost,
				path:   "/foo/bar/1",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodGet,
				path:   "/foo/bar/1/john",
			},
			expectedAction: &action{
				handler:     fooBarIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "id",
					value: "1",
				},
				{
					key:   "name",
					value: "john",
				},
			},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWildCardRegexp(t *testing.T) {
	tree := newTree()

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
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/wildcard",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "*",
					value: "wildcard",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: []Param{
				{
					key:   "*",
					value: "1234",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				method: http.MethodOptions,
				path:   "/1234/foo",
			},
			expectedAction: nil,
			expectedParams: []Param{},
		},
	}

	testWithFailure(t, tree, cases)
}

func testWithFailure(t *testing.T, tree *tree, cases []caseWithFailure) {
	for _, c := range cases {
		actualAction, actualParams, err := tree.Search(c.item.method, c.item.path)

		if c.hasError {
			if err == nil {
				t.Fatalf("actualAction: %v actualParams: %v expected err: %v", actualAction, actualParams, err)
			}

			if actualAction != c.expectedAction {
				t.Errorf("actualAction:%v expectedAction:%v", actualAction, c.expectedAction)
			}

			if len(actualParams) != len(c.expectedParams) {
				t.Errorf("actualParams: %v expectedParams: %v\n", len(actualParams), len(c.expectedParams))
			}

			for i, Param := range actualParams {
				if !reflect.DeepEqual(Param, c.expectedParams[i]) {
					t.Errorf("actualParams: %v expectedParams: %v\n", Param, c.expectedParams[i])
				}
			}

			continue
		}

		if err != nil {
			t.Fatalf("actualAction: %v actualParams: %v expected err: %v", actualAction, actualParams, err)
		}

		if reflect.ValueOf(actualAction.handler) != reflect.ValueOf(c.expectedAction.handler) {
			t.Errorf("actualActionHandler:%v expectedActionHandler:%v", actualAction.handler, c.expectedAction.handler)
		}

		if len(actualAction.middlewares) != len(c.expectedAction.middlewares) {
			t.Errorf("actualActionMiddlewares: %v expectedActionsMiddleware: %v\n", len(actualAction.middlewares), len(c.expectedAction.middlewares))
		}

		for i, mws := range actualAction.middlewares {
			if reflect.ValueOf(mws) != reflect.ValueOf(c.expectedAction.middlewares[i]) {
				t.Errorf("actualActionsMiddleware: %v expectedActionsMiddleware: %v\n", mws, c.expectedAction.middlewares[i])
			}
		}

		if len(actualParams) != len(c.expectedParams) {
			t.Errorf("actualParams: %v expectedParams: %v\n", len(actualParams), len(c.expectedParams))
		}

		for i, Param := range actualParams {
			if !reflect.DeepEqual(Param, c.expectedParams[i]) {
				t.Errorf("actualParam: %v expectedParam: %v\n", Param, c.expectedParams[i])
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
			expected: "",
		},
		{
			actual:   getPattern(`:id]`),
			expected: "",
		},
		{
			actual:   getPattern(`:id`),
			expected: "",
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

func TestCleanPath(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{
			path:     "",
			expected: "/",
		},
		{
			path:     "//",
			expected: "/",
		},
		{
			path:     "///",
			expected: "/",
		},
		{
			path:     "path",
			expected: "/path",
		},
		{
			path:     "/",
			expected: "/",
		},
		{
			path:     "/path/trailingslash/",
			expected: "/path/trailingslash/",
		},
		{
			path:     "path/trailingslash//",
			expected: "/path/trailingslash/",
		},
	}

	for _, c := range cases {
		actual := cleanPath(c.path)
		if actual != c.expected {
			t.Errorf("actual: %v expected: %v\n", actual, c.expected)
		}
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{
			path:     "//",
			expected: "/",
		},
		{
			path:     "///",
			expected: "//",
		},
		{
			path:     "path",
			expected: "path",
		},
		{
			path:     "/",
			expected: "",
		},
		{
			path:     "/path/trailingslash/",
			expected: "/path/trailingslash",
		},
		{
			path:     "/path/trailingslash//",
			expected: "/path/trailingslash/",
		},
	}

	for _, c := range cases {
		actual := removeTrailingSlash(c.path)
		if actual != c.expected {
			t.Errorf("actual: %v expected: %v\n", actual, c.expected)
		}
	}
}
