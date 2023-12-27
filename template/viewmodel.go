/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package template

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/env"
)

// Viewmodel is the template viewmodel.
// It surfaces a subset of the template api.
// It is set / accessed by the outer template.
type Viewmodel struct {
	vars	Vars
	env	env.Vars
}

// Vars returns the vars collection.
func (vm Viewmodel) Vars() Vars {
	return vm.vars
}

// Var returns the value of a variable, or panics if the variable is not set.
func (vm Viewmodel) Var(key string, defaults ...interface{}) (interface{}, error) {
	if value, hasVar := vm.vars[key]; hasVar {
		return value, nil
	}

	if len(defaults) > 0 {
		return defaults[0], nil
	}

	return nil, fmt.Errorf("template variable `%s` is unset and no default is provided", key)
}

// Has is an alias to `HasVar`.
func (vm Viewmodel) Has(key string) bool {
	return vm.HasVar(key)
}

// HasVar returns if a variable is set.
func (vm Viewmodel) HasVar(key string) bool {
	_, hasKey := vm.vars[key]
	return hasKey
}

// Env returns an environment variable.
func (vm Viewmodel) Env(key string, defaults ...string) (string, error) {
	if value, hasVar := vm.env[key]; hasVar {
		return value, nil
	}

	if len(defaults) > 0 {
		return defaults[0], nil
	}
	return "", fmt.Errorf("template environment variable `%s` is unset and no default is provided", key)
}

// HasEnv returns if an env var is set.
func (vm Viewmodel) HasEnv(key string) bool {
	_, hasKey := vm.env[key]
	return hasKey
}

// ExpandEnv replaces $var or ${var} based on the configured environment variables.
func (vm Viewmodel) ExpandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		if value, ok := vm.env[key]; ok {
			return value
		}
		return ""
	})
}
