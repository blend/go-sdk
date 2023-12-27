/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import "github.com/blend/go-sdk/ex"

// Uint16 errors
const (
	ErrUint16Min		ex.Class	= "uint16 should be above a minimum value"
	ErrUint16Max		ex.Class	= "uint16 should be below a maximum value"
	ErrUint16Zero		ex.Class	= "uint16 should be zero"
	ErrUint16NotZero	ex.Class	= "uint16 should not be zero"
)

// Uint16 returns validators for uint16s.
func Uint16(value *uint16) Uint16Validators {
	return Uint16Validators{value}
}

// Uint16Validators implements uint16 validators.
type Uint16Validators struct {
	Value *uint16
}

// Min returns a validator that an uint16 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Uint16Validators) Min(min uint16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint16Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint16Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an uint16 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Uint16Validators) Max(max uint16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrUint16Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an uint16 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Uint16Validators) Between(min, max uint16) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint16Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint16Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrUint16Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Zero returns a validator that an uint16 is zero.
func (i Uint16Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrUint16Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrUint16Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an uint16 is not zero.
func (i Uint16Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrUint16NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrUint16NotZero, *i.Value)
		}
		return nil
	}
}
