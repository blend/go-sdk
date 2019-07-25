package validate

import (
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrNonLengthType ex.Class = "instance is a non-length type"
)

// GetLength returns the length of an object or an error if it's not a thing that can have a length.
func GetLength(obj interface{}) (int, error) {
	objValue := reflect.ValueOf(obj)
	switch objValue.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.String:
		{
			if obj == nil {
				return 0, nil
			} else if obj == "" {
				return 0, nil
			}
			if objValue.IsValid() {
				return objValue.Len(), nil
			}
			return 0, nil
		}
	}
	return 0, ErrNonLengthType
}
