package selector

import (
	"fmt"
	"strings"
)

// NotIn returns if a key does not match a set of values.
type NotIn struct {
	Key             string
	Values          []string
	PermittedValues []map[rune]bool
}

// Matches returns the selector result.
func (ni NotIn) Matches(labels Labels) bool {
	if value, hasValue := labels[ni.Key]; hasValue {
		for _, iv := range ni.Values {
			if iv == value {
				// the key does not equal any of the values
				return false
			}
		}
	}
	// the value doesn't exist.
	return true
}

// Validate validates the selector.
func (ni *NotIn) Validate(options ...SelectorOption) (err error) {
	var selector Selector = ni
	for _, option := range options {
		option(selector)
	}

	err = CheckKey(ni.Key)
	if err != nil {
		return
	}
	for _, v := range ni.Values {
		err = CheckValue(v, ni.PermittedValues...)
		if err != nil {
			return
		}
	}
	return
}

// AddPermittedValues adds runes to be accepted in values
func (ni *NotIn) AddPermittedValues(permitted map[rune]bool) {
	ni.PermittedValues = append(ni.PermittedValues, permitted)
}

// String returns a string representation of the selector.
func (ni NotIn) String() string {
	return fmt.Sprintf("%s notin (%s)", ni.Key, strings.Join(ni.Values, ", "))
}
