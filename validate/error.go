package validate

import (
	"fmt"

	"github.com/blend/go-sdk/ex"
)

// The root error, all validation errors inherit from this type.
const (
	ErrValidation ex.Class = "validation error"
)

var (
	_ ex.ClassProvider = (*ValidationError)(nil)
)

// ValidationError is the inner error for validation exceptions.
type ValidationError struct {
	Cause   error
	Message string
	Value   interface{}
}

// Class implements
func (ve ValidationError) Class() error {
	return ve.Cause
}

// Error implements error.
func (ve ValidationError) Error() string {
	if ve.Value != nil && ve.Message != "" {
		return fmt.Sprintf("%v; %v [%v]", ve.Cause, ve.Message, ve.Value)
	}
	if ve.Value != nil {
		return fmt.Sprintf("%v [%v]", ve.Cause, ve.Value)
	}
	if ve.Message != "" {
		return fmt.Sprintf("%v; %v", ve.Cause, ve.Message)
	}
	return ve.Cause.Error()
}

// Error returns a new validation error.
// The root class of the error will be ErrValidation.
// The root stack will begin the frame above this call to error.
// The inner error will the cause of the validation vault.
func Error(cause error, value interface{}, messageArgs ...interface{}) error {
	return &ex.Ex{
		Class: ErrValidation,
		Inner: &ValidationError{
			Cause:   cause,
			Value:   value,
			Message: fmt.Sprint(messageArgs...),
		},
		StackTrace: ex.Callers(ex.DefaultNewStartDepth + 1),
	}
}

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

// Inner returns the inner validation error if it's present on
// the outer error.
func Inner(err error) *ValidationError {
	inner := ex.ErrInner(err)
	if inner == nil {
		return nil
	}
	if typed, ok := inner.(*ValidationError); ok {
		return typed
	}
	return nil
}

// Cause returns the underlying validation failure for an error.
// If the error is not a validation error, it returns the error class.
func Cause(err error) error {
	if exClass := ex.ErrClass(err); exClass != ErrValidation {
		return exClass
	}
	if inner := Inner(err); inner != nil {
		return inner.Cause
	}
	return nil
}

// Message returns the underlying validation error message.
func Message(err error) string {
	if inner := Inner(err); inner != nil {
		return inner.Message
	}
	return ""
}

// Value returns the validation error value.
func Value(err error) interface{} {
	if inner := Inner(err); inner != nil {
		return inner.Value
	}
	return nil
}

// Format formats an error.
func Format(err error) string {
	if err == nil {
		return "ok!"
	}
	return Inner(err).Error()
}

// Is returns if an error is a validation error.
func Is(err error) bool {
	return ex.Is(err, ErrValidation)
}
