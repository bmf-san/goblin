// Borrowed code from go-router-benchmark
// See: https://github.com/bmf-san/go-router-benchmark
package goblin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// routeSet is a struct for routeSet.
type routeSet struct {
	path    string
	reqPath string
}

var (
	staticRoutesRoot       = routeSet{"/", "/"}
	staticRoutes1          = routeSet{"/foo", "/foo"}
	staticRoutes5          = routeSet{"/foo/bar/baz/qux/quux", "/foo/bar/baz/qux/quux"}
	staticRoutes10         = routeSet{"/foo/bar/baz/qux/quux/corge/grault/garply/waldo/fred", "/foo/bar/baz/qux/quux/corge/grault/garply/waldo/fred"}
	pathParamRoutes1Colon  = routeSet{"/foo/:bar", "/foo/bar"}
	pathParamRoutes5Colon  = routeSet{"/foo/:bar/:baz/:qux/:quux/:corge", "/foo/bar/baz/qux/quux/corge"}
	pathParamRoutes10Colon = routeSet{"/foo/:bar/:baz/:qux/:quux/:corge/:grault/:garply/:waldo/:fred/:plugh", "/foo/bar/baz/qux/quux/corge/grault/garply/waldo/fred/plugh"}
)

func loadGoblin(r routeSet) http.Handler {
	router := NewRouter()
	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	router.Methods(http.MethodGet).Handler(r.path, handler)
	return router
}

func testServeHTTP(b *testing.B, r routeSet, router http.Handler) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, r.reqPath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rec, req)
		if rec.Code != 200 {
			panic(fmt.Sprintf("Request failed. path: %v request path:%v", r.path, r.reqPath))
		}
	}
}

func benchmark(b *testing.B, r routeSet, router http.Handler) {
	testServeHTTP(b, r, router)
}

// Benchmark tests
func BenchmarkStaticRoutesRootGoblin(b *testing.B) {
	router := loadGoblin(staticRoutesRoot)
	benchmark(b, staticRoutesRoot, router)
}

func BenchmarkStaticRoutes1Goblin(b *testing.B) {
	router := loadGoblin(staticRoutes1)
	benchmark(b, staticRoutes1, router)
}

func BenchmarkStaticRoutes5Goblin(b *testing.B) {
	router := loadGoblin(staticRoutes5)
	benchmark(b, staticRoutes5, router)
}

func BenchmarkStaticRoutes10Goblin(b *testing.B) {
	router := loadGoblin(staticRoutes10)
	benchmark(b, staticRoutes10, router)
}

func BenchmarkPathParamRoutes1ColonGoblin(b *testing.B) {
	router := loadGoblin(pathParamRoutes1Colon)
	benchmark(b, pathParamRoutes1Colon, router)
}

func BenchmarkPathParamRoutes5ColonGoblin(b *testing.B) {
	router := loadGoblin(pathParamRoutes5Colon)
	benchmark(b, pathParamRoutes5Colon, router)
}

func BenchmarkPathParamRoutes10ColonGoblin(b *testing.B) {
	router := loadGoblin(pathParamRoutes10Colon)
	benchmark(b, pathParamRoutes10Colon, router)
}
