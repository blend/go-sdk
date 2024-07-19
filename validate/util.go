/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

// IsNil returns if an object is nil or is a typed pointer to nil.
func IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	value := reflect.ValueOf(obj)
	kind := value.Kind()
	return kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil()
}
