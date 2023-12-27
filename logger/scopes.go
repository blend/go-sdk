/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/stringutil"
)

// NewScopes yields a scope set from a list of scopes.
//
// Each scope should be path formatted, e.g. `foo/bar/baz` and can include
// wildcard segments for glob matching, e.g. `foo/*/baz`.
//
// A single `*` is interpretted as `All` and will match any scope path.
func NewScopes(scopes ...string) *Scopes {
	scopeSet := &Scopes{
		scopes: make(map[string]bool),
	}
	for _, rawScope := range scopes {
		parsedScope := strings.ToLower(strings.TrimSpace(rawScope))
		if parsedScope == ScopeAll {
			scopeSet.all = true
			continue
		}
		if strings.HasPrefix(parsedScope, "-") {
			scopeSet.scopes[strings.TrimPrefix(parsedScope, "-")] = false
		} else {
			scopeSet.scopes[parsedScope] = true
		}
	}
	return scopeSet
}

// ScopesAll returns a preset scopes with the all flag flipped.
func ScopesAll() *Scopes {
	return &Scopes{scopes: make(map[string]bool), all: true}
}

// ScopesNone returns a preset empty scopes.
func ScopesNone() *Scopes {
	return &Scopes{scopes: make(map[string]bool)}
}

// Scopes is a set of scopes.
type Scopes struct {
	all	bool
	scopes	map[string]bool
}

// Enable enables a set of scopes.
//
// The scopes should be given in filepath form, e.g. `foo/bar/*`.
func (s *Scopes) Enable(scopes ...string) {
	for _, scope := range scopes {
		s.scopes[strings.ToLower(strings.TrimSpace(scope))] = true
	}
}

// Disable disables a set of scopes.
//
// The scopes should be given in filepath form, e.g. `foo/bar/*`.
func (s *Scopes) Disable(scopes ...string) {
	for _, scope := range scopes {
		s.scopes[strings.ToLower(strings.TrimSpace(scope))] = false
	}
}

// SetAll flips the `all` bit on the flag set to true.
//
// Note: flags that are explicitly disabled will remain disabled.
func (s *Scopes) SetAll() {
	s.all = true
}

// All returns if the all bit is flipped to true.
func (s *Scopes) All() bool {
	return s.all
}

// SetNone disables the `all` bit, and empties the scopes
// set, resulting in calls to `None()` to return true.
//
// You should view this method as a way to reset or zero a scopes set.
func (s *Scopes) SetNone() {
	s.all = false
	s.scopes = make(map[string]bool)
}

// None returns if the all bit is set to false, and
// there are no scope specific overrides.
//
// It is functionally equivalent to an `IsZero()` method.
func (s *Scopes) None() bool {
	return !s.all && len(s.scopes) == 0
}

// IsEnabled returns if a given logger scope is enabled.
func (s Scopes) IsEnabled(scopePath ...string) bool {
	scopeJoined := filepath.Join(scopePath...)
	if s.all {
		// check if we explicitly disabled the scope
		return !s.isScopeExplicitlyDisabled(scopeJoined)
	}
	// we treat empty scopes as `none`
	if len(s.scopes) == 0 {
		return false
	}
	// if there is no scope, assumed to pass
	if len(scopeJoined) == 0 {
		return true
	}
	// return if we either explicitly enabled or disabled the scope
	// with path globs in the scopes map.
	return s.isScopeEnabled(scopeJoined)
}

// String returns a string representation of the scopes.
func (s Scopes) String() string {
	return strings.Join(s.Scopes(), ", ")
}

// Scopes returns an array of scopes.
func (s Scopes) Scopes() []string {
	var scopes []string
	if s.all {
		scopes = []string{ScopeAll}
	}
	for key, enabled := range s.scopes {
		if key != FlagAll {
			if enabled {
				if !s.all {
					scopes = append(scopes, string(key))
				}
			} else {
				scopes = append(scopes, "-"+string(key))
			}
		}
	}
	return scopes
}

//
// internal helpers
//

// isScopeEnabled returns if a scopePath is enabled strictly by
// a lookup to the underlying scopes map.
func (s Scopes) isScopeEnabled(scopePath string) bool {
	for pattern, enabled := range s.scopes {
		if s.matches(scopePath, pattern) {
			return enabled
		}
	}
	// no matching entry is a failure.
	return false
}

// isScopeDisabled returns if a scopePath is explicitly disabled
// that is, has a matching glob in the scopes map that is set to false.
//
// it is differentiated from `isScopeEnabled`
func (s Scopes) isScopeExplicitlyDisabled(subj string) bool {
	for pattern, enabled := range s.scopes {
		if !enabled && s.matches(subj, pattern) {
			return true
		}
	}
	// if we didn't find a matching scope, assume it's not disabled
	return false
}

func (s Scopes) matches(subj, pattern string) (output bool) {
	output = stringutil.Glob(subj, pattern)
	return
}
