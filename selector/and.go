package selector

import "strings"

// And is a combination selector.
type And struct {
	selectors []Selector
	options   []SelectorOption
}

// Matches returns if both A and B match the labels.
func (a And) Matches(labels Labels) bool {
	for _, s := range a.selectors {
		if !s.Matches(labels) {
			return false
		}
	}
	return true
}

// Validate validates all the selectors in the clause.
func (a And) Validate(options ...SelectorOption) (err error) {
	for _, s := range a.selectors {
		err = s.Validate(append(a.options, options...)...)
		if err != nil {
			return
		}
	}
	return
}

// AddPermittedValues adds runes to be accepted in values
func (a *And) AddPermittedValues(permitted map[rune]bool) {
	p := []rune{}
	for key := range permitted {
		p = append(p, key)
	}

	a.options = append(a.options, SelectorOptPermittedValues(p...))
}

// And returns a string representation for the selector.
func (a And) String() string {
	var childValues []string
	for _, c := range a.selectors {
		childValues = append(childValues, c.String())
	}
	return strings.Join(childValues, ", ")
}
