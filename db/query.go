package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// --------------------------------------------------------------------------------
// Query Result
// --------------------------------------------------------------------------------

// Query is the intermediate result of a query.
type Query struct {
	context       context.Context
	statement     string
	cachedPlanKey string
	args          []interface{}

	rows *sql.Rows
	err  error

	conn *Connection
	inv  *Invocation
	tx   *sql.Tx
}

// Execute runs a given query, yielding the raw results.
func (q *Query) Execute() (rows *sql.Rows, err error) {
	defer func() { err = q.finish(recover(), err) }()
	rows, err = q.query()
	return
}

// Any returns if there are any results for the query.
func (q *Query) Any() (hasRows bool, err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, q.rows.Close()) }()

	hasRows = q.rows.Next()
	return
}

// None returns if there are no results for the query.
func (q *Query) None() (hasRows bool, err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.rows.Close())) }()
	hasRows = !q.rows.Next()
	return
}

// Scan writes the results to a given set of local variables.
func (q *Query) Scan(args ...interface{}) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.rows.Close())) }()

	if q.rows.Next() {
		if err = q.rows.Scan(args...); err != nil {
			err = Error(err)
			return
		}
	}

	return
}

// Out writes the query result to a single object via. reflection mapping.
func (q *Query) Out(object interface{}) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.rows.Close())) }()

	sliceType := ReflectType(object)
	if sliceType.Kind() != reflect.Struct {
		err = Error(ErrDestinationNotStruct)
		return
	}

	columnMeta := CachedColumnCollectionFromInstance(object)
	if q.rows.Next() {
		if populatable, ok := object.(Populatable); ok {
			err = populatable.Populate(q.rows)
		} else {
			err = PopulateByName(object, q.rows, columnMeta)
		}
		if err != nil {
			return
		}
	}

	return
}

// OutMany writes the query results to a slice of objects.
func (q *Query) OutMany(collection interface{}) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, q.rows.Close()) }()

	sliceType := ReflectType(collection)
	if sliceType.Kind() != reflect.Slice {
		err = Error(ErrCollectionNotSlice)
		return
	}

	sliceInnerType := ReflectSliceType(collection)
	collectionValue := ReflectValue(collection)
	v := makeNew(sliceInnerType)
	meta := CachedColumnCollectionFromType(newColumnCacheKey(sliceInnerType), sliceInnerType)

	isPopulatable := isPopulatable(v)

	didSetRows := false
	for q.rows.Next() {
		newObj := makeNew(sliceInnerType)

		if isPopulatable {
			err = asPopulatable(newObj).Populate(q.rows)
		} else {
			err = PopulateByName(newObj, q.rows, meta)
		}
		if err != nil {
			return
		}

		newObjValue := ReflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
		didSetRows = true
	}

	if !didSetRows {
		collectionValue.Set(reflect.MakeSlice(sliceType, 0, 0))
	}
	return
}

// Each executes the consumer for each result of the query (one to many).
func (q *Query) Each(consumer RowsConsumer) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.rows.Close())) }()

	for q.rows.Next() {
		if err = consumer(q.rows); err != nil {
			err = Error(err)
			return
		}
	}
	return
}

// First executes the consumer for the first result of a query.
func (q *Query) First(consumer RowsConsumer) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.rows, q.err = q.query()
	if q.err != nil {
		err = q.err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.rows.Close())) }()

	if q.rows.Next() {
		if err = consumer(q.rows); err != nil {
			return
		}
	}
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (q *Query) query() (rows *sql.Rows, err error) {
	if q.err != nil {
		err = q.err
		return
	}

	stmt, stmtErr := q.inv.Prepare(q.statement)
	if stmtErr != nil {
		err = Error(stmtErr)
		return
	}
	defer func() { err = q.inv.CloseStatement(stmt, err) }()

	rows, err = stmt.QueryContext(q.context, q.args...)
	if err != nil && !ex.Is(err, sql.ErrNoRows) {
		err = Error(err)
	}
	return
}

func (q *Query) finish(r interface{}, err error) error {
	return q.inv.Finish(q.statement, r, err)
}
