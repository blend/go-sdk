package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/blend/go-sdk/bufferutil"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

const (
	//DBNilError is a common error
	DBNilError = "connection is nil"
)

const (
	runeComma   = rune(',')
	runeNewline = rune('\n')
	runeTab     = rune('\t')
	runeSpace   = rune(' ')
)

// --------------------------------------------------------------------------------
// Connection
// --------------------------------------------------------------------------------

// New returns a new Connection.
// It will use very bare bones defaults for the config.
func New(options ...Option) (*Connection, error) {
	c := &Connection{
		PlanCache: NewPlanCache(),
	}
	var err error
	for _, o := range options {
		if err = o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
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
//    conn, err := db.Open(db.NewFromConfig(cfg))
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
	sync.Mutex
	Config               Config
	Tracer               Tracer
	StatementInterceptor StatementInterceptor
	Connection           *sql.DB
	BufferPool           *bufferutil.Pool
	Log                  logger.Log
	PlanCache            *PlanCache
}

// Close implements a closer.
func (dbc *Connection) Close() error {
	if dbc.PlanCache != nil {
		if err := dbc.PlanCache.Close(); err != nil {
			return err
		}
	}
	return dbc.Connection.Close()
}

// Open returns a connection object, either a cached connection object or creating a new one in the process.
func (dbc *Connection) Open() error {
	dbc.Lock()
	defer dbc.Unlock()

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
	if dbc.PlanCache == nil {
		dbc.PlanCache = NewPlanCache()
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

	dbc.PlanCache.WithConnection(dbConn)
	dbc.PlanCache.WithEnabled(!dbc.Config.PlanCacheDisabled)
	dbc.Connection = dbConn
	dbc.Connection.SetConnMaxLifetime(dbc.Config.MaxLifetimeOrDefault())
	dbc.Connection.SetMaxIdleConns(dbc.Config.IdleConnectionsOrDefault())
	dbc.Connection.SetMaxOpenConns(dbc.Config.MaxConnectionsOrDefault())
	return nil
}

// Begin starts a new transaction.
func (dbc *Connection) Begin(opts ...*sql.TxOptions) (*sql.Tx, error) {
	return dbc.BeginContext(context.Background(), opts...)
}

// BeginContext starts a new transaction in a givent context.
func (dbc *Connection) BeginContext(context context.Context, opts ...*sql.TxOptions) (*sql.Tx, error) {
	if dbc.Connection == nil {
		return nil, ex.New(ErrConnectionClosed)
	}
	if len(opts) > 0 {
		tx, err := dbc.Connection.BeginTx(context, opts[0])
		return tx, Error(err)
	}
	tx, err := dbc.Connection.BeginTx(context, nil)
	return tx, Error(err)
}

// PrepareContext prepares a statement potentially returning a cached version of the statement.
func (dbc *Connection) PrepareContext(context context.Context, cachedPlanKey, statement string, tx *sql.Tx) (stmt *sql.Stmt, err error) {
	if dbc.Tracer != nil {
		tf := dbc.Tracer.Prepare(context, dbc, statement)
		if tf != nil {
			defer func() { tf.Finish(err) }()
		}
	}
	if tx != nil {
		stmt, err = tx.PrepareContext(context, statement)
		return
	}
	if dbc.PlanCache != nil && dbc.PlanCache.Enabled() && cachedPlanKey != "" {
		stmt, err = dbc.PlanCache.PrepareContext(context, cachedPlanKey, statement)
		return
	}
	stmt, err = dbc.Connection.PrepareContext(context, statement)
	return
}

// --------------------------------------------------------------------------------
// Invocation
// --------------------------------------------------------------------------------

// Invoke returns a new invocation.
func (dbc *Connection) Invoke(options ...InvocationOption) *Invocation {
	i := &Invocation{
		Context:              context.Background(),
		Tracer:               dbc.Tracer,
		StatementInterceptor: dbc.StatementInterceptor,
		Conn:                 dbc,
		StartTime:            time.Now().UTC(),
	}
	for _, option := range options {
		option(i)
	}
	return i
}

// Ping checks the db connection.
func (dbc *Connection) Ping() error {
	return Error(dbc.Connection.Ping())
}

// PingContext checks the db connection.
func (dbc *Connection) PingContext(context context.Context) (err error) {
	if dbc.Tracer != nil {
		tf := dbc.Tracer.Ping(context, dbc)
		if tf != nil {
			defer func() {
				tf.Finish(err)
			}()
		}
	}

	err = Error(dbc.Connection.PingContext(context))
	return
}

// Exec is a helper stub for .Invoke(...).Exec(...).
func (dbc *Connection) Exec(statement string, args ...interface{}) error {
	return dbc.Invoke().Exec(statement, args...)
}

// ExecContext is a helper stub for .Invoke(OptContext(ctx)).Exec(...).
func (dbc *Connection) ExecContext(ctx context.Context, statement string, args ...interface{}) error {
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
