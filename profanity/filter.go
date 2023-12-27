/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

import (
	"fmt"
	"strings"
)

// Filter is the base rule helper.
type Filter struct {
	// Include sets a glob filter for file inclusion by name.
	Include	[]string	`yaml:"include,omitempty"`
	// ExcludeGlob sets a glob filter for file exclusion by name.
	Exclude	[]string	`yaml:"exclude,omitempty"`
}

// IsZero returns if the filter is set or not.
func (f Filter) IsZero() bool {
	return len(f.Include) == 0 && len(f.Exclude) == 0
}

// Match returns the matching glob filter for a given value.
func (f Filter) Match(value string, filter func(string, string) bool) (includeMatch, excludeMatch string) {
	if len(f.Include) > 0 {
		for _, include := range f.Include {
			if filter(value, include) {
				includeMatch = include
				break
			}
		}
	}
	if len(f.Include) == 0 || includeMatch != "" {
		if len(f.Exclude) > 0 {
			for _, exclude := range f.Exclude {
				if filter(value, exclude) {
					excludeMatch = exclude
					break
				}
			}
		}
	}
	return
}

// AllowMatch returns if the filter should allow a given match set.
func (f Filter) AllowMatch(includeMatch, excludeMatch string) bool {
	if len(f.Include) > 0 && len(f.Exclude) > 0 {
		return includeMatch != "" && excludeMatch == ""
	}
	if len(f.Include) > 0 {
		return includeMatch != ""
	}
	if len(f.Exclude) > 0 {
		return excludeMatch == ""
	}
	return true
}

// Allow returns if the filters include or exclude a given value.
func (f Filter) Allow(value string, filter func(string, string) bool) bool {
	return f.AllowMatch(f.Match(value, filter))
}

// Validate doesn't do anything right now for Filter.
func (f Filter) Validate() error {
	return nil
}

// String implements fmt.Stringer.
func (f Filter) String() string {
	if len(f.Include) > 0 && len(f.Exclude) > 0 {
		return fmt.Sprintf("[include: %s, exclude: %s]",
			strings.Join(f.Include, ", "),
			strings.Join(f.Exclude, ", "),
		)
	} else if len(f.Include) > 0 {
		return fmt.Sprintf("[include: %s]", strings.Join(f.Include, ", "))
	} else if len(f.Exclude) > 0 {
		return fmt.Sprintf("[exclude: %s]", strings.Join(f.Exclude, ", "))
	}
	return ""
}
