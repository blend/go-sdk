package validate

import "github.com/blend/go-sdk/ex"

// Int errors
const (
	ErrIntMin      ex.Class = "int should be above a minimum value"
	ErrIntMax      ex.Class = "int should be below a maximum value"
	ErrIntPositive ex.Class = "int should be positive"
	ErrIntNegative ex.Class = "int should be negative"
)

// Int validator singleton.
func Int(v int) IntValidators {
	return IntValidators(v)
}

// IntValidators contains helpers for int validation.
type IntValidators int

// Min returns a validator that an int is above a minimum value.
func (i IntValidators) Min(min int) Validator {
	return func() error {
		if int(i) < min {
			return Errorf(ErrIntMin, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that a int is below a max value.
func (i IntValidators) Max(max int) Validator {
	return func() error {
		if int(i) > max {
			return Errorf(ErrIntMax, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int is between a given min and max exclusive.
func (i IntValidators) Between(min, max int) Validator {
	return func() error {
		if int(i) < min {
			return Errorf(ErrIntMin, "min: %d", min)
		}
		if int(i) > max {
			return Errorf(ErrIntMax, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int is positive.
func (i IntValidators) Positive() error {
	if int(i) < 0 {
		return Error(ErrIntPositive)
	}
	return nil
}

// Negative returns a validator that an int is negative.
func (i IntValidators) Negative() error {
	if int(i) > 0 {
		return Error(ErrIntNegative)
	}
	return nil
}
