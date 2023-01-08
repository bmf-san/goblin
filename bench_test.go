package goblin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func handler() http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
}

func benchmarkSetRoutes(n int, b *testing.B) {
	r := NewRouter()
	path := "/path"
	pathN := "/pathN"

	for i := 0; i < n; i++ {
		path += pathN
	}

	b.ResetTimer()
	r.Methods(http.MethodGet).Handler(path, handler())
}

func benchmarkStatic(n int, b *testing.B) {
	r := NewRouter()
	path := "/static"
	pathN := "/pathN"

	for i := 0; i < n; i++ {
		path += pathN
	}

	r.Methods(http.MethodGet).Handler(path, handler())
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ServeHTTP(rec, req)
		}
	})
}

func benchmarkWildcard(n int, b *testing.B) {
	r := NewRouter()
	var path, rpath string
	wildcard := "/:wildcard"
	pathN := "/pathN"

	for i := 0; i < n; i++ {
		path += wildcard
		rpath += pathN
	}

	r.Methods(http.MethodGet).Handler(path, handler())
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ServeHTTP(rec, req)
		}
	})
}

func benchmarkRegexp(n int, b *testing.B) {
	r := NewRouter()
	var path, rpath string
	wildcard := "/:*[(.+)]"
	pathN := "/pathN"

	for i := 0; i < n; i++ {
		path += wildcard
		rpath += pathN
	}

	r.Methods(http.MethodGet).Handler(path, handler())
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ServeHTTP(rec, req)
		}
	})
}

func BenchmarkSetRoutes1(b *testing.B) {
	benchmarkSetRoutes(1, b)
}

func BenchmarkSetRoutes5(b *testing.B) {
	benchmarkSetRoutes(5, b)
}

func BenchmarkSetRoutes10(b *testing.B) {
	benchmarkSetRoutes(10, b)
}

func BenchmarkStatic1(b *testing.B) {
	benchmarkStatic(1, b)
}

func BenchmarkStatic5(b *testing.B) {
	benchmarkStatic(5, b)
}

func BenchmarkStatic10(b *testing.B) {
	benchmarkStatic(10, b)
}

func BenchmarkWildCard1(b *testing.B) {
	benchmarkWildcard(1, b)
}

func BenchmarkWildCard5(b *testing.B) {
	benchmarkWildcard(5, b)
}

func BenchmarkWildCard10(b *testing.B) {
	benchmarkWildcard(10, b)
}

func BenchmarkRegexp1(b *testing.B) {
	benchmarkRegexp(1, b)
}

func BenchmarkRegexp5(b *testing.B) {
	benchmarkRegexp(5, b)
}

func BenchmarkRegexp10(b *testing.B) {
	benchmarkRegexp(10, b)
}
