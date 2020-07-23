package web

import (
	"context"
	"time"
)

type appKey struct{}

// WithApp adds an app to a context.
func WithApp(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appKey{}, app)
}

// GetApp gets an app off a context.
func GetApp(ctx context.Context) *App {
	if value := ctx.Value(appKey{}); value != nil {
		if typed, ok := value.(*App); ok {
			return typed
		}
	}
	return nil
}

type requestStartedKey struct{}

// WithRequestStarted sets the request started time on a context.
func WithRequestStarted(ctx context.Context, requestStarted time.Time) context.Context {
	return context.WithValue(ctx, requestStartedKey{}, requestStarted)
}

// GetRequestStarted gets the request started time from a context.
func GetRequestStarted(ctx context.Context) time.Time {
	if value := ctx.Value(requestStartedKey{}); value != nil {
		if typed, ok := value.(time.Time); ok {
			return typed
		}
	}
	return time.Now().UTC()
}
