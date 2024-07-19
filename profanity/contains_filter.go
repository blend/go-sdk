/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import "strings"

// ContainsFilter is the contains filter.
type ContainsFilter struct {
	Filter `yaml:",inline"`
}

// Match applies the filter.
func (c ContainsFilter) Match(value string) (includeMatch, excludeMatch string) {
	return c.Filter.Match(value, strings.Contains)
}

// Allow returns if apply returns a result.
func (c ContainsFilter) Allow(value string) bool {
	return c.Filter.Allow(value, strings.Contains)
}
