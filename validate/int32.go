/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import "github.com/blend/go-sdk/ex"

// Int32 errors
const (
	ErrInt32Min		ex.Class	= "int32 should be above a minimum value"
	ErrInt32Max		ex.Class	= "int32 should be below a maximum value"
	ErrInt32Positive	ex.Class	= "int32 should be positive"
	ErrInt32Negative	ex.Class	= "int32 should be negative"
	ErrInt32Zero		ex.Class	= "int32 should be zero"
	ErrInt32NotZero		ex.Class	= "int32 should not be zero"
)

// Int32 returns validators for int32s.
func Int32(value *int32) Int32Validators {
	return Int32Validators{value}
}

// Int32Validators implements int32 validators.
type Int32Validators struct {
	Value *int32
}

// Min returns a validator that an int32 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Int32Validators) Min(min int32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt32Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt32Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an int32 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Int32Validators) Max(max int32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrInt32Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an int32 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Int32Validators) Between(min, max int32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrInt32Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrInt32Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrInt32Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Positive returns a validator that an int32 is positive.
func (i Int32Validators) Positive() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be positive
			return Error(ErrInt32Positive, nil)
		}
		if *i.Value < 0 {
			return Error(ErrInt32Positive, *i.Value)
		}
		return nil
	}
}

// Negative returns a validator that an int32 is negative.
func (i Int32Validators) Negative() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be negative
			return Error(ErrInt32Negative, nil)
		}
		if *i.Value > 0 {
			return Error(ErrInt32Negative, *i.Value)
		}
		return nil
	}
}

// Zero returns a validator that an int32 is zero.
func (i Int32Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrInt32Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrInt32Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an int32 is not zero.
func (i Int32Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrInt32NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrInt32NotZero, *i.Value)
		}
		return nil
	}
}
