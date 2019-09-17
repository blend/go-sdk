package selector

// SelectorOption is a tweak to selector parsing.
type SelectorOption func(s Selector)

// SelectorOptPermittedValues is an option to extend the set of symbols that are valid in name values
func SelectorOptPermittedValues(permitted ...rune) SelectorOption {
	m := map[rune]bool{}
	for _, r := range permitted {
		m[r] = true
	}
	return func(s Selector) {
		s.AddPermittedValues(m)
	}
}
