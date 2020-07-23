package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

type tracerKey struct{}

// WithTracer adds a tracer to a context.
func WithTracer(ctx context.Context, tracer opentracing.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

// GetTracer gets a tracer from a context.
func GetTracer(ctx context.Context) opentracing.Tracer {
	if value := ctx.Value(tracerKey{}); value != nil {
		if typed, ok := value.(opentracing.Tracer); ok {
			return typed
		}
	}
	return nil
}
