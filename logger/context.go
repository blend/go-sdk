package logger

import "context"

type subScopeMetaKey struct{}

type subScopeMeta struct {
	path        []string
	labels      Labels
	annotations Annotations
}

// WithSubScopeMeta adds a sub context path to a context.
func WithSubScopeMeta(ctx context.Context, scope Scope) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, subScopeMetaKey{}, subScopeMeta{scope.Path, scope.Labels, scope.Annotations})
	}
	return context.WithValue(context.Background(), subScopeMetaKey{}, subScopeMeta{scope.Path, scope.Labels, scope.Annotations})
}

// GetSubScopeMeta adds a subscope meta to a context.
func GetSubScopeMeta(ctx context.Context) (path []string, labels Labels, annotations Annotations) {
	if rawValue := ctx.Value(subScopeMetaKey{}); rawValue != nil {
		if typed, ok := rawValue.(subScopeMeta); ok {
			path = typed.path
			labels = typed.labels
			annotations = typed.annotations
			return
		}
	}
	return
}
