package env

import (
	"context"
)

type varsKey struct{}

// WithVars adds environment variables to a context.
func WithVars(ctx context.Context, vars Vars) context.Context {
	return context.WithValue(ctx, varsKey{}, vars)
}

// GetVars gets environment variables from a context.
func GetVars(ctx context.Context) Vars {
	if raw := ctx.Value(varsKey{}); raw != nil {
		if typed, ok := raw.(Vars); ok {
			return typed
		}
	}
	return Env()
}
