package validate

import (
	"math"

	"github.com/blend/go-sdk/ex"
)

// Float64 errors
const (
	ErrFloat64Min      ex.Class = "float64 should be above a minimum value"
	ErrFloat64Max      ex.Class = "float64 should be below a maximum value"
	ErrFloat64Positive ex.Class = "float64 should be positive"
	ErrFloat64Negative ex.Class = "float64 should be negative"
	ErrFloat64Epsilon  ex.Class = "float64 should be within an epsilon of a value"
)

// Float64 returns validators for float64s.
func Float64(value *float64) Float64Validators {
	return Float64Validators{value}
}

// Float64Validators implements float64 validators.
type Float64Validators struct {
	Value *float64
}

// Min returns a validator that a float64 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (f Float64Validators) Min(min float64) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrFloat64Min, nil, "min: %v", min)
		}
		if *f.Value < min {
			return Errorf(ErrFloat64Min, *f.Value, "min: %v", min)
		}
		return nil
	}
}

// Max returns a validator that a float64 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (f Float64Validators) Max(max float64) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *f.Value > max {
			return Errorf(ErrFloat64Max, *f.Value, "max: %v", max)
		}
		return nil
	}
}

// Between returns a validator that a float64 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (f Float64Validators) Between(min, max float64) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrFloat64Min, nil, "min: %v", min)
		}
		if *f.Value < min {
			return Errorf(ErrFloat64Min, *f.Value, "min: %v", min)
		}
		if *f.Value > max {
			return Errorf(ErrFloat64Max, *f.Value, "max: %v", max)
		}
		return nil
	}
}

// Positive returns a validator that a float64 is positive.
func (f Float64Validators) Positive() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be positive
			return Error(ErrFloat64Positive, nil)
		}
		if *f.Value < 0 {
			return Error(ErrFloat64Positive, *f.Value)
		}
		return nil
	}
}

// Negative returns a validator that a float64 is negative.
func (f Float64Validators) Negative() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be negative
			return Error(ErrFloat64Negative, nil)
		}
		if *f.Value > 0 {
			return Error(ErrFloat64Negative, *f.Value)
		}
		return nil
	}
}

// Epsilon returns if a value is comparable to another value within an epsilon.
// It will return a failure if the absolute difference between the target value
// and a given value is greater than the given epsilon.
func (f Float64Validators) Epsilon(value, epsilon float64) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be negative
			return Errorf(ErrFloat64Epsilon, nil, "value: %v epsilon: %v", value, epsilon)
		}
		if math.Abs(*f.Value-value) > epsilon {
			return Errorf(ErrFloat64Epsilon, *f.Value, "value: %v epsilon: %v", value, epsilon)
		}
		return nil
	}
}
