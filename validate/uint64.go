/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import "github.com/blend/go-sdk/ex"

// Uint64 errors
const (
	ErrUint64Min     ex.Class = "uint64 should be above a minimum value"
	ErrUint64Max     ex.Class = "uint64 should be below a maximum value"
	ErrUint64Zero    ex.Class = "uint64 should be zero"
	ErrUint64NotZero ex.Class = "uint64 should not be zero"
)

// Uint64 returns validators for uint64s.
func Uint64(value *uint64) Uint64Validators {
	return Uint64Validators{value}
}

// Uint64Validators implements uint64 validators.
type Uint64Validators struct {
	Value *uint64
}

// Min returns a validator that an uint64 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Uint64Validators) Min(min uint64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint64Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint64Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an uint64 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Uint64Validators) Max(max uint64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrUint64Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an uint64 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Uint64Validators) Between(min, max uint64) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint64Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint64Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrUint64Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Zero returns a validator that an uint64 is zero.
func (i Uint64Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrUint64Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrUint64Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an uint64 is not zero.
func (i Uint64Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrUint64NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrUint64NotZero, *i.Value)
		}
		return nil
	}
}
