/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import (
	"regexp"
)

// RegexFilter represents rules around matching (or excluding) based on
// regular expressions.
type RegexFilter struct {
	Filter `yaml:",inline"`

	compiledExpressions map[string]*regexp.Regexp
}

// Match returns the matching glob filter for a given value.
func (rf *RegexFilter) Match(value string) (includeMatch, excludeMatch string) {
	return rf.Filter.Match(value, rf.MustMatch)
}

// Allow returns if the filters include or exclude a given filename.
func (rf *RegexFilter) Allow(value string) (result bool) {
	return rf.Filter.Allow(value, rf.MustMatch)
}

// MustMatch regexp but panics
func (rf *RegexFilter) MustMatch(value, expr string) bool {
	if rf.compiledExpressions == nil {
		rf.compiledExpressions = make(map[string]*regexp.Regexp)
	}
	if expr, ok := rf.compiledExpressions[expr]; ok {
		return expr.MatchString(value)
	}
	compiled, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	rf.compiledExpressions[expr] = compiled
	return compiled.MatchString(value)
}
