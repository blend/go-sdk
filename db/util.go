package db

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

// TableNameByType returns the table name for a given reflect.Type by instantiating it and calling o.TableName().
// The type must implement DatabaseMapped or an exception will be returned.
func TableNameByType(t reflect.Type) string {
	instance := reflect.New(t).Interface()
	if typed, isTyped := instance.(TableNameProvider); isTyped {
		return typed.TableName()
	}
	return strings.ToLower(t.Name())
}

// TableName returns the mapped table name for a given instance; it will sniff for the `TableName()` function on the type.
func TableName(obj DatabaseMapped) string {
	if typed, isTyped := obj.(TableNameProvider); isTyped {
		return typed.TableName()
	}
	return strings.ToLower(ReflectType(obj).Name())
}

// --------------------------------------------------------------------------------
// String Utility Methods
// --------------------------------------------------------------------------------

// ParamTokens returns a csv token string in the form "$1,$2,$3...$N" if passed (1, N).
func ParamTokens(startAt, count int) string {
	if count < 1 {
		return ""
	}
	var str string
	for i := startAt; i < startAt+count; i++ {
		str = str + fmt.Sprintf("$%d", i)
		if i < (startAt + count - 1) {
			str = str + ","
		}
	}
	return str
}

// --------------------------------------------------------------------------------
// Internal / Reflection Utility Methods
// --------------------------------------------------------------------------------

// AsPopulatable casts an object as populatable.
func AsPopulatable(object interface{}) Populatable {
	return object.(Populatable)
}

// IsPopulatable returns if an object is populatable
func IsPopulatable(object interface{}) bool {
	_, isPopulatable := object.(Populatable)
	return isPopulatable
}

// ReflectValue returns the reflect.Value for an object following pointers.
func ReflectValue(obj interface{}) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// ReflectType retruns the reflect.Type for an object following pointers.
func ReflectType(obj interface{}) reflect.Type {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	return t
}

// ReflectSliceType returns the inner type of a slice following pointers.
func ReflectSliceType(collection interface{}) reflect.Type {
	v := reflect.ValueOf(collection)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Len() == 0 {
		t := v.Type()
		for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		return t
	}
	v = v.Index(0)
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Type()
}

// MakeWhereClause returns the sql `where` clause for a column collection, starting at a given index (used in sql $1 parameterization).
func MakeWhereClause(pks *ColumnCollection, startAt int) string {
	whereClause := " WHERE "
	for i, pk := range pks.Columns() {
		whereClause = whereClause + fmt.Sprintf("%s = %s", pk.ColumnName, "$"+strconv.Itoa(i+startAt))
		if i < (pks.Len() - 1) {
			whereClause = whereClause + " AND "
		}
	}

	return whereClause
}

// ParamTokensCSV returns a csv token string in the form "$1,$2,$3...$N"
func ParamTokensCSV(num int) string {
	str := ""
	for i := 1; i <= num; i++ {
		str = str + fmt.Sprintf("$%d", i)
		if i != num {
			str = str + ","
		}
	}
	return str
}

// makeNew creates a new object.
func makeNew(t reflect.Type) interface{} {
	return reflect.New(t).Interface()
}

func makeSliceOfType(t reflect.Type) interface{} {
	return reflect.New(reflect.SliceOf(t)).Interface()
}
