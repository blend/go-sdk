package selector

import "strings"

// And is a combination selector.
type And []Selector

// Matches returns if both A and B match the labels.
func (a And) Matches(labels Labels) bool {
	for _, s := range a {
		if !s.Matches(labels) {
			return false
		}
	}
	return true
}

// Validate validates all the selectors in the clause.
func (a And) Validate(permittedValues ...map[rune]bool) (err error) {
	for _, s := range a {
		err = s.Validate(permittedValues...)
		if err != nil {
			return
		}
	}
	return
}

// And returns a string representation for the selector.
func (a And) String() string {
	var childValues []string
	for _, c := range a {
		childValues = append(childValues, c.String())
	}
	return strings.Join(childValues, ", ")
}
