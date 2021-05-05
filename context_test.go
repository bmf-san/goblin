package goblin

import "testing"

func TestGetParameter(t *testing.T) {
	cases := []struct {
		actual   string
		expected string
	}{
		{
			actual:   getParameter(`:id[^\d+$]`),
			expected: "id",
		},
		{
			actual:   getParameter(`:id[`),
			expected: "id",
		},
		{
			actual:   getParameter(`:id]`),
			expected: "id]",
		},
		{
			actual:   getParameter(`:id`),
			expected: "id",
		},
	}

	for _, c := range cases {
		if c.actual != c.expected {
			t.Errorf("actual:%v expected:%v", c.actual, c.expected)
		}
	}
}
