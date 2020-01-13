package db

import (
	"context"
	"database/sql"
	"time"
	"github.com/blend/go-sdk/logger"
)

// InvocationOption is an option for invocations.
type InvocationOption func(*Invocation)

// OptLabel sets the Label on the invocation.
func OptLabel(label string) InvocationOption {
	return func(i *Invocation) {
		i.Label = label
	}
}

// OptInvocationStatementInterceptor sets the invocation statement interceptor.
func OptInvocationStatementInterceptor(interceptor StatementInterceptor) InvocationOption {
	return func(i *Invocation) { i.StatementInterceptor = interceptor }
}

// OptInvocationLog sets the invocation logger.
func OptInvocationLog(log logger.Log) InvocationOption {
	return func(i *Invocation) { i.Log = log }
}

// OptContext sets a context on an invocation.
func OptContext(ctx context.Context) InvocationOption {
	return func(i *Invocation) {
		i.Context = ctx
	}
}

// OptCancel sets the context cancel func..
func OptCancel(cancel context.CancelFunc) InvocationOption {
	return func(i *Invocation) {
		i.Cancel = cancel
	}
}

// OptTimeout sets a command timeout for the invocation.
func OptTimeout(d time.Duration) InvocationOption {
	return func(i *Invocation) {
		i.Context, i.Cancel = context.WithTimeout(i.Context, d)
	}
}

// OptTx is an invocation option that sets the invocation transaction.
func OptTx(tx *sql.Tx) InvocationOption {
	return func(i *Invocation) {
		if tx != nil {
			i.DB = tx
		}
	}
}

// OptDB is an invocation option that sets the underlying invocation db.
func OptDB(db DB) InvocationOption {
	return func(i *Invocation) {
		i.DB = db
	}
}
