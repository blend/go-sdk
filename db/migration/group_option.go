/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import "database/sql"

// GroupOption is an option for migration Groups (Group)
type GroupOption func(g *Group)

// OptGroupActions allows you to add actions to the NewGroup. If you want, multiple OptActions can be applied to the same NewGroup.
// They are additive.
func OptGroupActions(actions ...Action) GroupOption {
	return func(g *Group) {
		if len(g.Actions) == 0 {
			g.Actions = actions
		} else {
			g.Actions = append(g.Actions, actions...)
		}
	}
}

// OptGroupSkipTransaction will allow this group to be run outside of a transaction. Use this to concurrently create indices
// and perform other actions that cannot be executed in a Tx
func OptGroupSkipTransaction() GroupOption {
	return func(g *Group) {
		g.SkipTransaction = true
	}
}

// OptGroupTx sets a transaction on the group.
func OptGroupTx(tx *sql.Tx) GroupOption {
	return func(g *Group) {
		g.Tx = tx
	}
}
