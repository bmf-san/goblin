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
		name     string
		actual   string
		expected string
	}{
		{
			name:     "id_param",
			actual:   GetParam(ctx, "id"),
			expected: "123",
		},
		{
			name:     "name_param",
			actual:   GetParam(ctx, "name"),
			expected: "john",
		},
		{
			name:     "not_exist_param",
			actual:   GetParam(ctx, "not-exist-key"),
			expected: "",
		},
		{
			name:     "param_value_wrong_type",
			actual:   GetParam(ngCtx, "ng ctx"),
			expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.actual != c.expected {
				t.Errorf("actual:%v expected:%v", c.actual, c.expected)
			}
		})
	}
}
