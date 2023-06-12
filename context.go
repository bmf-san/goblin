package goblin

import (
	"context"
)

// paramsKey represents the key for parameters
type paramsKey struct{}

// ParamsKey is the request context key under which URL params are stored.
var ParamsKey = paramsKey{}

// GetParam gets parameters from request.
func GetParam(ctx context.Context, name string) string {
	params, ok := ctx.Value(ParamsKey).(Params)
	if !ok {
		return ""
	}

	for i := 0; i < len(params); i++ {
		if params[i].key == name {
			return params[i].value
		}
	}

	return ""
}
