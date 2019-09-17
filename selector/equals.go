package selector

import "fmt"

// Equals returns if a key strictly equals a value.
type Equals struct {
	Key, Value      string
	PermittedValues []map[rune]bool
}

// Matches returns the selector result.
func (e Equals) Matches(labels Labels) bool {
	if value, hasValue := labels[e.Key]; hasValue {
		return e.Value == value
	}
	return false
}

// Validate validates the selector.
func (e Equals) Validate(options ...SelectorOption) (err error) {
	var selector Selector = &e
	for _, option := range options {
		option(selector)
	}

	err = CheckKey(e.Key)
	if err != nil {
		return
	}
	err = CheckValue(e.Value, e.PermittedValues...)
	return
}

// AddPermittedValues adds runes to be accepted in values
func (e *Equals) AddPermittedValues(permitted map[rune]bool) {
	e.PermittedValues = append(e.PermittedValues, permitted)
}

// String returns the string representation of the selector.
func (e Equals) String() string {
	return fmt.Sprintf("%s == %s", e.Key, e.Value)
}
