/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import "github.com/blend/go-sdk/ex"

// Int64 errors
const (
	ErrInt64Min      ex.Class = "int64 should be above a minimum value"
	ErrInt64Max      ex.Class = "int64 should be below a maximum value"
	ErrInt64Positive ex.Class = "int64 should be positive"
	ErrInt64Negative ex.Class = "int64 should be negative"
	ErrInt64Zero     ex.Class = "int64 should be zero"
	ErrInt64NotZero  ex.Class = "int64 should not be zero"
)

// Int64 returns validators for int64s.
func Int64(value *int64) Int64Validators {
	return Int64Validators{value}
}

// Int64Validators implements int64 validators.
type Int64Validators struct {
	Value *int64
}

// Min returns a validator that an int64 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Int64Validators) Min(min int64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt64Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt64Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an int64 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Int64Validators) Max(max int64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrInt64Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int64 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Int64Validators) Between(min, max int64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt64Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt64Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrInt64Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int64 is positive.
func (i Int64Validators) Positive() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be positive
			return Error(ErrInt64Positive, nil)
		}
		if *i.Value < 0 {
			return Error(ErrInt64Positive, *i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int64 is negative.
func (i Int64Validators) Negative() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be negative
			return Error(ErrInt64Negative, nil)
		}
		if *i.Value > 0 {
			return Error(ErrInt64Negative, *i.Value)
		}
		return nil
	}
}

// Zero returns a validator that an int64 is zero.
func (i Int64Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrInt64Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrInt64Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an int64 is not zero.
func (i Int64Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrInt64NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrInt64NotZero, *i.Value)
		}
		return nil
	}
}
