package configutil

// ConfigResolver is an interface configs can implement to do basic config operations.
//
// DEPRECATION(1.2021*): this interface is deprecated and will be removed before 2021.
type ConfigResolver interface {
	Resolve() error
}

// ReturnFirst returns the first non-nil error in a list.
//
// DEPRECATION(1.2021*): this function is deprecated and will be removed before 2021.
func ReturnFirst(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

// AnyError returns the first non-nil error in a list.
//
// DEPRECATION(1.2021*): this function is deprecated and will be removed before 2021.
func AnyError(errors ...error) error {
	return ReturnFirst(errors...)
}
