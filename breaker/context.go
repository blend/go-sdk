package breaker

import (
	"context"
	"time"
)

type contextTimeKey struct{}

// WithContextTime adds a relevant timestamp to a context.
func WithContextTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, contextTimeKey{}, t)
}

// GetContextTime gets the context timestamp off a context.
func GetContextTime(ctx context.Context) (t time.Time) {
	if val := ctx.Value(contextTimeKey{}); val != nil {
		if typed, ok := val.(time.Time); ok {
			t = typed
			return
		}
	}
	return
}
