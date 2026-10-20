package goblin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestNewRouter(t *testing.T) {
	actual := NewRouter()
	expected := &Router{
		tree:    map[string]*tree{},
		pending: &route{},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v\n", actual, expected)
	}
}

func TestRouterIndependentRegistration(t *testing.T) {
	for _, firstFinishesFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("firstFinishesFirst=%t", firstFinishesFirst), func(t *testing.T) {
			r1 := NewRouter()
			r1.Methods(http.MethodGet).Use(first).Methods(http.MethodPatch)
			r2 := NewRouter()
			r2.Use(second).Methods(http.MethodPost)

			h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "one\n")
			})
			h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "two\n")
			})
			if firstFinishesFirst {
				r1.Handler("/route", h1)
				r2.Handler("/route", h2)
			} else {
				r2.Handler("/route", h2)
				r1.Handler("/route", h1)
			}

			cases := []struct {
				name   string
				router *Router
				method string
				code   int
				body   string
			}{
				{"first GET", r1, http.MethodGet, http.StatusOK, "first: before\none\nfirst: after\n"},
				{"first PATCH", r1, http.MethodPatch, http.StatusOK, "first: before\none\nfirst: after\n"},
				{"first POST", r1, http.MethodPost, http.StatusMethodNotAllowed, ""},
				{"second POST", r2, http.MethodPost, http.StatusOK, "second: before\ntwo\nsecond: after\n"},
				{"second GET", r2, http.MethodGet, http.StatusMethodNotAllowed, ""},
				{"second PATCH", r2, http.MethodPatch, http.StatusMethodNotAllowed, ""},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					rec := httptest.NewRecorder()
					c.router.ServeHTTP(rec, httptest.NewRequest(c.method, "/route", nil))
					if rec.Code != c.code || rec.Body.String() != c.body {
						t.Errorf("got (%d, %q), want (%d, %q)", rec.Code, rec.Body.String(), c.code, c.body)
					}
				})
			}

			r1.Methods(http.MethodDelete).Handler("/next", h1)
			r2.Methods(http.MethodPut).Handler("/next", h2)
			for _, c := range []struct {
				name   string
				router *Router
				method string
				code   int
				body   string
			}{
				{"first reset", r1, http.MethodDelete, http.StatusOK, "one\n"},
				{"second reset", r2, http.MethodPut, http.StatusOK, "two\n"},
				{"first old method", r1, http.MethodGet, http.StatusNotFound, "404 page not found\n"},
				{"second old method", r2, http.MethodPost, http.StatusNotFound, "404 page not found\n"},
			} {
				t.Run(c.name, func(t *testing.T) {
					rec := httptest.NewRecorder()
					c.router.ServeHTTP(rec, httptest.NewRequest(c.method, "/next", nil))
					if rec.Code != c.code || rec.Body.String() != c.body {
						t.Errorf("got (%d, %q), want (%d, %q)", rec.Code, rec.Body.String(), c.code, c.body)
					}
				})
			}
		})
	}
}

func TestRouterParallelRegistration(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		wg.Add(1)
		go func(method string) {
			defer wg.Done()
			r := NewRouter()
			<-start
			r.Methods(method).Use(first).Handler("/route", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, method+"\n")
			}))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(method, "/route", nil))
			body := "first: before\n" + method + "\nfirst: after\n"
			if rec.Code != http.StatusOK || rec.Body.String() != body {
				t.Errorf("%s got (%d, %q), want (200, %q)", method, rec.Code, rec.Body.String(), body)
			}
		}(method)
	}
	close(start)
	wg.Wait()
}

type routerTest struct {
	path   string
	method string
	code   int
	body   string
}

func (tc routerTest) name() string {
	return fmt.Sprintf("%s_%s_%d", tc.method, tc.path, tc.code)
}

