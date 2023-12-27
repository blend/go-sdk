/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"context"
	"database/sql"
	"reflect"
	"unsafe"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

// NOTE: Ensure that
//   - `AlwaysFailDB` satisfies `db.DB`.
//   - `PseudoQueryDB` satisfies `db.DB`.
var (
	_ db.DB = (*AlwaysFailDB)(nil)
	_ db.DB = (*PseudoQueryDB)(nil)
)

// AlwaysFailDB implements the `db.DB` interface, but each method always fails.
type AlwaysFailDB struct {
	Errors ErrorProducer
}

// ExecContext implements the `db.DB` interface and returns and error.
func (afd *AlwaysFailDB) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, afd.Errors.NextError()
}

// QueryContext implements the `db.DB` interface and returns and error.
func (afd *AlwaysFailDB) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, afd.Errors.NextError()
}

// QueryRowContext implements the `db.DB` interface; the error value is embedded
// in the `sql.Row` value returned.
func (afd *AlwaysFailDB) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return sqlRowWithError(afd.Errors.NextError())
}

// PseudoQueryDB implements the `db.DB` interface, it intercepts calls to
// `QueryContext` and replaces the `query` / `args` arguments with custom
// values.
type PseudoQueryDB struct {
	DB    *sql.DB
	Query string
	Args  []interface{}
}

// ExecContext implements the `db.DB` interface; this is not supported in
// `PseudoQueryDB`. It will **always** return a "not implemented" error.
func (pqd *PseudoQueryDB) ExecContext(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, ex.New("Not Implemented: ExecContext")
}

// QueryContext implements the `db.DB` interface. It intercepts the **actual**
// query and arguments and replaces them with the query and arguments stored
// on the current `PseudoQueryDB`.
func (pqd *PseudoQueryDB) QueryContext(ctx context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return pqd.DB.QueryContext(ctx, pqd.Query, pqd.Args...)
}

// QueryRowContext implements the `db.DB` interface; this is not supported in
// `PseudoQueryDB`. It will **always** return a `sql.Row` with a "not implemented"
// error set on the row.
func (pqd *PseudoQueryDB) QueryRowContext(_ context.Context, _ string, _ ...interface{}) *sql.Row {
	err := ex.New("Not Implemented: QueryRowContext")
	return sqlRowWithError(err)
}

func sqlRowWithError(err error) *sql.Row {
	row := &sql.Row{}
	e := reflect.ValueOf(row).Elem().FieldByName("err")
	eMutable := reflect.NewAt(e.Type(), unsafe.Pointer(e.UnsafeAddr())).Elem()
	eMutable.Set(reflect.ValueOf(err))
	return row
}
