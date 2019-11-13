package validate

// When returns the result of the "passes" validator if the predicate returns true,
// otherwise it returns nil.
func When(predicate func() bool, passes Validator) Validator {
	return func() error {
		if predicate() {
			return passes()
		}
		return nil
	}
}
