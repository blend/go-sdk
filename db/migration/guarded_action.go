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

// NewStep returns a new Step, given a GuardFunc and an Action
func NewStep(guard GuardFunc, action Action) *Step {
	return &Step{
		Guard:	guard,
		Body:	action,
	}
}

// Step is a guarded action. The GuardFunc will decide whether to execute this Action
type Step struct {
	Guard	GuardFunc
	Body	Action
}

// Action implements the Actionable interface and runs the body if the provided guard passes.
func (ga *Step) Action(ctx context.Context, c *db.Connection, tx *sql.Tx) error {
	return ga.Guard(ctx, c, tx, ga.Body)
}
