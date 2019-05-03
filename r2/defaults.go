package r2

// Defaults is a helper to create requests with a common set of base options.
type Defaults []Option

// Add adds new options to the default set.
func (d Defaults) Add(options ...Option) Defaults {
	return append(d, options...)
}
