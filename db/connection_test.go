/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

// Test_Connection_useBeforeOpen tests if we can connect to the db, a.k.a., if the underlying driver works.
func Test_Connection_useBeforeOpen(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	conn, err := New()
	its.Nil(err)

	tx, err := conn.Begin()
	its.NotNil(err)
	its.True(ex.Is(ErrConnectionClosed, err))
	its.Nil(tx)

	inv := conn.Invoke()
	its.Nil(inv.DB)
	its.True(inv.DB == nil)

	any, err := conn.Query("select 1").Any()
	its.NotNil(err)
	its.True(ex.Is(ErrConnectionClosed, err), err.Error())
	its.False(any)
}

// TestConnectionSanityCheck tests if we can connect to the db, a.k.a., if the underlying driver works.
func TestConnectionSanityCheck(t *testing.T) {
	assert := assert.New(t)

	conn, err := OpenTestConnection()
	assert.Nil(err)
	str := conn.Config.CreateDSN()
	_, err = sql.Open("pgx", str)
	assert.Nil(err)
}

func TestPrepareContext(t *testing.T) {
	a := assert.New(t)

	conn, err := OpenTestConnection()
	a.Nil(err)

	var calledPrepare, calledFinish bool
	conn.Tracer = mockTracer{
		PrepareHandler: func(_ context.Context, _ Config, _ string) {
			calledPrepare = true
		},
		FinishPrepareHandler: func(_ context.Context, _ error) {
			calledFinish = true
		},
	}

	stmt, err := conn.PrepareContext(context.TODO(), "select 'ok!'", nil)
	a.Nil(err)
	defer stmt.Close()
	a.NotNil(stmt)
	a.True(calledPrepare)
	a.True(calledFinish)
}

func TestQuery(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer func() { _ = tx.Rollback() }()

	a.Equal(DefaultSchema, defaultDB().Config.SchemaOrDefault())
	err = seedObjects(100, tx)
	a.Nil(err)

	objs := []benchObj{}
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").OutMany(&objs)
	a.Nil(err)
	a.NotEmpty(objs)

	err = defaultDB().Invoke(OptTx(tx), OptInvocationStatementInterceptor(failInterceptor)).Query("select * from bench_object").OutMany(&objs)
	a.Equal("this is just an interceptor error", err.Error())
}

func Test_Connection_Open(t *testing.T) {
	a := assert.New(t)

	conn, err := New(OptConfigFromEnv())
	a.Nil(err)
	a.Nil(conn.Open())
	defer conn.Close()

	a.NotNil(conn.BufferPool)
	a.NotNil(conn.Connection)
}

func TestExec(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("select 'ok!'"))
	a.Nil(err)

	err = IgnoreExecResult(defaultDB().Invoke(OptTx(tx), OptInvocationStatementInterceptor(failInterceptor)).Exec("select 'ok!'"))
	a.Equal("this is just an interceptor error", err.Error())
}

// TestConnectionConfigSetsDatabase tests if we set the .database property on open.
func TestConnectionConfigSetsDatabase(t *testing.T) {
	assert := assert.New(t)
	conn, err := New(OptConfigFromEnv())
	assert.Nil(err)
	assert.Nil(conn.Open())
	defer conn.Close()
	assert.NotEmpty(conn.Config.DatabaseOrDefault())
}
