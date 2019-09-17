package selector

// Option is a tweak to selector parsing.
type Option func(p *Parser)

// SkipValidation is an option to skip checking the values of selector expressions.
func SkipValidation(p *Parser) {
	p.skipValidation = true
}

// OptPermittedValues is an option to extend the set of symbols that are valid in name values
func OptPermittedValues(permitted []rune) Option {
	return func(p *Parser) {
		if p.permittedValues == nil {
			p.permittedValues = map[rune]bool{}
		}

		for _, r := range permitted {
			p.permittedValues[r] = true
		}

		p.selectorOptions = append(p.selectorOptions, SelectorOptPermittedValues(permitted...))
	}
}
