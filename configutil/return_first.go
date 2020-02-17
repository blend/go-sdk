package configutil

// ReturnFirst returns the first non-nil error in a list.
func ReturnFirst(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
