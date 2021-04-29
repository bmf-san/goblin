package goblin

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	actual := NewMiddlewares(nil)
	var expected middlewares

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v\n", actual, expected)
	}
}

func handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestThen(t *testing.T) {
	t.Skip("This method covers the tests in an associative manner with router tests, so skip it.")
}
