package logger

import "context"

type subContextMetaKey struct{}

type subContextMeta struct {
	path   []string
	fields Fields
}

// WithSubContextMeta adds a sub context path to a context.
func WithSubContextMeta(ctx context.Context, path []string, fields Fields) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, subContextMetaKey{}, subContextMeta{path, fields})
	}
	return context.WithValue(context.Background(), subContextMetaKey{}, subContextMeta{path, fields})
}

// GetSubContextMeta adds a sub context path to a context.
func GetSubContextMeta(ctx context.Context) (path []string, fields Fields) {
	if rawValue := ctx.Value(subContextMetaKey{}); rawValue != nil {
		if typed, ok := rawValue.(subContextMeta); ok {
			path = typed.path
			fields = typed.fields
			return
		}
	}
	return
}
