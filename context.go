package goblin

import "context"

type key int

const (
	// ParamsKey is the key in a request context.
	ParamsKey key = iota
)

// GetParam gets parameters from request.
func GetParam(ctx context.Context, name string) string {
	params, _ := ctx.Value(ParamsKey).(Params)

	for i := range params {
		if params[i].key == name {
			return params[i].value
		}
	}

	return ""
}
