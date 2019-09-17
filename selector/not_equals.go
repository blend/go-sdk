package selector

import "fmt"

// NotEquals returns if a key strictly equals a value.
type NotEquals struct {
	Key, Value      string
	PermittedValues []map[rune]bool
}

// Matches returns the selector result.
func (ne NotEquals) Matches(labels Labels) bool {
	if value, hasValue := labels[ne.Key]; hasValue {
		return ne.Value != value
	}
	return true
}

// Validate validates the selector.
func (ne *NotEquals) Validate(options ...SelectorOption) (err error) {
	var selector Selector = ne
	for _, option := range options {
		option(selector)
	}

	err = CheckKey(ne.Key)
	if err != nil {
		return
	}
	err = CheckValue(ne.Value, ne.PermittedValues...)
	return
}

// AddPermittedValues adds runes to be accepted in values
func (ne *NotEquals) AddPermittedValues(permitted map[rune]bool) {
	ne.PermittedValues = append(ne.PermittedValues, permitted)
}

// String returns a string representation of the selector.
func (ne NotEquals) String() string {
	return fmt.Sprintf("%s != %s", ne.Key, ne.Value)
}
