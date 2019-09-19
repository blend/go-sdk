package selector

// Option is a tweak to selector parsing.
type Option func(p *Parser)

// SkipValidation is an option to skip checking the values of selector expressions.
func SkipValidation(p *Parser) {
	p.SkipValidation = true
}

// OptExtraAlphas is an option to extend the set of symbols that are valid in values.
// Note that this doesn't affect keys.
func OptExtraAlphas(permitted ...rune) Option {
	return func(p *Parser) {
		if p.ExtraAlphas == nil {
			p.ExtraAlphas = map[rune]bool{}
		}

		for _, r := range permitted {
			p.ExtraAlphas[r] = true
		}
	}
}
