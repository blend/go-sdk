package env

import (
	"os"
	"strings"
)

// Option is a mutator for the options set.
type Option func(Vars)

// OptSet overrides values in the set with a specific set of values.
func OptSet(overides Vars) Option {
	return func(vars Vars) {
		for key, value := range overides {
			vars[key] = value
		}
	}
}

// OptRemove removes keys from a set.
func OptRemove(keys ...string) Option {
	return func(vars Vars) {
		for _, key := range keys {
			delete(vars, key)
		}
	}
}

// OptFromEnv sets the vars from the current os environment.
func OptFromEnv() Option {
	return func(v Vars) {
		envVars := os.Environ()
		for _, ev := range envVars {
			parts := strings.SplitN(ev, "=", 2)
			if len(parts) > 1 {
				v[parts[0]] = parts[1]
			}
		}
	}
}
