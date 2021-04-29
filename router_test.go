package goblin

import (
	"context"
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

func first(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "first: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "first: after\n")
	})
}

func second(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "second: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "second: after\n")
	})
}

func third(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "third: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "third: after\n")
	})
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
	r.OPTION(`/option`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/option")
	}))
	r.OPTION(`:id`).Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			path:   "/option",
			method: http.MethodOptions,
			code:   http.StatusOK,
			body:   "/option",
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

		if c.code != rec.Code {
			t.Errorf("path:%v code:%v body:%v rec.Code:%v", c.path, c.code, c.body, rec.Code)
		}

		recBody, _ := ioutil.ReadAll(rec.Body)
		body := string(recBody)
		if c.body != body {
			t.Errorf("path:%v code:%v body:%v rec.Body:%v", c.path, c.code, c.body, body)
		}
	}
}

func TestGetParam(t *testing.T) {
	params := &Params{
		&Param{
			key:   "id",
			value: "123",
		},
		&Param{
			key:   "name",
			value: "john",
		},
	}

	ctx := context.WithValue(context.Background(), ParamsKey, *params)

	cases := []struct {
		actual   string
		expected string
	}{
		{
			actual:   GetParam(ctx, "id"),
			expected: "123",
		},
		{
			actual:   GetParam(ctx, "name"),
			expected: "john",
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
