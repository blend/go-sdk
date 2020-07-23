package validate

import "github.com/blend/go-sdk/ex"

// Uint8 errors
const (
	ErrUint8Min     ex.Class = "uint8 should be above a minimum value"
	ErrUint8Max     ex.Class = "uint8 should be below a maximum value"
	ErrUint8Zero    ex.Class = "uint8 should be zero"
	ErrUint8NotZero ex.Class = "uint8 should not be zero"
)

// Uint8 returns validators for uint8s.
func Uint8(value *uint8) Uint8Validators {
	return Uint8Validators{value}
}

// Uint8Validators implements uint8 validators.
type Uint8Validators struct {
	Value *uint8
}

// Min returns a validator that an uint8 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Uint8Validators) Min(min uint8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint8Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint8Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an uint8 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Uint8Validators) Max(max uint8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrUint8Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an uint8 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Uint8Validators) Between(min, max uint8) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint8Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint8Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrUint8Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Zero returns a validator that an uint8 is zero.
func (i Uint8Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrUint8Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrUint8Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an uint8 is not zero.
func (i Uint8Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrUint8NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrUint8NotZero, *i.Value)
		}
		return nil
	}
}
