package validate

import "github.com/blend/go-sdk/ex"

// Int errors
const (
	ErrIntMin      ex.Class = "int should be above a minimum value"
	ErrIntMax      ex.Class = "int should be below a maximum value"
	ErrIntPositive ex.Class = "int should be positive"
	ErrIntNegative ex.Class = "int should be negative"
)

// IntValidator is a validator for strings.
type IntValidator func(int) error

// Int validator singleton.
var Int intValidators

// intValidators contains helpers for int validation.
type intValidators struct{}

// Min returns a validator that an int is above a minimum value.
func (i intValidators) Min(min int) IntValidator {
	return func(v int) error {
		if v < min {
			return Errorf(ErrIntMin, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that a int is below a max value.
func (i intValidators) Max(max int) IntValidator {
	return func(v int) error {
		if v > max {
			return Errorf(ErrIntMax, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int is between a given min and max exclusive.
func (i intValidators) Between(min, max int) IntValidator {
	return func(v int) error {
		if v < min {
			return Errorf(ErrIntMin, "min: %d", min)
		}
		if v > max {
			return Errorf(ErrIntMax, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int is positive.
func (i intValidators) Positive(v int) error {
	if v < 0 {
		return Error(ErrIntPositive)
	}
	return nil
}

// Negative returns a validator that an int is negative.
func (i intValidators) Negative(v int) error {
	if v > 0 {
		return Error(ErrIntNegative)
	}
	return nil
}
