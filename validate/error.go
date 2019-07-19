package validate

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
)

// The root error, all validation errors inherit from this type.
const (
	ErrValidation ex.Class = "validation error"
)

// Error returns a new validation error.
// The root class of the error will be ErrValidation.
// The root stack will begin the frame above this call to error.
// The inner error will the cause of the validation vault.
func Error(cause error, args ...interface{}) error {
	return &ex.Ex{
		Class:   ErrValidation,
		Message: fmt.Sprint(args...),
		Inner:   cause,
		Stack:   ex.Callers(ex.DefaultNewStartDepth + 1),
	}
}

// Errorf returns a new validation error.
// The root class of the error will be ErrValidation.
// The root stack will begin the frame above this call to error.
// The inner error will the cause of the validation vault.
func Errorf(cause error, format string, args ...interface{}) error {
	return &ex.Ex{
		Class:   ErrValidation,
		Message: fmt.Sprintf(format, args...),
		Inner:   cause,
		Stack:   ex.Callers(ex.DefaultNewStartDepth + 1),
	}
}

// Format formats an error.
func Format(err error) string {
	if err == nil {
		return "ok!"
	}
	class := ex.ErrClass(err)
	message := ex.ErrMessage(err)
	inner := ex.ErrInner(err)
	innerClass := ex.ErrClass(inner)
	if message != "" {
		return fmt.Sprintf("%v; %v (%v)", class, innerClass, message)
	}
	return fmt.Sprintf("%v; %v", class, innerClass)
}

// Is returns if an error is a validation error.
func Is(err error) bool {
	return ex.Is(err, ErrValidation)
}
