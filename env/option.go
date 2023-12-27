/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package env

// Option is a mutator for the options set.
type Option func(Vars)

// OptEnviron parses the output of `os.Environ()`
func OptEnviron(environ ...string) Option {
	return func(vars Vars) {
		var key, value string
		for _, envVar := range environ {
			key, value = Split(envVar)
			if key != "" {
				vars[key] = value
			}
		}
	}
}

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
