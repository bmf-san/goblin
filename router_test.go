package goblin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	actual := NewRouter()
	expected := &Router{
		tree: NewTree(),
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v\n", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.GET(`/`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.GET(`/middleware`).Use(first).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middleware\n")
	}))
	r.GET(`/middlewares`).Use(first, second, third).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/middlewares\n")
	}))
	r.GET(`/foo`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.GET(`/foo/bar`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.GET(`/foo/bar/:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.GET(`/foo/bar/:id/:name`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.GET(`/foo/:id[^\d+$]`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.GET(`/foo/:id[^\d+$]/:name`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))
	r.POST(`/`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.POST(`/foo`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo")
	}))
	r.POST(`/foo/bar`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar")
	}))
	r.POST(`/foo/bar/:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v", id)
	}))
	r.POST(`/foo/bar/:id/:name`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v", id, name)
	}))
	r.POST(`/foo/:id[^\d+$]`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v", id)
	}))
	r.POST(`/foo/:id[^\d+$]/:name`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v", id, name)
	}))
	r.OPTIONS(`/options`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/options")
	}))
	r.OPTIONS(`:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		recBody, _ := ioutil.ReadAll(rec.Body)
		body := string(recBody)
		if body != c.body {
			t.Errorf("actual: %v expected: %v\n", body, c.body)
		}
	}
}
