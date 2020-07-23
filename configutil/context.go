package configutil

import (
	"context"
)

type configFilePathsKey struct{}

// WithConfigFilePaths adds config file paths to the context.
func WithConfigFilePaths(ctx context.Context, paths []string) context.Context {
	return context.WithValue(ctx, configFilePathsKey{}, paths)
}

// GetConfigFilePaths gets the config file paths from a context..
func GetConfigFilePaths(ctx context.Context) []string {
	if raw := ctx.Value(configFilePathsKey{}); raw != nil {
		if typed, ok := raw.([]string); ok {
			return typed
		}
	}
	return nil
}
