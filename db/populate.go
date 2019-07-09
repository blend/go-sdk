package db

import (
	"database/sql"
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// PopulateEmpty populates all the column fields of a struct with empty value
func PopulateEmpty(object interface{}, cols *ColumnCollection) {
	for _, v := range cols.columns {
		v.SetZero(object)
	}
}

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
	for i, v := range values {
		colName = rowColumns[i]
		if field, ok = columnLookup[colName]; ok {
			err = field.SetValue(object, v)
			if err != nil {
				return ex.New(Error(err), ex.OptMessagef("column: %s", colName))
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

	columns := cols.Columns()
	var field Column
	for i, v := range values {
		field = columns[i]
		if err = field.SetValue(object, v); err != nil {
			err = ex.New(err)
			return
		}
	}

	return
}

// initColumnValue inserts the correct placeholder in the scan array of values.
// it will use `sql.Null` forms where appropriate.
// JSON fields are implicitly nullable.
func initColumnValue(index int, values []interface{}, col *Column) {
	if col.IsJSON {
		values[index] = &sql.NullString{}
	} else {
		values[index] = reflect.New(reflect.PtrTo(col.FieldType)).Interface()
	}
}
