package validate

// All returns all the failing validations.
func All(validators ...Validator) []error {
	var output []error
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			output = append(output, err)
		}
	}
	return output
}
