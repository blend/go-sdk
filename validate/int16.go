/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import "github.com/blend/go-sdk/ex"

// Int16 errors
const (
	ErrInt16Min		ex.Class	= "int16 should be above a minimum value"
	ErrInt16Max		ex.Class	= "int16 should be below a maximum value"
	ErrInt16Positive	ex.Class	= "int16 should be positive"
	ErrInt16Negative	ex.Class	= "int16 should be negative"
	ErrInt16Zero		ex.Class	= "int16 should be zero"
	ErrInt16NotZero		ex.Class	= "int16 should not be zero"
)

// Int16 returns validators for int16s.
func Int16(value *int16) Int16Validators {
	return Int16Validators{value}
}

// Int16Validators implements int16 validators.
type Int16Validators struct {
	Value *int16
}

// Min returns a validator that an int16 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Int16Validators) Min(min int16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt16Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt16Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an int16 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Int16Validators) Max(max int16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrInt16Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int16 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Int16Validators) Between(min, max int16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt16Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt16Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrInt16Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int16 is positive.
func (i Int16Validators) Positive() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be positive
			return Error(ErrInt16Positive, nil)
		}
		if *i.Value < 0 {
			return Error(ErrInt16Positive, *i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int16 is negative.
func (i Int16Validators) Negative() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be negative
			return Error(ErrInt16Negative, nil)
		}
		if *i.Value > 0 {
			return Error(ErrInt16Negative, *i.Value)
		}
		return nil
	}
}

// Zero returns a validator that an int16 is zero.
func (i Int16Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrInt16Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrInt16Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an int16 is not zero.
func (i Int16Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrInt16NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrInt16NotZero, *i.Value)
		}
		return nil
	}
}
