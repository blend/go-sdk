package validate

import "github.com/blend/go-sdk/ex"

// Uint32 errors
const (
	ErrUint32Min     ex.Class = "uint32 should be above a minimum value"
	ErrUint32Max     ex.Class = "uint32 should be below a maximum value"
	ErrUint32Zero    ex.Class = "uint32 should be zero"
	ErrUint32NotZero ex.Class = "uint32 should not be zero"
)

// Uint32 returns validators for uint32s.
func Uint32(value *uint32) Uint32Validators {
	return Uint32Validators{value}
}

// Uint32Validators implements uint32 validators.
type Uint32Validators struct {
	Value *uint32
}

// Min returns a validator that an uint32 is above a minimum value inclusive.
// Min will pass for a value 1 if the min is set to 1, that is no error
// would be returned.
func (i Uint32Validators) Min(min uint32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint32Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint32Min, *i.Value, "min: %d", min)
		}
		return nil
	}
}

// Max returns a validator that an uint32 is below a max value inclusive.
// Max will pass for a value 10 if the max is set to 10, that is no error
// would be returned.
func (i Uint32Validators) Max(max uint32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot _violate_ a maximum because it has no value.
			return nil
		}
		if *i.Value > max {
			return Errorf(ErrUint32Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Between returns a validator that an uint32 is between a given min and max inclusive,
// that is, `.Between(1,5)` will _fail_ for [0] and [6] respectively, but pass
// for [1] and [5].
func (i Uint32Validators) Between(min, max uint32) Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot satisfy a minimum because it has no value.
			return Errorf(ErrUint32Min, nil, "min: %d", min)
		}
		if *i.Value < min {
			return Errorf(ErrUint32Min, *i.Value, "min: %d", min)
		}
		if *i.Value > max {
			return Errorf(ErrUint32Max, *i.Value, "max: %d", max)
		}
		return nil
	}
}

// Zero returns a validator that an uint32 is zero.
func (i Uint32Validators) Zero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be zero
			return Error(ErrUint32Zero, nil)
		}
		if *i.Value != 0 {
			return Error(ErrUint32Zero, *i.Value)
		}
		return nil
	}
}

// NotZero returns a validator that an uint32 is not zero.
func (i Uint32Validators) NotZero() Validator {
	return func() error {
		if i.Value == nil {
			// an unset value cannot be not zero
			return Error(ErrUint32NotZero, nil)
		}
		if *i.Value == 0 {
			return Error(ErrUint32NotZero, *i.Value)
		}
		return nil
	}
}
