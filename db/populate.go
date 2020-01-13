package db

import (
	"database/sql"
	"reflect"
	"github.com/blend/go-sdk/ex"
)

// PopulateByName sets the values of an object from the values of a sql.Rows object using column names.
func PopulateByName(object interface{}, row Rows, cols *ColumnCollection) error {
	rowColumns, err := row.Columns()
	if err != nil {
		return Error(err)
	}

	var values = make([]interface{}, len(rowColumns))
	var columnLookup = cols.Lookup()
	for i, name := range rowColumns {
		if col, ok := columnLookup[name]; ok {
			initColumnValue(i, values, col)
		} else {
			var value interface{}
			values[i] = &value
		}
	}

	err = row.Scan(values...)
	if err != nil {
		return Error(err)
	}

	var colName string
	var field *Column
	var ok bool

	objectValue := ReflectValue(object)
	for i, v := range values {
		colName = rowColumns[i]
		if field, ok = columnLookup[colName]; ok {
			err = field.SetValueReflected(objectValue, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PopulateInOrder sets the values of an object in order from a sql.Rows object.
// Only use this method if you're certain of the column order. It is faster than populateByName.
// Optionally if your object implements Populatable this process will be skipped completely, which is even faster.
func PopulateInOrder(object DatabaseMapped, row Scanner, cols *ColumnCollection) (err error) {
	var values = make([]interface{}, cols.Len())

	for i, col := range cols.Columns() {
		initColumnValue(i, values, &col)
	}
	if err = row.Scan(values...); err != nil {
		return Error(err)
	}

	objectValue := ReflectValue(object)
	columns := cols.Columns()
	var field Column
	for i, v := range values {
		field = columns[i]
		if err = field.SetValueReflected(objectValue, v); err != nil {
			err = ex.New(err)
			return
		}
	}

	return
}

// Zero resets an object.
func Zero(object interface{}) error {
	objectValue := reflect.ValueOf(object)
	if !objectValue.Elem().CanSet() {
		return ex.New("zero; cannot set object, did you pass a reference?")
	}
	objectValue.Elem().Set(reflect.Zero(objectValue.Type().Elem()))
	return nil
}

// initColumnValue inserts the correct placeholder in the scan array of values.
// it will use `sql.Null` forms where appropriate.
// JSON fields are implicitly nullable.
func initColumnValue(index int, values []interface{}, col *Column) {
	if col.IsJSON {
		values[index] = &sql.NullString{}
	} else if col.FieldType.Kind() == reflect.Ptr {
		values[index] = reflect.New(col.FieldType).Interface()
	} else {
		values[index] = reflect.New(reflect.PtrTo(col.FieldType)).Interface()
	}
}
