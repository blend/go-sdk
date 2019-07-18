package validate

// Any returns the first set error if one exists.
func Any(verrs ...error) error {
	for _, err := range verrs {
		if err != nil {
			return err
		}
	}
	return nil
}
