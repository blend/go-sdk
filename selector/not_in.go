package selector

import (
	"fmt"
	"strings"
)

// NotIn returns if a key does not match a set of values.
type NotIn struct {
	Key    string
	Values []string
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
func (ni NotIn) Validate(permittedValues ...map[rune]bool) (err error) {
	err = CheckKey(ni.Key)
	if err != nil {
		return
	}
	for _, v := range ni.Values {
		err = CheckValue(v, permittedValues...)
		if err != nil {
			return
		}
	}
	return
}

// String returns a string representation of the selector.
func (ni NotIn) String() string {
	return fmt.Sprintf("%s notin (%s)", ni.Key, strings.Join(ni.Values, ", "))
}
