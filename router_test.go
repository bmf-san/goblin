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

func TestRouter(t *testing.T) {
	r := NewRouter()

	r.GET(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.GET(`/foo/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/")
	}))
	r.GET(`/foo/bar/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar/")
	}))
	r.GET(`/foo/bar/:id/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v/", id)
	}))
	r.GET(`/foo/bar/:id/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v/", id, name)
	}))
	r.GET(`/foo/:id[^\d+$]/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v/", id)
	}))
	r.GET(`/foo/:id[^\d+$]/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v/", id, name)
	}))

	r.POST(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/")
	}))
	r.POST(`/foo/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/")
	}))
	r.POST(`/foo/bar/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/foo/bar/")
	}))
	r.POST(`/foo/bar/:id/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/bar/%v/", id)
	}))
	r.POST(`/foo/bar/:id/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/bar/%v/%v/", id, name)
	}))
	r.POST(`/foo/:id[^\d+$]/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		fmt.Fprintf(w, "/foo/%v/", id)
	}))
	r.POST(`/foo/:id[^\d+$]/:name/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r.Context(), "id")
		name := GetParam(r.Context(), "name")
		fmt.Fprintf(w, "/foo/%v/%v/", id, name)
	}))

	cases := []struct {
		url    string
		method string
		code   int
		body   string
	}{
		{
			url:    "/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			url:    "/foo/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/",
		},
		{
			url:    "/foo/bar/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar/",
		},
		{
			url:    "/foo/bar/123/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar/123/",
		},
		{
			url:    "/foo/bar/123/john/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/bar/123/john/",
		},
		{
			url:    "/foo/123/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/123/",
		},
		{
			url:    "/foo/123/john/",
			method: http.MethodGet,
			code:   http.StatusOK,
			body:   "/foo/123/john/",
		},
		{
			url:    "/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/",
		},
		{
			url:    "/foo/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/",
		},
		{
			url:    "/foo/bar/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar/",
		},
		{
			url:    "/foo/bar/123/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar/123/",
		},
		{
			url:    "/foo/bar/123/john/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/bar/123/john/",
		},
		{
			url:    "/foo/123/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/123/",
		},
		{
			url:    "/foo/123/john/",
			method: http.MethodPost,
			code:   http.StatusOK,
			body:   "/foo/123/john/",
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest(c.method, c.url, nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if c.code != rec.Code {
			t.Errorf("url:%v code:%v body:%v rec.Code:%v", c.url, c.code, c.body, rec.Code)
		}

		recBody, _ := ioutil.ReadAll(rec.Body)
		body := string(recBody)
		if c.body != body {
			t.Errorf("url:%v code:%v body:%v rec.Body:%v", c.url, c.code, c.body, body)
		}
	}
}

func TestGetParam(t *testing.T) {
	params := &Params{
		Param{
			key:   "id",
			value: "123",
		},
		Param{
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
