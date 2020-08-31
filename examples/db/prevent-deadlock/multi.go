package main

import (
	"fmt"
	"strings"
)

// NOTE: Ensure that
//       * `multiError` satisfies `error`.
var (
	_ error = (*multiError)(nil)
)

type multiError struct {
	Errors []error
}

func (me *multiError) Error() string {
	if me == nil || len(me.Errors) == 0 {
		return "<nil>"
	}
	parts := []string{}
	for _, err := range me.Errors {
		parts = append(parts, fmt.Sprintf("- %#v", err))
	}
	return strings.Join(parts, "\n")
}

func nest(err1, err2 error) error {
	// We know below here that both errors are non-nil.
	asMulti1, ok1 := err1.(*multiError)
	asMulti2, ok2 := err2.(*multiError)

	if err1 == nil {
		if err2 == nil {
			return nil
		}
		if ok2 {
			return err2
		}
		return &multiError{Errors: []error{err2}}
	}

	if err2 == nil {
		if ok1 {
			return err1
		}
		return &multiError{Errors: []error{err1}}
	}

	if ok1 {
		if ok2 {
			return &multiError{Errors: append(asMulti1.Errors, asMulti2.Errors...)}
		}

		return &multiError{Errors: append(asMulti1.Errors, err2)}
	}

	if ok2 {
		return &multiError{Errors: append(asMulti2.Errors, err1)}
	}

	return &multiError{Errors: []error{err1, err2}}
}
