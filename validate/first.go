package validate

// First returns the first set error if one exists.
func First(validators ...func() error) error {
	var err error
	for _, validator := range validators {
		if err = validator(); err != nil {
			return err
		}
	}
	return nil
}
