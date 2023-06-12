package goblin

import (
	"context"
	"testing"
)

func TestGetParam(t *testing.T) {
	params := Params{
		{
			key:   "id",
			value: "123",
		},
		{
			key:   "name",
			value: "john",
		},
	}

	ctx := context.WithValue(context.Background(), ParamsKey, params)
	ngCtx := context.WithValue(context.Background(), ParamsKey, "not a []Param")

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
		{
			actual:   GetParam(ctx, "not-exist-key"),
			expected: "",
		},
		{
			actual:   GetParam(ngCtx, "ng ctx"),
			expected: "",
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
