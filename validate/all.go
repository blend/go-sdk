package validate

// All returns a validator that returns all non-nil validation errors from a given
// set of validators.
func All(validators ...Validator) Validator {
	return func() error {
		var output []error
		var err error
		for _, validator := range validators {
			if err = validator(); err != nil {
				output = append(output, err)
			}
		}
		if len(output) > 0 {
			return ValidationErrors(output)
		}
		return nil
	}
}
