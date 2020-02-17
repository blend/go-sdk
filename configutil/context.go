package configutil

import (
	"context"

	"github.com/blend/go-sdk/env"
)

type envVarsKey struct{}

// WithEnvVars adds environment variables to a context.
func WithEnvVars(ctx context.Context, vars env.Vars) context.Context {
	return context.WithValue(ctx, envVarsKey{}, vars)
}

// GetEnvVars gets environment variables from a context.
func GetEnvVars(ctx context.Context) env.Vars {
	if raw := ctx.Value(envVarsKey{}); raw != nil {
		if typed, ok := raw.(env.Vars); ok {
			return typed
		}
	}
	return nil
}

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
