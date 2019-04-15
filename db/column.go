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

// SetValue sets the field on a database mapped object to the instance of `value`.
func (c Column) SetValue(object interface{}, value interface{}) error {
	objValue := ReflectValue(object)
	field := objValue.FieldByName(c.FieldName)
	fieldType := field.Type()
	if !field.CanSet() {
		return ex.New("hit a field we can't set: '" + c.FieldName + "', did you forget to pass the object as a reference?")
	}

	valueReflected := ReflectValue(value)
	if !valueReflected.IsValid() {
		return nil
	}

	if c.IsJSON {
		valueContents, ok := valueReflected.Interface().(sql.NullString)
		if ok && valueContents.Valid && len(valueContents.String) > 0 {
			fieldAddr := field.Addr().Interface()
			jsonErr := json.Unmarshal([]byte(valueContents.String), fieldAddr)
			if jsonErr != nil {
				return ex.New(jsonErr)
			}
			field.Set(reflect.ValueOf(fieldAddr).Elem())
		}
		return nil
	}

	if valueReflected.Type().AssignableTo(fieldType) {
		if field.Kind() == reflect.Ptr && valueReflected.CanAddr() {
			field.Set(valueReflected.Addr())
		} else {
			field.Set(valueReflected)
		}
		return nil
	}

	if field.Kind() == reflect.Ptr {
		if valueReflected.CanAddr() {
			if fieldType.Elem() == valueReflected.Type() {
				field.Set(valueReflected.Addr())
			} else {
				convertedValue := valueReflected.Convert(fieldType.Elem())
				if convertedValue.CanAddr() {
					field.Set(convertedValue.Addr())
				}
			}
			return nil
		}

		return ex.New("cannot take address of value")
	}

	convertedValue := valueReflected.Convert(fieldType)
	field.Set(convertedValue)
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
