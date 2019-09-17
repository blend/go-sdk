package selector

// Any matches everything
type Any struct{}

// Matches returns true
func (a Any) Matches(labels Labels) bool {
	return true
}

// Validate validates the selector
func (a Any) Validate(options ...SelectorOption) (err error) {
	return nil
}

// AddPermittedValues adds runes to be accepted in values
func (a Any) AddPermittedValues(permitted map[rune]bool) {
	return
}

// String returns a string representation of the selector
func (a Any) String() string {
	return ""
}
