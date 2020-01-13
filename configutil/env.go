package configutil

import (
	"time"

	"github.com/blend/go-sdk/env"
)

var (
	_ StringSource   = (*Env)(nil)
	_ BoolSource     = (*Env)(nil)
	_ IntSource      = (*Env)(nil)
	_ Float64Source  = (*Env)(nil)
	_ DurationSource = (*Env)(nil)
)

// Env is a value provider where the string represents the environment variable name.
// It can be used with *any* config.Set___ type.
type Env string

// String returns a given environment variable as a string.
func (e Env) String() (*string, error) {
	key := string(e)
	if env.Env().Has(key) {
		value := env.Env().String(key)
		return &value, nil
	}
	return nil, nil
}

// Strings returns a given environment variable as strings.
func (e Env) Strings() ([]string, error) {
	key := string(e)
	if env.Env().Has(key) {
		return env.Env().CSV(key), nil
	}
	return nil, nil
}

// Bool returns a given environment variable as a bool.
func (e Env) Bool() (*bool, error) {
	key := string(e)
	if env.Env().Has(key) {
		value := env.Env().Bool(key)
		return &value, nil
	}
	return nil, nil
}

// Int returns a given environment variable as an int.
func (e Env) Int() (*int, error) {
	key := string(e)
	if env.Env().Has(key) {
		value, err := env.Env().Int(key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// Float64 returns a given environment variable as a float64.
func (e Env) Float64() (*float64, error) {
	key := string(e)
	if env.Env().Has(key) {
		value, err := env.Env().Float64(key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}

// Duration returns a given environment variable as a time.Duration.
func (e Env) Duration() (*time.Duration, error) {
	key := string(e)
	if env.Env().Has(key) {
		value, err := env.Env().Duration(key)
		if err != nil {
			return nil, err
		}
		return &value, nil
	}
	return nil, nil
}
