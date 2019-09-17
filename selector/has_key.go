package selector

// HasKey returns if a label set has a given key.
type HasKey string

// Matches returns the selector result.
func (hk HasKey) Matches(labels Labels) bool {
	_, hasKey := labels[string(hk)]
	return hasKey
}

// Validate validates the selector.
func (hk HasKey) Validate(options ...SelectorOption) (err error) {
	err = CheckKey(string(hk))
	return
}

// AddPermittedValues satisfies the Selector interface
func (hk HasKey) AddPermittedValues(permitted map[rune]bool) {
	// since this selector only deals with keys, this method does nothing
	return
}

// String returns a string representation of the selector.
func (hk HasKey) String() string {
	return string(hk)
}
