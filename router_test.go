package goblin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	actual := NewRouter()
	expected := Router{
		tree: newTree(),
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v\n", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middleware\n")
	}))
	r.Methods(http.MethodGet).Use(first, second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middlewares\n")
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

	cases := []struct {
		path   string
		method string
		code   int
		body   string
	}{
		{
			path:   "/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			path:   "/middleware",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "first: before\n/middleware\nfirst: after\n",
		},
		{
			path:   "/middlewares",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "first: before\nsecond: before\nthird: before\n/middlewares\nthird: after\nsecond: after\nfirst: after\n",
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
	}
}

func TestDefaultErrorHandler(t *testing.T) {
	r := NewRouter()
	r.Methods(http.MethodGet).Handler(`/defaulterrorhandler`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	r.Methods(http.MethodGet).Handler(`/methodnotallowed`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	cases := []struct {
		path   string
		method string
		code   int
	}{
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
		req := httptest.NewRequest(c.method, c.path, nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != c.code {
			t.Errorf("actual: %v expected: %v\n", rec.Code, c.code)
		}
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

	cases := []struct {
		path   string
		method string
		code   int
		body   string
	}{
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
