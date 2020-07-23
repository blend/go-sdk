package configutil

import (
	"context"
	"time"

	"github.com/blend/go-sdk/env"
)

var (
	_ StringSource   = (*EnvVars)(nil)
	_ BoolSource     = (*EnvVars)(nil)
	_ IntSource      = (*EnvVars)(nil)
	_ Float64Source  = (*EnvVars)(nil)
	_ DurationSource = (*EnvVars)(nil)
)

// Env returns a new environment value provider.
func Env(key string) EnvVars {
	return EnvVars{
		Key: key,
	}
}

// EnvVars is a value provider where the string represents the environment variable name.
// It can be used with *any* config.Set___ type.
type EnvVars struct {
	Key string
}

// String returns a given environment variable as a string.
func (e EnvVars) String(ctx context.Context) (*string, error) {
	vars := e.vars(ctx)
	if vars.Has(e.Key) {
		value := vars.String(e.Key)
		return &value, nil
	}
	return nil, nil
}

// Strings returns a given environment variable as strings.
func (e EnvVars) Strings(ctx context.Context) ([]string, error) {
	vars := e.vars(ctx)

	if vars.Has(e.Key) {
		return vars.CSV(e.Key), nil
	}
	return nil, nil
}

// Bool returns a given environment variable as a bool.
func (e EnvVars) Bool(ctx context.Context) (*bool, error) {
	vars := e.vars(ctx)

	if vars.Has(e.Key) {
		value := vars.Bool(e.Key)
		return &value, nil
	}
	return nil, nil
}

// Int returns a given environment variable as an int.
func (e EnvVars) Int(ctx context.Context) (*int, error) {
	vars := e.vars(ctx)

	if vars.Has(e.Key) {
		value, err := vars.Int(e.Key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// Float64 returns a given environment variable as a float64.
func (e EnvVars) Float64(ctx context.Context) (*float64, error) {
	vars := e.vars(ctx)

	if vars.Has(e.Key) {
		value, err := vars.Float64(e.Key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// Duration returns a given environment variable as a time.Duration.
func (e EnvVars) Duration(ctx context.Context) (*time.Duration, error) {
	vars := e.vars(ctx)

	if vars.Has(e.Key) {
		value, err := vars.Duration(e.Key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// vars resolves the env vars
func (e EnvVars) vars(ctx context.Context) env.Vars {
	if vars := env.GetVars(ctx); vars != nil {
		return vars
	}
	return env.Env()
}
