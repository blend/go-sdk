package validate

// When returns the result of the "passes" validator if the predicate returns true,
// otherwise it returns the result of the "fails" validator.
func When(predicate func() bool, passes, fails Validator) Validator {
	return func() error {
		if predicate() {
			return passes()
		}
		return fails()
	}
}
