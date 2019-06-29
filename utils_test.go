package goblin

import (
	"reflect"
	"strings"
	"testing"
)

func TestDeleteEmpty(t *testing.T) {
	s := strings.Split("/foo/bar/baz", "/")

	actual := deleteEmpty(s)
	expected := []string{"foo", "bar", "baz"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v", actual, expected)
	}
}
