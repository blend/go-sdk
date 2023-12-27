/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"context"
	"database/sql"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

// NewGroup creates a new Group from a given list of actionable.
func NewGroup(options ...GroupOption) *Group {
	g := Group{}
	for _, o := range options {
		o(&g)
	}
	return &g
}

// NewGroupWithAction returns a new group with a single action.
func NewGroupWithAction(guard GuardFunc, action Action, options ...GroupOption) *Group {
	return NewGroup(
		append([]GroupOption{OptGroupActions(NewStep(guard, action))}, options...)...,
	)
}

// Group is an series of migration actions.
// It uses normally transactions to apply these actions as an atomic unit, but this transaction can be bypassed by
// setting the SkipTransaction flag to true. This allows the use of CONCURRENT index creation and other operations that
// postgres will not allow within a transaction.
type Group struct {
	Actions         []Action
	Tx              *sql.Tx
	SkipTransaction bool
}

// Action runs the groups actions within a transaction.
func (ga *Group) Action(ctx context.Context, c *db.Connection) (err error) {
	var tx *sql.Tx
	if ga.Tx != nil { // if we have a transaction provided to us
		tx = ga.Tx
	} else if !ga.SkipTransaction { // if we aren't told to skip transactions
		tx, err = c.Begin()
		if err != nil {
			return
		}
		defer func() {
			if err != nil {
				if txErr := tx.Rollback(); txErr != nil {
					err = ex.Nest(err, txErr)
				}
			} else {
				if txErr := tx.Commit(); txErr != nil {
					err = ex.Nest(err, txErr)
				}
			}
		}()
	}

	for _, a := range ga.Actions {
		err = a.Action(ctx, c, tx)
		if err != nil {
			return
		}
	}

	return
}
