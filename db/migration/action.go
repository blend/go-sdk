/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
)

// Action is a type that represents a migration action.
type Action interface {
	Action(context.Context, *db.Connection, *sql.Tx) error
}

// ActionFunc is a function that can be run during a migration step.
type ActionFunc func(context.Context, *db.Connection, *sql.Tx) error

// Action implements actioner.
func (a ActionFunc) Action(ctx context.Context, conn *db.Connection, tx *sql.Tx) error {
	return a(ctx, conn, tx)
}

// NoOp performs no action.
func NoOp(_ context.Context, _ *db.Connection, _ *sql.Tx) error { return nil }

// Statements returns a body func that executes the statments serially.
func Statements(statements ...string) Action {
	return ActionFunc(func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		for _, statement := range statements {
			err = db.IgnoreExecResult(c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Exec(statement))
			if err != nil {
				return
			}
		}
		return
	})
}

// Exec creates an Action that will run a statement with a given set of arguments.
// It can be used in lieu of Statements, when parameterization is needed
func Exec(statement string, args ...interface{}) Action {
	return ActionFunc(func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		err = db.IgnoreExecResult(c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Exec(statement, args...))
		return
	})
}

// Actions creates an Action with a single body func that executes all the variadic argument actions serially
func Actions(actions ...Action) Action {
	return ActionFunc(func(ctx context.Context, c *db.Connection, tx *sql.Tx) (err error) {
		for _, action := range actions {
			err = action.Action(ctx, c, tx)
			if err != nil {
				return err
			}
		}
		return
	})
}
