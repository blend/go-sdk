/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// Basic errors
const (
	ErrZero		ex.Class	= "object should be its default value or unset"
	ErrNotZero	ex.Class	= "object should not be its default value or unset"
	ErrRequired	ex.Class	= "field is required"
	ErrForbidden	ex.Class	= "field is forbidden"
	ErrEmpty	ex.Class	= "object should be empty"
	ErrNotEmpty	ex.Class	= "object should not be empty"
	ErrLen		ex.Class	= "object should have a given length"
	ErrNil		ex.Class	= "object should be nil"
	ErrNotNil	ex.Class	= "object should not be nil"
	ErrEquals	ex.Class	= "objects should be equal"
	ErrNotEquals	ex.Class	= "objects should not be equal"
	ErrAllowed	ex.Class	= "objects should be one of a given set of allowed values"
	ErrDisallowed	ex.Class	= "objects should not be one of a given set of disallowed values"
)

// Any returns a new AnyRefValidators.
// AnyRef can be used to validate any type, but will be more limited
// than using type specific validators.
func Any(obj interface{}) AnyValidators {
	return AnyValidators{Obj: obj}
}

// AnyValidators are validators for the empty interface{}.
type AnyValidators struct {
	Obj interface{}
}

// Forbidden mirrors Zero but uses a specific error.
// This is useful if you want to have more aggressive failure cases.
func (a AnyValidators) Forbidden() Validator {
	return func() error {
		if err := a.Zero()(); err != nil {
			return Error(ErrForbidden, a.Obj)
		}
		return nil
	}
}

// Required mirrors NotZero but uses a specific error.
// This is useful if you want to have more aggressive failure cases.
func (a AnyValidators) Required() Validator {
	return func() error {
		if err := a.NotZero()(); err != nil {
			return Error(ErrRequired, a.Obj)
		}
		return nil
	}
}

// Zero returns a validator that asserts an object is it's zero value.
// This nil for pointers, slices, maps, channels.
// And whatever equality passes for everything else with it's initialized value.
// Note: this method uses reflect.Zero, there are faster .Zero evaluators
// for the relevant numeric types.
func (a AnyValidators) Zero() Validator {
	return func() error {
		if a.Obj == nil {
			return nil
		}

		zero := reflect.Zero(reflect.TypeOf(a.Obj)).Interface()
		if verr := a.Equals(zero)(); verr == nil {
			return nil
		}
		return Error(ErrZero, a.Obj)
	}
}

// NotZero returns a validator that a given field is set.
// It will return an error if the field is unset.
// Note: this method uses reflect.Zero, there are faster .NotZero evaluators
// for the relevant numeric types.
func (a AnyValidators) NotZero() Validator {
	return func() error {
		if err := a.Zero()(); err == nil {
			return Error(ErrNotZero, a.Obj)
		}
		return nil
	}
}

// Empty returns if a slice, map or channel is empty.
// It will error if the object is not a slice, map or channel.
func (a AnyValidators) Empty() Validator {
	return func() error {
		objLen, err := GetLength(a.Obj)
		if err != nil {
			return err
		}
		if objLen == 0 {
			return nil
		}
		return Error(ErrEmpty, a.Obj)
	}
}

// NotEmpty returns if a slice, map or channel is not empty.
// It will error if the object is not a slice, map or channel.
func (a AnyValidators) NotEmpty() Validator {
	return func() error {
		objLen, err := GetLength(a.Obj)
		if err != nil {
			return err
		}
		if objLen > 0 {
			return nil
		}
		return Error(ErrNotEmpty, a.Obj)
	}
}

// Len validates the length is a given value.
func (a AnyValidators) Len(length int) Validator {
	return func() error {
		objLen, err := GetLength(a.Obj)
		if err != nil {
			return err
		}
		if objLen == length {
			return nil
		}
		return Error(ErrLen, a.Obj)
	}
}

// Nil validates the object is nil.
func (a AnyValidators) Nil() Validator {
	return func() error {
		if IsNil(a.Obj) {
			return nil
		}
		return Error(ErrNil, a.Obj)
	}
}

// NotNil validates the object is not nil.
// It also validates that the object is not an unset pointer.
func (a AnyValidators) NotNil() Validator {
	return func() error {
		if verr := a.Nil()(); verr != nil {
			return nil
		}
		return Error(ErrNotNil, a.Obj)
	}
}

// Equals validates an object equals another object.
func (a AnyValidators) Equals(expected interface{}) Validator {
	return func() error {
		actual := a.Obj

		if IsNil(expected) && IsNil(actual) {
			return nil
		}
		if (IsNil(expected) && !IsNil(actual)) || (!IsNil(expected) && IsNil(actual)) {
			return Error(ErrEquals, a.Obj)
		}

		actualType := reflect.TypeOf(actual)
		if actualType == nil {
			return Error(ErrEquals, a.Obj)
		}
		expectedValue := reflect.ValueOf(expected)
		if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
			if !reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual) {
				return Error(ErrEquals, a.Obj)
			}
		}

		if !reflect.DeepEqual(expected, actual) {
			return Error(ErrEquals, a.Obj)
		}
		return nil
	}
}

// NotEquals validates an object does not equal another object.
func (a AnyValidators) NotEquals(expected interface{}) Validator {
	return func() error {
		if verr := a.Equals(expected)(); verr != nil {
			return nil
		}
		return Error(ErrNotEquals, a.Obj)
	}
}

// Allow validates a field is one of a given set of allowed values.
func (a AnyValidators) Allow(values ...interface{}) Validator {
	return func() error {
		for _, expected := range values {
			if verr := a.Equals(expected)(); verr == nil {
				return nil
			}
		}
		return Error(ErrAllowed, a.Obj)
	}
}

// Disallow validates a field is one of a given set of allowed values.
func (a AnyValidators) Disallow(values ...interface{}) Validator {
	return func() error {
		for _, expected := range values {
			if verr := a.Equals(expected)(); verr == nil {
				return Error(ErrDisallowed, a.Obj)
			}
		}
		return nil
	}
}
