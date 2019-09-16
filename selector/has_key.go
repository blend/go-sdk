package selector

// HasKey returns if a label set has a given key.
type HasKey string

// Matches returns the selector result.
func (hk HasKey) Matches(labels Labels) bool {
	_, hasKey := labels[string(hk)]
	return hasKey
}

// Validate validates the selector.
func (hk HasKey) Validate(permittedValues ...map[rune]bool) (err error) {
	err = CheckKey(string(hk))
	return
}

// String returns a string representation of the selector.
func (hk HasKey) String() string {
	return string(hk)
}
