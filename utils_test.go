package goblin

import (
	"reflect"
	"strings"
	"testing"
)

func TestDeleteEmpty(t *testing.T) {
	actual := deleteEmpty(strings.Split("/foo/bar/baz", "/"))
	expected := []string{"foo", "bar", "baz"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual:%v expected:%v", actual, expected)
	}
}
