package validate

import "github.com/blend/go-sdk/ex"

// Int errors
const (
	ErrIntMin      ex.Class = "int should be above a minimum value"
	ErrIntMax      ex.Class = "int should be below a maximum value"
	ErrIntPositive ex.Class = "int should be positive"
	ErrIntNegative ex.Class = "int should be negative"
)

// Int contains helpers for int validation.
func Int(value *int) IntValidators {
	return IntValidators{value}
}

// IntValidators implements int validators.
type IntValidators struct {
	Value *int
}

// Min returns a validator that an int is above a minimum value.
func (i IntValidators) Min(min int) Validator {
	return func() error {
		if i.Value == nil || *i.Value < min {
			return Errorf(ErrIntMin, i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that a int is below a max value.
func (i IntValidators) Max(max int) Validator {
	return func() error {
		if i.Value == nil {
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrIntMax, i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int is between a given min and max exclusive.
func (i IntValidators) Between(min, max int) Validator {
	return func() error {
		if i.Value == nil || *i.Value < min {
			return Errorf(ErrIntMin, i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrIntMax, i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int is positive.
func (i IntValidators) Positive() Validator {
	return func() error {
		if i.Value == nil || *i.Value < 0 {
			return Error(ErrIntPositive, i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int is negative.
func (i IntValidators) Negative() Validator {
	return func() error {
		if i.Value == nil || *i.Value > 0 {
			return Error(ErrIntNegative, i.Value)
		}
		return nil
	}
}
