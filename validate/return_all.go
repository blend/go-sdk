package validate

// ReturnAll returns all the failing validations as an
// array of validation errors.
func ReturnAll(validators ...Validator) []error {
	var output []error
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			output = append(output, err)
		}
	}
	return output
}
