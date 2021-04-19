package goblin

import (
	"reflect"
	"strings"
	"testing"
)

func TestDeleteEmpty(t *testing.T) {
	cases := []struct {
		item     string
		expected []string
	}{
		{
			item:     "/foo/bar/baz",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			item:     "/foo/baz",
			expected: []string{"foo", "baz"},
		},
		{
			item:     "/foo/",
			expected: []string{"foo"},
		},
		{
			item:     "/",
			expected: []string(nil),
		},
	}

	for _, c := range cases {
		actual := deleteEmpty(strings.Split(c.item, "/"))
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("actual:%v expected:%v", actual, c.expected)
		}
	}
}
