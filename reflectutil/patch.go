/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reflectutil

import (
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// Patcher describes an object that can be patched with raw values.
type Patcher interface {
	Patch(map[string]interface{}) error
}

// Patch updates an object based on a map of field names to values.
func Patch(obj interface{}, patchValues map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(r)
		}
	}()

	if patchable, isPatchable := obj.(Patcher); isPatchable {
		return patchable.Patch(patchValues)
	}

	targetValue := Value(obj)
	targetType := targetValue.Type()

	for key, value := range patchValues {
		err = SetValue(obj, targetType, targetValue, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetValue sets a value on an object by its field name.
func SetValue(obj interface{}, objType reflect.Type, objValue reflect.Value, fieldName string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ex.New("panic setting value by name", ex.OptMessagef("field: %s panic: %v", fieldName, r))
		}
	}()

	relevantField, hasField := objType.FieldByName(fieldName)
	if !hasField {
		err = ex.New("unknown field", ex.OptMessagef("%s `%s`", objType.Name(), fieldName))
		return
	}

	return doSetValue(relevantField, objType, objValue, fieldName, value)
}

func doSetValue(relevantField reflect.StructField, objType reflect.Type, objValue reflect.Value, name string, value interface{}) (err error) {
	field := objValue.FieldByName(relevantField.Name)
	if !field.CanSet() {
		err = ex.New("cannot set field", ex.OptMessagef("%s `%s`", objType.Name(), name))
		return
	}

	valueReflected := Value(value)
	if !valueReflected.IsValid() {
		err = ex.New("invalid value", ex.OptMessagef("%s `%s`", objType.Name(), name))
		return
	}

	assigned, assignErr := tryAssignment(field, valueReflected)
	if assignErr != nil {
		err = assignErr
		return
	}
	if !assigned {
		err = ex.New("cannot set field", ex.OptMessagef("%s `%s`", objType.Name(), name))
		return
	}
	return
}
