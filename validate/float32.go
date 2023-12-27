/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"math"

	"github.com/blend/go-sdk/ex"
)

// Float32 errors
const (
	ErrFloat32Min		ex.Class	= "float32 should be above a minimum value"
	ErrFloat32Max		ex.Class	= "float32 should be below a maximum value"
	ErrFloat32Positive	ex.Class	= "float32 should be positive"
	ErrFloat32Negative	ex.Class	= "float32 should be negative"
	ErrFloat32Epsilon	ex.Class	= "float32 should be within an epsilon of a value"
	ErrFloat32Zero		ex.Class	= "float32 should be zero"
	ErrFloat32NotZero	ex.Class	= "float32 should not be zero"
)

// Float32 returns validators for float32s.
func Float32(value *float32) Float32Validators {
	return Float32Validators{value}
}

// Float32Validators implements float32 validators.
type Float32Validators struct {
	Value *float32
}

// Min returns a validator that a float32 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (f Float32Validators) Min(min float32) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrFloat32Min, nil, "min: %v", min)
		}
		if *f.Value < min {
			return Errorf(ErrFloat32Min, *f.Value, "min: %v", min)
		}
		return nil
	}
}

// Max returns a validator that a float32 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (f Float32Validators) Max(max float32) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *f.Value > max {
			return Errorf(ErrFloat32Max, *f.Value, "max: %v", max)
		}
		return nil
	}
}

// Between returns a validator that a float32 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (f Float32Validators) Between(min, max float32) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrFloat32Min, nil, "min: %v", min)
		}
		if *f.Value < min {
			return Errorf(ErrFloat32Min, *f.Value, "min: %v", min)
		}
		if *f.Value > max {
			return Errorf(ErrFloat32Max, *f.Value, "max: %v", max)
		}
		return nil
	}
}

// Positive returns a validator that a float32 is positive.
func (f Float32Validators) Positive() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be positive
			return Error(ErrFloat32Positive, nil)
		}
		if *f.Value < 0 {
			return Error(ErrFloat32Positive, *f.Value)
		}
		return nil
	}
}

// Negative returns a validator that a float32 is negative.
func (f Float32Validators) Negative() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be negative
			return Error(ErrFloat32Negative, nil)
		}
		if *f.Value > 0 {
			return Error(ErrFloat32Negative, *f.Value)
		}
		return nil
	}
}

// Epsilon returns if a value is comparable to another value within an epsilon.
// It will return a failure if the absolute difference between the target value
// and a given value is greater than the given epsilon.
func (f Float32Validators) Epsilon(value, epsilon float32) Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be negative
			return Errorf(ErrFloat32Epsilon, nil, "value: %v epsilon: %v", value, epsilon)
		}
		if math.Abs(float64(*f.Value-value)) > float64(epsilon) {
			return Errorf(ErrFloat32Epsilon, *f.Value, "value: %v epsilon: %v", value, epsilon)
		}
		return nil
	}
}

// Zero returns a validator that a float32 is zero.
func (f Float32Validators) Zero() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be zero
			return Error(ErrFloat32Zero, nil)
		}
		if *f.Value != 0 {
			return Error(ErrFloat32Zero, *f.Value)
		}
		return nil
	}
}

// NotZero returns a validator that a float32 is not zero.
func (f Float32Validators) NotZero() Validator {
	return func() error {
		if f.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrFloat32NotZero, nil)
		}
		if *f.Value == 0 {
			return Error(ErrFloat32NotZero, *f.Value)
		}
		return nil
	}
}
