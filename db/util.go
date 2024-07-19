/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
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
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		instance = reflect.New(t).Interface()
		if typed, isTyped := instance.(TableNameProvider); isTyped {
			return typed.TableName()
		}
	}
	return t.Name()
}

// TableName returns the mapped table name for a given instance; it will sniff for the `TableName()` function on the type.
func TableName(obj DatabaseMapped) string {
	if typed, isTyped := obj.(TableNameProvider); isTyped {
		return typed.TableName()
	}
	return ReflectType(obj).Name()
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
// Result utility methods
// --------------------------------------------------------------------------------

// IgnoreExecResult is a helper for use with .Exec() (sql.Result, error)
// that ignores the result return.
func IgnoreExecResult(_ sql.Result, err error) error {
	return err
}

// ExecRowsAffected is a helper for use with .Exec() (sql.Result, error)
// that returns the rows affected.
func ExecRowsAffected(i sql.Result, inputErr error) (int64, error) {
	if inputErr != nil {
		return 0, inputErr
	}
	ra, err := i.RowsAffected()
	if err != nil {
		return 0, Error(err)
	}
	return ra, nil
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
