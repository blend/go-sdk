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
		col.ColumnName = strings.ToLower(field.Name)
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

// SetZero sets the empty value on a field on a database mapped object
func (c Column) SetZero(object interface{}) {
	objValue := ReflectValue(object)
	field := objValue.FieldByName(c.FieldName)
	field.Set(reflect.Zero(field.Type()))
}

// SetValue sets the field on a database mapped object to the instance of `value`.
func (c Column) SetValue(object interface{}, value interface{}) error {
	objectValue := ReflectValue(object)

	objectField := objectValue.FieldByName(c.FieldName)
	objectFieldType := objectField.Type()

	// check if we've been passed a reference for the target object
	if !objectField.CanSet() {
		return ex.New("hit a field we can't set; did you forget to pass the object as a reference?").WithMessagef("field: %s", c.FieldName)
	}

	valueReflected := ReflectValue(value)
	if !valueReflected.IsValid() { // if the value is nil
		objectField.Set(reflect.Zero(objectFieldType)) // zero the field
		return nil
	}

	// special case for `db:"...,json"` fields.
	if c.IsJSON {
		// check if we have a driver nullable string for the given value
		valueContents, ok := valueReflected.Interface().(sql.NullString)
		// if it's a nullable string, and it's valid (set) and not empty
		if ok && valueContents.Valid && len(valueContents.String) > 0 {
			// grab a pointer to the field on the object
			fieldAddr := objectField.Addr().Interface()
			jsonErr := json.Unmarshal([]byte(valueContents.String), fieldAddr)
			if jsonErr != nil {
				return ex.New(jsonErr)
			}
			// set the field to the indirect of the value we just set
			objectField.Set(reflect.ValueOf(fieldAddr).Elem())
		}
		return nil
	}

	// if we can direct assign
	// this will handle cases where you're setting a * to a *
	// or the intrinsic type to itself.
	if valueReflected.Type().AssignableTo(objectFieldType) {
		if objectField.Kind() == reflect.Ptr && valueReflected.CanAddr() {
			objectField.Set(valueReflected.Addr())
		} else {
			objectField.Set(valueReflected)
		}
		return nil
	}

	if objectField.Kind() == reflect.Ptr {
		if valueReflected.CanAddr() {
			if valueReflected.Type().AssignableTo(objectFieldType.Elem()) {
				objectField.Set(valueReflected.Addr())
				return nil
			}

			convertedValue := valueReflected.Convert(objectFieldType.Elem())
			if convertedValue.CanAddr() {
				objectField.Set(convertedValue.Addr())
				return nil
			}
			// what should we do here?
			return ex.New("cannot convert value to destination pointer type")
		}
		return ex.New("cannot take address of value for assignment to object field (which is itself a pointer)")
	}

	// convert and assign
	convertedValue := valueReflected.Convert(objectFieldType)
	objectField.Set(convertedValue)
	return nil
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
