package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
)

// Actionable is a type that represents a migration action.
type Actionable interface {
	Action(context.Context, *db.Connection, *sql.Tx) error
}

// Action is a function that can be run during a migration step.
type Action func(context.Context, *db.Connection, *sql.Tx) error

// NoOp performs no action.
func NoOp(_ context.Context, _ *db.Connection, _ *sql.Tx) error { return nil }

// Statements returns a body func that executes the statments serially.
func Statements(statements ...string) Action {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		for _, statement := range statements {
			err = c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Exec(statement)
			if err != nil {
				return
			}
		}
		return
	}
}

// Exec creates an Action that will run a statement with a given set of arguments.
// It can be used in lieu of Statements, when parameterization is needed
func Exec(statement string, args ...interface{}) Action {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		err = c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Exec(statement, args...)
		return
	}
}

// Actions creates an Action with a single body func that executes all the variadic argument actions serially
func Actions(actions ...Action) Action {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		for _, action := range actions {
			err = action(ctx, c, tx)
			if err != nil {
				return err
			}
		}
		return
	}
}
