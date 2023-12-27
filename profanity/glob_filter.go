/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

// GlobFilter rules for if we should include or exclude file or directory by name.
type GlobFilter struct {
	Filter `yaml:",inline"`
}

// Match returns the matching glob filter for a given value.
func (gf GlobFilter) Match(value string) (includeMatch, excludeMatch string) {
	return gf.Filter.Match(value, Glob)
}

// Allow returns if the filters include or exclude a given value.
func (gf GlobFilter) Allow(value string) bool {
	return gf.Filter.Allow(value, Glob)
}
