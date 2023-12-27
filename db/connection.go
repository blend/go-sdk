/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/bufferutil"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

var (
	_ async.Checker = (*Connection)(nil)
)

// --------------------------------------------------------------------------------
// Connection
// --------------------------------------------------------------------------------

// New returns a new Connection.
// It will use very bare bones defaults for the config.
func New(options ...Option) (*Connection, error) {
	var c Connection
	var err error
	for _, opt := range options {
		if err = opt(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// MustNew returns a new connection and panics on error.
func MustNew(options ...Option) *Connection {
	c, err := New(options...)
	if err != nil {
		panic(err)
	}
	return c
}

// Open opens a connection, testing an error and returning it if not nil, and if nil, opening the connection.
// It's designed ot be used in conjunction with a constructor, i.e.
//
//	conn, err := db.Open(db.NewFromConfig(cfg))
func Open(conn *Connection, err error) (*Connection, error) {
	if err != nil {
		return nil, err
	}
	if err = conn.Open(); err != nil {
		return nil, err
	}
	return conn, nil
}

// Connection is the basic wrapper for connection parameters and saves a reference to the created sql.Connection.
type Connection struct {
	Connection		*sql.DB
	BufferPool		*bufferutil.Pool
	Config			Config
	Log			logger.Log
	Tracer			Tracer
	StatementInterceptor	StatementInterceptor
}

// Close implements a closer.
func (dbc *Connection) Close() error {
	return dbc.Connection.Close()
}

// Open returns a connection object, either a cached connection object or creating a new one in the process.
func (dbc *Connection) Open() error {
	// bail if we've already opened the connection.
	if dbc.Connection != nil {
		return Error(ErrConnectionAlreadyOpen)
	}
	if dbc.Config.IsZero() {
		return Error(ErrConfigUnset)
	}
	if dbc.BufferPool == nil {
		dbc.BufferPool = bufferutil.NewPool(dbc.Config.BufferPoolSizeOrDefault())
	}

	dsn := dbc.Config.CreateDSN()
	namedValues, err := ParseURL(dsn)
	if err != nil {
		return err
	}

	// open the connection
	dbConn, err := sql.Open(dbc.Config.EngineOrDefault(), namedValues)
	if err != nil {
		return Error(err)
	}

	dbc.Connection = dbConn
	dbc.Connection.SetConnMaxLifetime(dbc.Config.MaxLifetimeOrDefault())
	dbc.Connection.SetConnMaxIdleTime(dbc.Config.MaxIdleTimeOrDefault())
	dbc.Connection.SetMaxIdleConns(dbc.Config.IdleConnectionsOrDefault())
	dbc.Connection.SetMaxOpenConns(dbc.Config.MaxConnectionsOrDefault())
	return nil
}

// Begin starts a new transaction.
func (dbc *Connection) Begin(opts ...func(*sql.TxOptions)) (*sql.Tx, error) {
	if dbc.Connection == nil {
		return nil, ex.New(ErrConnectionClosed)
	}
	return dbc.BeginContext(context.Background(), opts...)
}

// BeginContext starts a new transaction in a givent context.
func (dbc *Connection) BeginContext(ctx context.Context, opts ...func(*sql.TxOptions)) (*sql.Tx, error) {
	if dbc.Connection == nil {
		return nil, ex.New(ErrConnectionClosed)
	}
	var txOptions sql.TxOptions
	for _, opt := range opts {
		opt(&txOptions)
	}
	tx, err := dbc.Connection.BeginTx(ctx, &txOptions)
	return tx, Error(err)
}

// PrepareContext prepares a statement within a given context.
// If a tx is provided, the tx is the target for the prepare.
// This will trigger tracing on prepare.
func (dbc *Connection) PrepareContext(ctx context.Context, statement string, tx *sql.Tx) (stmt *sql.Stmt, err error) {
	if dbc.Tracer != nil {
		tf := dbc.Tracer.Prepare(ctx, dbc.Config, statement)
		if tf != nil {
			defer func() { tf.FinishPrepare(ctx, err) }()
		}
	}
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, statement)
		return
	}
	stmt, err = dbc.Connection.PrepareContext(ctx, statement)
	return
}

// --------------------------------------------------------------------------------
// Invocation
// --------------------------------------------------------------------------------

// Invoke returns a new invocation.
func (dbc *Connection) Invoke(options ...InvocationOption) *Invocation {
	i := Invocation{
		Config:			dbc.Config,
		BufferPool:		dbc.BufferPool,
		Context:		context.Background(),
		Log:			dbc.Log,
		Tracer:			dbc.Tracer,
		StatementInterceptor:	dbc.StatementInterceptor,
	}
	if dbc.Connection != nil {
		i.DB = dbc.Connection
	}
	for _, option := range options {
		option(&i)
	}
	return &i
}

// Exec is a helper stub for .Invoke(...).Exec(...).
func (dbc *Connection) Exec(statement string, args ...interface{}) (sql.Result, error) {
	return dbc.Invoke().Exec(statement, args...)
}

// ExecContext is a helper stub for .Invoke(OptContext(ctx)).Exec(...).
func (dbc *Connection) ExecContext(ctx context.Context, statement string, args ...interface{}) (sql.Result, error) {
	return dbc.Invoke(OptContext(ctx)).Exec(statement, args...)
}

// Query is a helper stub for .Invoke(...).Query(...).
func (dbc *Connection) Query(statement string, args ...interface{}) *Query {
	return dbc.Invoke().Query(statement, args...)
}

// QueryContext is a helper stub for .Invoke(OptContext(ctx)).Query(...).
func (dbc *Connection) QueryContext(ctx context.Context, statement string, args ...interface{}) *Query {
	return dbc.Invoke(OptContext(ctx)).Query(statement, args...)
}

// Check implements a status check.
func (dbc *Connection) Check(ctx context.Context) error {
	_, err := dbc.QueryContext(ctx, "select 1").Any()
	return err
}
