package validate

import (
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// Basic errors
const (
	ErrNil       ex.Class = "object should be nil"
	ErrNotNil    ex.Class = "object should not be nil"
	ErrEquals    ex.Class = "objects should be equal"
	ErrNotEquals ex.Class = "objects should not be equal"
)

// BasicValidator validates any object.
type BasicValidator func(interface{}) error

// Nil validates the object is nil.
func Nil(obj interface{}) error {
	if obj == nil {
		return Error(ErrNil)
	}

	value := reflect.ValueOf(obj)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return Error(ErrNil)
	}
	return nil
}

// NotNil validates the object is not nil.
// It also validates that the object is not an unset pointer.
func NotNil(obj interface{}) error {
	if verr := Nil(obj); verr != nil {
		return nil
	}
	return Error(ErrNotNil)
}

// Equals validates an object equals another object.
func Equals(expected interface{}) BasicValidator {
	return func(actual interface{}) error {
		if expected == nil && actual == nil {
			return nil
		}
		if (expected == nil && actual != nil) || (expected != nil && actual == nil) {
			return Error(ErrEquals)
		}

		actualType := reflect.TypeOf(actual)
		if actualType == nil {
			return Error(ErrEquals)
		}
		expectedValue := reflect.ValueOf(expected)
		if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
			if !reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual) {
				return Error(ErrEquals)
			}
		}

		if !reflect.DeepEqual(expected, actual) {
			return Error(ErrEquals)
		}
		return nil
	}
}

// NotEquals validates an object does not equal another object.
func NotEquals(expected interface{}) BasicValidator {
	return func(actual interface{}) error {
		if verr := Equals(expected)(actual); verr != nil {
			return nil
		}
		return Error(ErrNotEquals)
	}
}
