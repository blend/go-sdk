package validate

import "github.com/blend/go-sdk/ex"

// Int errors
const (
	ErrIntMin      ex.Class = "int should be above a minimum value"
	ErrIntMax      ex.Class = "int should be below a maximum value"
	ErrIntPositive ex.Class = "int should be positive"
	ErrIntNegative ex.Class = "int should be negative"
)

// Int returns validators for ints.
func Int(value *int) IntValidators {
	return IntValidators{value}
}

// IntValidators implements int validators.
type IntValidators struct {
	Value *int
}

// Min returns a validator that an int is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i IntValidators) Min(min int) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrIntMin, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrIntMin, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that a int is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i IntValidators) Max(max int) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrIntMax, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i IntValidators) Between(min, max int) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrIntMin, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrIntMin, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrIntMax, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int is positive.
func (i IntValidators) Positive() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be positive
			return Error(ErrIntPositive, nil)
		}
		if *i.Value < 0 {
			return Error(ErrIntPositive, *i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int is negative.
func (i IntValidators) Negative() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be negative
			return Error(ErrIntNegative, nil)
		}
		if *i.Value > 0 {
			return Error(ErrIntNegative, *i.Value)
		}
		return nil
	}
}
