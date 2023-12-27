/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// --------------------------------------------------------------------------------
// Column
// --------------------------------------------------------------------------------

// NewColumnFromFieldTag reads the contents of a field tag, ex: `json:"foo" db:"bar,isprimarykey,isserial"
func NewColumnFromFieldTag(field reflect.StructField) *Column {
	db := field.Tag.Get("db")
	if db != "-" {
		col := Column{}
		col.FieldName = field.Name
		col.ColumnName = field.Name
		col.FieldType = field.Type
		if db != "" {
			pieces := strings.Split(db, ",")

			if !strings.HasPrefix(db, ",") {
				col.ColumnName = pieces[0]
			}

			if len(pieces) >= 1 {
				args := strings.ToLower(strings.Join(pieces[1:], ","))

				col.IsPrimaryKey = strings.Contains(args, "pk")
				col.IsUniqueKey = strings.Contains(args, "uk")
				col.IsAuto = strings.Contains(args, "serial") || strings.Contains(args, "auto")
				col.IsReadOnly = strings.Contains(args, "readonly")
				col.Inline = strings.Contains(args, "inline")
				col.IsJSON = strings.Contains(args, "json")
			}
		}
		return &col
	}

	return nil
}

// Column represents a single field on a struct that is mapped to the database.
type Column struct {
	Parent       *Column
	TableName    string
	FieldName    string
	FieldType    reflect.Type
	ColumnName   string
	Index        int
	IsPrimaryKey bool
	IsUniqueKey  bool
	IsAuto       bool
	IsReadOnly   bool
	IsJSON       bool
	Inline       bool
}

// SetValue sets the field on a database mapped object to the instance of `value`.
func (c Column) SetValue(object, value interface{}) error {
	return c.SetValueReflected(ReflectValue(object), value)
}

// SetValueReflected sets the field on a reflect value object to the instance of `value`.
func (c Column) SetValueReflected(objectValue reflect.Value, value interface{}) error {
	objectField := objectValue.FieldByName(c.FieldName)

	// check if we've been passed a reference for the target object
	if !objectField.CanSet() {
		return ex.New("hit a field we can't set; did you forget to pass the object as a reference?").WithMessagef("field: %s", c.FieldName)
	}

	// special case for `db:"...,json"` fields.
	if c.IsJSON {
		var deserialized interface{}
		if objectField.Kind() == reflect.Ptr {
			deserialized = reflect.New(objectField.Type().Elem()).Interface()
		} else {
			deserialized = objectField.Addr().Interface()
		}

		switch valueContents := value.(type) {
		case *sql.NullString:
			if !valueContents.Valid {
				objectField.Set(reflect.Zero(objectField.Type()))
				return nil
			}
			if err := json.Unmarshal([]byte(valueContents.String), deserialized); err != nil {
				return ex.New(err).WithMessage(valueContents.String)
			}
		case sql.NullString:
			if !valueContents.Valid {
				objectField.Set(reflect.Zero(objectField.Type()))
				return nil
			}
			if err := json.Unmarshal([]byte(valueContents.String), deserialized); err != nil {
				return ex.New(err)
			}
		case *string:
			if err := json.Unmarshal([]byte(*valueContents), deserialized); err != nil {
				return ex.New(err)
			}
		case string:
			if err := json.Unmarshal([]byte(valueContents), deserialized); err != nil {
				return ex.New(err)
			}
		case *[]byte:
			if err := json.Unmarshal(*valueContents, deserialized); err != nil {
				return ex.New(err)
			}
		case []byte:
			if err := json.Unmarshal(valueContents, deserialized); err != nil {
				return ex.New(err)
			}
		default:
			return ex.New("set value; invalid type for assignment to json field").WithMessagef("field: %s, value: %t", value)
		}

		if rv := reflect.ValueOf(deserialized); !rv.IsValid() {
			objectField.Set(reflect.Zero(objectField.Type()))
		} else {
			if objectField.Kind() == reflect.Ptr {
				objectField.Set(rv)
			} else {
				objectField.Set(rv.Elem())
			}
		}
		return nil
	}

	valueReflected := ReflectValue(value)
	if !valueReflected.IsValid() { // if the value is nil
		objectField.Set(reflect.Zero(objectField.Type())) // zero the field
		return nil
	}

	// if we can direct assign the value to the field
	if valueReflected.Type().AssignableTo(objectField.Type()) {
		objectField.Set(valueReflected)
		return nil
	}

	// convert and assign
	if valueReflected.Type().ConvertibleTo(objectField.Type()) ||
		haveSameUnderlyingTypes(objectField, valueReflected) {
		objectField.Set(valueReflected.Convert(objectField.Type()))
		return nil
	}

	if objectField.Kind() == reflect.Ptr && valueReflected.CanAddr() {
		if valueReflected.Addr().Type().AssignableTo(objectField.Type()) {
			objectField.Set(valueReflected.Addr())
			return nil
		}
		if valueReflected.Addr().Type().ConvertibleTo(objectField.Type()) {
			objectField.Set(valueReflected.Convert(objectField.Elem().Type()).Addr())
			return nil
		}
		return ex.New("set value; can addr value but can't figure out how to assign or convert").WithMessagef("field: %s, value: %#v", c.FieldName, value)
	}

	return ex.New("set value; ran out of ways to set the field").WithMessagef("field: %s, value: %#v", c.FieldName, value)
}

// GetValue returns the value for a column on a given database mapped object.
func (c Column) GetValue(object DatabaseMapped) interface{} {
	value := ReflectValue(object)
	if c.Parent != nil {
		embedded := value.Field(c.Parent.Index)
		valueField := embedded.Field(c.Index)
		return valueField.Interface()
	}
	valueField := value.Field(c.Index)
	return valueField.Interface()
}
