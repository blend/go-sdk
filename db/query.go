/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"database/sql"
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// --------------------------------------------------------------------------------
// Query Result
// --------------------------------------------------------------------------------

// Query is the intermediate result of a query.
type Query struct {
	Invocation *Invocation
	Statement  string
	Err        error
	Args       []interface{}
}

// Do runs a given query, yielding the raw results.
func (q *Query) Do() (rows *sql.Rows, err error) {
	defer func() {
		err = q.finish(recover(), err)
	}()
	rows, err = q.query()
	return
}

// Any returns if there are any results for the query.
func (q *Query) Any() (found bool, err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()
	rows, err = q.query()
	if err != nil {
		return
	}
	found = rows.Next()
	return
}

// None returns if there are no results for the query.
func (q *Query) None() (notFound bool, err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()
	rows, err = q.query()
	if err != nil {
		return
	}
	notFound = !rows.Next()
	return
}

// Scan writes the results to a given set of local variables.
// It returns if the query produced a row, and returns `ErrTooManyRows` if there
// are multiple row results.
func (q *Query) Scan(args ...interface{}) (found bool, err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()

	rows, err = q.query()
	if err != nil {
		return
	}
	found, err = Scan(rows, args...)
	return
}

// Out writes the query result to a single object via. reflection mapping. If there is more than one result, the first
// result is mapped to to object, and ErrTooManyRows is returned. Out() will apply column values for any colums
// in the row result to the object, potentially zeroing existing values out.
func (q *Query) Out(object interface{}) (found bool, err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()

	rows, err = q.query()
	if err != nil {
		return
	}
	found, err = Out(rows, object)
	return
}

// OutMany writes the query results to a slice of objects.
func (q *Query) OutMany(collection interface{}) (err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()

	rows, err = q.query()
	if err != nil {
		return
	}
	err = OutMany(rows, collection)
	return
}

// Each executes the consumer for each result of the query (one to many).
func (q *Query) Each(consumer RowsConsumer) (err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()

	rows, err = q.query()
	if err != nil {
		return
	}

	err = Each(rows, consumer)
	return
}

// First executes the consumer for the first result of a query.
// It returns `ErrTooManyRows` if more than one result is returned.
func (q *Query) First(consumer RowsConsumer) (found bool, err error) {
	var rows *sql.Rows
	defer func() {
		err = q.finish(recover(), err)
		err = q.rowsClose(rows, err)
	}()
	rows, err = q.query()
	if err != nil {
		return
	}
	found, err = First(rows, consumer)
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (q *Query) rowsClose(rows *sql.Rows, err error) error {
	if rows == nil {
		return err
	}
	return ex.Nest(err, rows.Close())
}

func (q *Query) query() (rows *sql.Rows, err error) {
	// fast abort if there was an issue ahead of returning the query.
	if q.Err != nil {
		err = q.Err
		return
	}

	var queryError error
	db := q.Invocation.DB
	ctx := q.Invocation.Context
	rows, queryError = db.QueryContext(ctx, q.Statement, q.Args...)
	if queryError != nil && !ex.Is(queryError, sql.ErrNoRows) {
		err = Error(queryError)
	}
	return
}

func (q *Query) finish(r interface{}, err error) error {
	return q.Invocation.finish(q.Statement, r, nil, err)
}

// Out reads a given rows set out into an object reference.
func Out(rows *sql.Rows, object interface{}) (found bool, err error) {
	sliceType := ReflectType(object)
	if sliceType.Kind() != reflect.Struct {
		err = Error(ErrDestinationNotStruct)
		return
	}
	columnMeta := Columns(object)
	if rows.Next() {
		found = true
		if populatable, ok := object.(Populatable); ok {
			err = populatable.Populate(rows)
		} else {
			err = PopulateByName(object, rows, columnMeta)
		}
		if err != nil {
			return
		}
	} else if err = Zero(object); err != nil {
		return
	}
	if rows.Next() {
		err = Error(ErrTooManyRows)
	}
	return
}

// OutMany reads a given result set into a given collection.
func OutMany(rows *sql.Rows, collection interface{}) (err error) {
	sliceType := ReflectType(collection)
	if sliceType.Kind() != reflect.Slice {
		err = Error(ErrCollectionNotSlice)
		return
	}

	sliceInnerType := ReflectSliceType(collection)
	collectionValue := ReflectValue(collection)
	v := makeNew(sliceInnerType)
	meta := ColumnsFromType(newColumnCacheKey(sliceInnerType), sliceInnerType)

	isPopulatable := IsPopulatable(v)

	var didSetRows bool
	for rows.Next() {
		newObj := makeNew(sliceInnerType)
		if isPopulatable {
			err = AsPopulatable(newObj).Populate(rows)
		} else {
			err = PopulateByName(newObj, rows, meta)
		}
		if err != nil {
			return
		}

		newObjValue := ReflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
		didSetRows = true
	}

	// this initializes the slice if we didn't add elements to it.
	if !didSetRows {
		collectionValue.Set(reflect.MakeSlice(sliceType, 0, 0))
	}
	return
}

// Each iterates over a given result set, calling the rows consumer.
func Each(rows *sql.Rows, consumer RowsConsumer) (err error) {
	for rows.Next() {
		if err = consumer(rows); err != nil {
			err = Error(err)
			return
		}
	}
	return
}

// First returns the first result of a result set to a consumer.
// If there are more than one row in the result, they are ignored.
func First(rows *sql.Rows, consumer RowsConsumer) (found bool, err error) {
	if found = rows.Next(); found {
		if err = consumer(rows); err != nil {
			return
		}
	}
	return
}

// Scan reads the first row from a resultset and scans it to a given set of args.
// If more than one row is returned it will return ErrTooManyRows.
func Scan(rows *sql.Rows, args ...interface{}) (found bool, err error) {
	if rows.Next() {
		found = true
		if err = rows.Scan(args...); err != nil {
			err = Error(err)
			return
		}
	}
	if rows.Next() {
		err = Error(ErrTooManyRows)
	}
	return
}
