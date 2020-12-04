package validate

import "github.com/blend/go-sdk/ex"

// Int8 errors
const (
	ErrInt8Min      ex.Class = "int8 should be above a minimum value"
	ErrInt8Max      ex.Class = "int8 should be below a maximum value"
	ErrInt8Positive ex.Class = "int8 should be positive"
	ErrInt8Negative ex.Class = "int8 should be negative"
	ErrInt8Zero     ex.Class = "int8 should be zero"
	ErrInt8NotZero  ex.Class = "int8 should not be zero"
)

// Int8 returns validators for int8s.
func Int8(value *int8) Int8Validators {
	return Int8Validators{value}
}

// Int8Validators implements int8 validators.
type Int8Validators struct {
	Value *int8
}

// Min returns a validator that an int8 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Int8Validators) Min(min int8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt8Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt8Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an int8 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Int8Validators) Max(max int8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrInt8Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int8 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Int8Validators) Between(min, max int8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt8Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt8Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrInt8Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int8 is positive.
func (i Int8Validators) Positive() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be positive
			return Error(ErrInt8Positive, nil)
		}
		if *i.Value < 0 {
			return Error(ErrInt8Positive, *i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int8 is negative.
func (i Int8Validators) Negative() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be negative
			return Error(ErrInt8Negative, nil)
		}
		if *i.Value > 0 {
			return Error(ErrInt8Negative, *i.Value)
		}
		return nil
	}
}

// Zero returns a validator that an int8 is zero.
func (i Int8Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrInt8Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrInt8Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an int8 is not zero.
func (i Int8Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrInt8NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrInt8NotZero, *i.Value)
		}
		return nil
	}
}
