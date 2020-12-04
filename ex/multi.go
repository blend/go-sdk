package ex

import (
	"fmt"
	"strings"
)

// Append appends errors together, creating a multi-error.
func Append(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	var all []error
	if err != nil {
		all = append(all, NewWithStackDepth(err, DefaultNewStartDepth+1))
	}
	for _, extra := range errs {
		if extra != nil {
			all = append(all, NewWithStackDepth(extra, DefaultNewStartDepth+1))
		}
	}
	if len(all) == 0 {
		return nil
	}
	if len(all) == 1 {
		return all[0]
	}
	return Multi(all)
}

// Unwrap unwraps multi-errors.
func Unwrap(err error) []error {
	if typed, ok := err.(Multi); ok {
		return []error(typed)
	}
	return []error{err}
}

// Multi represents an array of errors.
type Multi []error

// Error implements error.
func (m Multi) Error() string {
	if len(m) == 0 {
		return ""
	}
	if len(m) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %v\n\n", m[0])
	}

	points := make([]string, len(m))
	for i, err := range m {
		points[i] = fmt.Sprintf("* %v", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(m), strings.Join(points, "\n\t"))
}

// WrappedErrors implements something in errors.
func (m Multi) WrappedErrors() []error {
	return m
}

// Unwrap returns an error from Error (or nil if there are no errors).
// This error returned will further support Unwrap to get the next error,
// etc.
//
// The resulting error supports errors.As/Is/Unwrap so you can continue
// to use the stdlib errors package to introspect further.
//
// This will perform a shallow copy of the errors slice. Any errors appended
// to this error after calling Unwrap will not be available until a new
// Unwrap is called on the multierror.Error.
func (m Multi) Unwrap() error {
	if m == nil || len(m) == 0 {
		return nil
	}
	if len(m) == 1 {
		return m[0]
	}
	errs := make([]error, len(m))
	copy(errs, m)
	return Nest(errs...)
}
