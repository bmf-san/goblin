package goblin

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func TestNewTree(t *testing.T) {
	actual := newTree()
	expected := &tree{
		node: &node{
			label:    "/",
			action:   &action{},
			children: []*node{},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v expected: %v\n", actual, expected)
	}
}

func TestGetParamsAndPutParams(t *testing.T) {
	params := &Params{
		{
			key:   "id",
			value: "123",
		},
		{
			key:   "name",
			value: "john",
		},
	}

	tree := newTree()
	tree.paramsPool.New = func() interface{} {
		// NOTE: It is better to set the maximum value of paramters to capacity.
		return &Params{}
	}
	params = tree.getParams()
	tree.putParams(params)

	expectedParams := &Params{}
	actualParams := tree.getParams()
	if !reflect.DeepEqual(actualParams, expectedParams) {
		t.Errorf("actual:%v expected:%v", actualParams, expectedParams)
	}
}

// item is a set of routing definition.
type item struct {
	path string
}

// caseWithFailure is a struct for testWithFailure.
type caseWithFailure struct {
	hasError       bool
	item           *item
	expectedAction *action
	expectedParams Params
}

// insertItem is a struct for insert method.
type insertItem struct {
	path        string
	handler     http.Handler
	middlewares []middleware
}

// searchItem is a struct for search method.
type searchItem struct {
	path string
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
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				path: "/foo/bar",
			},
			expected: ErrNotFound,
		},
		{
			// no matching Param was found.
			insertItems: []insertItem{
				{
					path:        `/foo/:id[^\d+$]`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				path: "/foo/name",
			},
			expected: ErrNotFound,
		},
		{
			// no matching Param was found.
			insertItems: []insertItem{
				{
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				path: "/bar",
			},
			expected: ErrNotFound,
		},
		{
			// no matching handler and middlewares was found.
			insertItems: []insertItem{
				{
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				path: "/",
			},
			expected: ErrNotFound,
		},
		{
			// no matching handler and middlewares was found.
			insertItems: []insertItem{
				{
					path:        `/foo`,
					handler:     fooHandler,
					middlewares: []middleware{first},
				},
			},
			searchItem: &searchItem{
				path: "/fo",
			},
			expected: ErrNotFound,
		},
	}

	for _, c := range cases {
		tree := newTree()
		for _, i := range c.insertItems {
			tree.Insert(i.path, i.handler, i.middlewares)
		}
		actualAction, actualParams, err := tree.Search(c.searchItem.path)
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

	tree.Insert(`/`, rootHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "//",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo/bar",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWithoutRoot(t *testing.T) {
	tree := newTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	barHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(`/foo`, fooHandler, []middleware{first})
	tree.Insert(`/bar`, barHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/bar",
			},
			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchCommonPrefix(t *testing.T) {
	tree := newTree()

	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(`/foo`, fooHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/b",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/f",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/fo",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/fooo",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo/bar",
			},
			expectedAction: nil,
			expectedParams: Params{},
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

	tree.Insert(`/`, rootHandler, []middleware{first})
	tree.Insert(`/foo/`, fooHandler, []middleware{first})
	tree.Insert(`/bar/`, barHandler, []middleware{first})
	tree.Insert(`/foo/bar/`, fooBarHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "//",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},

			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/bar/",
			},

			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/bar",
			},

			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar/",
			},

			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar",
			},

			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
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

	tree.Insert(`/`, rootHandler, []middleware{first})
	tree.Insert(`/foo`, fooHandler, []middleware{first})
	tree.Insert(`/bar`, barHandler, []middleware{first})
	tree.Insert(`/foo/bar`, fooBarHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/bar",
			},
			expectedAction: &action{
				handler:     barHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/baz",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo/baz",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo/bar/baz",
			},
			expectedAction: nil,
			expectedParams: Params{},
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

	tree.Insert(`/:id`, idHandler, []middleware{first})
	tree.Insert(`/foo/:id`, fooIDHandler, []middleware{first})
	tree.Insert(`/foo/:id/:name`, fooIDNameHandler, []middleware{first})
	tree.Insert(`/foo/:id/:name/:date`, fooIDNameDateHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/1",
			},
			expectedAction: &action{
				handler:     idHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/1",
			},
			expectedAction: &action{
				handler:     fooIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/1/john",
			},
			expectedAction: &action{
				handler:     fooIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
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
				path: "/foo/1/john/2020",
			},
			expectedAction: &action{
				handler:     fooIDNameDateHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
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

	tree.Insert(`/`, rootHandler, []middleware{first})
	tree.Insert(`/`, rootPriorityHandler, []middleware{first})
	tree.Insert(`/foo`, fooHandler, []middleware{first})
	tree.Insert(`/foo`, fooPriorityHandler, []middleware{first})
	tree.Insert(`/:id`, IDHandler, []middleware{first})
	tree.Insert(`/:id`, IDPriorityHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},
			expectedAction: &action{
				handler:     fooPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/1",
			},
			expectedAction: &action{
				handler:     IDPriorityHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
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
	bazInvalidIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(`/`, rootHandler, []middleware{first})
	tree.Insert(`/:*[(.+)]`, rootWildCardHandler, []middleware{first})
	tree.Insert(`/foo`, fooHandler, []middleware{first})
	tree.Insert(`/foo/:id[^\d+$]`, fooIDHandler, []middleware{first})
	tree.Insert(`/foo/:id[^\d+$]/:name[^\D+$]`, fooIDNameHandler, []middleware{first})
	tree.Insert(`/foo/bar`, fooBarHandler, []middleware{first})
	tree.Insert(`/foo/bar/:id`, fooBarIDHandler, []middleware{first})
	tree.Insert(`/foo/bar/:id/:name`, fooBarIDNameHandler, []middleware{first})
	tree.Insert(`/baz/:id[[\d+]`, bazInvalidIDHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/wildcard",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "*",
					value: "wildcard",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/1234",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "*",
					value: "1234",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo",
			},
			expectedAction: &action{
				handler:     fooHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/bar",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "*",
					value: "bar",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/1",
			},
			expectedAction: &action{
				handler:     fooIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				path: "/foo/notnumber",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/1/john",
			},
			expectedAction: &action{
				handler:     fooIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
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
				path: "/foo/1/1",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar",
			},
			expectedAction: &action{
				handler:     fooBarHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar/1",
			},
			expectedAction: &action{
				handler:     fooBarIDHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "id",
					value: "1",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/foo/bar/1/john",
			},
			expectedAction: &action{
				handler:     fooBarIDNameHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
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
				path: "/baz/1",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
	}

	testWithFailure(t, tree, cases)
}

func TestSearchWildCardRegexp(t *testing.T) {
	tree := newTree()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rootWildCardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tree.Insert(`/`, rootHandler, []middleware{first})
	tree.Insert(`/:*[(.+)]`, rootWildCardHandler, []middleware{first})

	cases := []caseWithFailure{
		{
			hasError: false,
			item: &item{
				path: "/",
			},
			expectedAction: &action{
				handler:     rootHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{},
		},
		{
			hasError: false,
			item: &item{
				path: "/wildcard",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "*",
					value: "wildcard",
				},
			},
		},
		{
			hasError: false,
			item: &item{
				path: "/1234",
			},
			expectedAction: &action{
				handler:     rootWildCardHandler,
				middlewares: []middleware{first},
			},
			expectedParams: Params{
				{
					key:   "*",
					value: "1234",
				},
			},
		},
		{
			hasError: true,
			item: &item{
				path: "/1234/foo",
			},
			expectedAction: nil,
			expectedParams: Params{},
		},
	}

	testWithFailure(t, tree, cases)
}

func testWithFailure(t *testing.T, tree *tree, cases []caseWithFailure) {
	for _, c := range cases {
		actualAction, actualParams, err := tree.Search(c.item.path)

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

func TestGetReg(t *testing.T) {
	cases := []struct {
		name        string
		ptn         string
		isCached    bool
		expectedReg *regexp.Regexp
		expectedErr error
	}{
		{
			name:        "Valid - no cache",
			ptn:         `\d+`,
			isCached:    false,
			expectedReg: regexp.MustCompile(`\d+`),
			expectedErr: nil,
		},
		{
			name:        "Valid - cached",
			ptn:         `\d+`,
			isCached:    true,
			expectedReg: regexp.MustCompile(`\d+`),
			expectedErr: nil,
		},
		{
			name:        "Invalid - regexp compile error",
			ptn:         `[\d+`,
			isCached:    false,
			expectedReg: nil,
			expectedErr: fmt.Errorf("error parsing regexp: missing closing ]: `[\\d+`"),
		},
	}

	cache := regCache{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.isCached {
				cache.s.Store(c.ptn, c.expectedReg)
			}
			reg, err := cache.getReg(c.ptn)

			if !reflect.DeepEqual(reg, c.expectedReg) {
				t.Errorf("actual:%v expected:%v", reg, c.expectedReg)
			}

			if err != nil {
				if err.Error() != c.expectedErr.Error() {
					t.Errorf("actual:%v expected:%v", err.Error(), c.expectedErr.Error())
				}
			}
		})
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
