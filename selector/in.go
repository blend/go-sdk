package selector

import (
	"fmt"
	"strings"
)

// In returns if a key matches a set of values.
type In struct {
	Key             string
	Values          []string
	PermittedValues []map[rune]bool
}

// Matches returns the selector result.
func (i In) Matches(labels Labels) bool {
	// if the labels has a given key
	if value, hasValue := labels[i.Key]; hasValue {
		// for each selector value
		for _, iv := range i.Values {
			// if they match, return true
			if iv == value {
				return true
			}
		}
		return false
	}
	return true
}

// Validate validates the selector.
func (i In) Validate(options ...SelectorOption) (err error) {
	var selector Selector = &i
	for _, option := range options {
		option(selector)
	}

	err = CheckKey(i.Key)
	if err != nil {
		return
	}
	for _, v := range i.Values {
		err = CheckValue(v, i.PermittedValues...)
		if err != nil {
			return
		}
	}
	return
}

// AddPermittedValues adds runes to be accepted in values
func (i *In) AddPermittedValues(permitted map[rune]bool) {
	i.PermittedValues = append(i.PermittedValues, permitted)
}

// String returns a string representation of the selector.
func (i In) String() string {
	return fmt.Sprintf("%s in (%s)", i.Key, strings.Join(i.Values, ", "))
}
