package goblin

import (
	"context"
	"testing"
)

func TestGetParam(t *testing.T) {
	params := &params{
		param{
			key:   "id",
			value: "123",
		},
		param{
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
