package profanity

import (
	"regexp"
)

// RegexFilter represents rules around matching (or excluding) based on
// regular expressions.
type RegexFilter struct {
	Filter `yaml:",inline"`
}

// Match returns the matching glob filter for a given value.
func (rf RegexFilter) Match(value string) (includeMatch, excludeMatch string) {
	return rf.Filter.Match(value, rf.MustMatch)
}

// Allow returns if the filters include or exclude a given filename.
func (rf RegexFilter) Allow(value string) (result bool) {
	return rf.Filter.Allow(value, rf.MustMatch)
}

// MustMatch regexp but panics
func (rf RegexFilter) MustMatch(value, expr string) bool {
	result, err := regexp.MatchString(expr, value)
	if err != nil {
		panic(err)
	}
	return result
}
