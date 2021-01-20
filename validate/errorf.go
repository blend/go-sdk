package validate

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
)

// Errorf returns a new validation error.
// The root class of the error will be ErrValidation.
// The root stack will begin the frame above this call to error.
// The inner error will the cause of the validation vault.
func Errorf(cause error, value interface{}, format string, args ...interface{}) error {
	return &ex.Ex{
		Class: ErrValidation,
		Inner: &ValidationError{
			Cause:   cause,
			Value:   value,
			Message: fmt.Sprintf(format, args...),
		},
		StackTrace: ex.Callers(ex.DefaultNewStartDepth + 1),
	}
}