func TestRouterMiddleware(t *testing.T) {
	r := NewRouter()

	r.UseGlobal(global)
	r.Methods(http.MethodGet).Handler(`/globalmiddleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/globalmiddleware\n")
	}))
	r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middleware\n")
	}))
	r.Methods(http.MethodGet).Use(first, second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middlewares\n")
	}))

	cases := []routerTest{
		{
			path:   "/globalmiddleware",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "global: before\n/globalmiddleware\nglobal: after\n",
		},
		{
			path:   "/middleware",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "global: before\nfirst: before\n/middleware\nfirst: after\nglobal: after\n",
		},
		{
			path:   "/middlewares",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "global: before\nfirst: before\nsecond: before\nthird: before\n/middlewares\nthird: after\nsecond: after\nfirst: after\nglobal: after\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name(), func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
			}

			recBody, _ := io.ReadAll(rec.Body)
			body := string(recBody)
			if body != c.body {
				t.Errorf("actual: %v expected: %v\n", body, c.body)
			}
		})
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodGet).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/bar/:id/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))
	r.Methods(http.MethodPost).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodPost).Handler(`/foo`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.Methods(http.MethodPost).Handler(`/foo/bar`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.Methods(http.MethodPost).Handler(`/foo/bar/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.Methods(http.MethodPost).Handler(`/foo/bar/:id/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.Methods(http.MethodPost).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.Methods(http.MethodPost).Handler(`/foo/:id[^\d+$]/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))
	r.Methods(http.MethodOptions).Handler(`/options`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/options")
	}))
	r.Methods(http.MethodOptions).Handler(`/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/%v", id)
	}))

	cases := []routerTest{
		{
			path:   "/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			path:   "/foo",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo",
		},
		{
			path:   "/foo/bar",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar",
		},
		{
			path:   "/foo/bar/123",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar/123",
		},
		{
			path:   "/foo/bar/123/john",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar/123/john",
		},
		{
			path:   "/foo/123",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/123",
		},
		{
			path:   "/foo/123/john",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/123/john",
		},
		{
			path:   "/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			path:   "/foo",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo",
		},
		{
			path:   "/foo/bar",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar",
		},
		{
			path:   "/foo/bar/123",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar/123",
		},
		{
			path:   "/foo/bar/123/john",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar/123/john",
		},
		{
			path:   "/foo/123",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/123",
		},
		{
			path:   "/foo/123/john",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/123/john",
		},
		{
			path:   "/options",
			method: http.MethodOptions,
			code:   http.StatusOK,
			body:   "/options",
		},
		{
			path:   "/1",
			method: http.MethodOptions,
			code:   http.StatusOK,
			body:   "/1",
		},
	}

	for _, c := range cases {
		t.Run(c.name(), func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
			}

			recBody, _ := io.ReadAll(rec.Body)
			body := string(recBody)
			if body != c.body {
				t.Errorf("actual: %v expected: %v\n", body, c.body)
			}
		})
	}
}

func TestRouterPathValue(t *testing.T) {
	r := NewRouter()

	var mwID string
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mwID = r.PathValue("id")
			next.ServeHTTP(w, r)
		})
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "id=%v name=%v getparam_id=%v", r.PathValue("id"), r.PathValue("name"), GetParam(r.Context(), "id"))
	})

	r.Methods(http.MethodGet).Handler(`/foo/:id/:name`, h)
	r.Methods(http.MethodGet).Handler(`/regex/:id[^\d+$]`, h)
	r.Methods(http.MethodGet).Handler(`/static`, h)
	r.Methods(http.MethodGet).Handler(`/dup/:id/:id`, h)
	r.Methods(http.MethodGet).Use(mw).Handler(`/middleware/:id`, h)

	cases := []struct {
		name string
		path string
		body string
		mwID string
	}{
		{
			name: "named params",
			path: "/foo/1/goblin",
			body: "id=1 name=goblin getparam_id=1",
		},
		{
			name: "regex param",
			path: "/regex/123",
			body: "id=123 name= getparam_id=123",
		},
		{
			name: "static route has no path values",
			path: "/static",
			body: "id= name= getparam_id=",
		},
		{
			name: "duplicated param name resolves to the first value like GetParam",
			path: "/dup/first/second",
			body: "id=first name= getparam_id=first",
		},
		{
			name: "middleware can read path values",
			path: "/middleware/1",
			body: "id=1 name= getparam_id=1",
			mwID: "1",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mwID = ""
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, c.path, nil))

			if rec.Code != http.StatusOK {
				t.Errorf("actual: %v expected: %v\n", rec.Code, http.StatusOK)
			}
			if rec.Body.String() != c.body {
				t.Errorf("actual: %v expected: %v\n", rec.Body.String(), c.body)
			}
			if mwID != c.mwID {
				t.Errorf("actual: %v expected: %v\n", mwID, c.mwID)
			}
		})
	}
}

func TestDefaultErrorHandler(t *testing.T) {
	r := NewRouter()
	r.Methods(http.MethodGet).Handler(`/defaulterrorhandler`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	r.Methods(http.MethodGet).Handler(`/methodnotallowed`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	cases := []routerTest{
		{
			path:   "/",
			method: http.MethodGet,
			code:   http.StatusNotFound,
		},
		{
			path:   "/methodnotallowed",
			method: http.MethodPost,
			code:   http.StatusMethodNotAllowed,
		},
	}

	for _, c := range cases {
		t.Run(c.name(), func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
			}
		})
	}
}

func TestCustomErrorHandler(t *testing.T) {
	r := NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "statusnotfound")
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "methodnotallowed")
	})
	r.Methods(http.MethodGet).Handler(`/custommethodnotfound`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	r.Methods(http.MethodGet).Handler(`/custommethodnotallowed`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	cases := []routerTest{
		{
			path:   "/",
			method: http.MethodGet,
			code:   http.StatusNotFound,
			body:   "statusnotfound",
		},
		{
			path:   "/custommethodnotallowed",
			method: http.MethodPost,
			code:   http.StatusMethodNotAllowed,
			body:   "methodnotallowed",
		},
	}

	for _, c := range cases {
		t.Run(c.name(), func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
			}

			recBody, _ := io.ReadAll(rec.Body)
			body := string(recBody)
			if body != c.body {
				t.Errorf("actual: %v expected: %v\n", body, c.body)
			}
		})
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	srv := httptest.NewServer(methodNotAllowedHandler())
	defer srv.Close()
	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("actual: %v expected: %v\n", err, nil)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("actual: %v expected: %v\n", res.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestDefaultOPTIONSHandler(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r := NewRouter()
	r.DefaultOPTIONSHandler = fn
	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("actual: %v expected: %v\n", rec.Code, http.StatusNoContent)
	}
}
