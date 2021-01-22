package dbutil

import (
	"context"

	"github.com/blend/go-sdk/db"
)

// NewBaseManager creates a new manager.
func NewBaseManager(conn *db.Connection, opts ...db.InvocationOption) BaseManager {
	return BaseManager{
		Conn:    conn,
		Options: opts,
	}
}

// BaseManager is the manager for database tasks.
//
// It is a base type you can use to build your own models
// that provides an `Invoke` method that will add default
// invocation options to a given invocation.
type BaseManager struct {
	Conn    *db.Connection
	Options []db.InvocationOption
}

// Invoke runs a command with a given set of options merged with the manager defaults.
func (m BaseManager) Invoke(ctx context.Context, opts ...db.InvocationOption) *db.Invocation {
	return m.Conn.Invoke(append(m.Options, append(opts, db.OptContext(ctx))...)...)
}
