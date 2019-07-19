package validate

// All returns the all the failing validations.
func All(validators ...func() error) []error {
	var output []error
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			output = append(output, err)
		}
	}
	return output
}
