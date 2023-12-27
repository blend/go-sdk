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

// GuardFunc is a control for migration steps.
// It should internally evaluate if the action should be called.
// The action is typically given separately so these two components can be composed.
type GuardFunc func(context.Context, *db.Connection, *sql.Tx, Action) error

// GuardPredicateFunc is a function that can act as a guard
type GuardPredicateFunc func(context.Context, *db.Connection, *sql.Tx) (bool, error)

// --------------------------------------------------------------------------------
// Guards
// --------------------------------------------------------------------------------

// Guard returns a function that determines if a step in a group should run.
func Guard(description string, predicate GuardPredicateFunc) GuardFunc {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx, step Action) error {
		proceed, err := predicate(ctx, c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}

		if !proceed {
			if suite := GetContextSuite(ctx); suite != nil {
				suite.Skipf(ctx, description)
			}
			return nil
		}

		err = step.Action(ctx, c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}
		if suite := GetContextSuite(ctx); suite != nil {
			suite.Applyf(ctx, description)
		}
		return nil
	}
}

// guardPredicate wraps a predicate in a GuardFunc
func guardPredicate(description string, p predicate, arg1 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(ctx, c, tx, arg1)
	})
}

// guardNotPredicate inverts a predicate, and wraps that in a GuardFunc
func guardNotPredicate(description string, p predicate, arg1 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(ctx, c, tx, arg1))
	})
}

// guardPredicate2 wraps a predicate2 in a GuardFunc
func guardPredicate2(description string, p predicate2, arg1, arg2 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(ctx, c, tx, arg1, arg2)
	})
}

// guardNotPredicate2 inverts a predicate2, and wraps that in a GuardFunc
func guardNotPredicate2(description string, p predicate2, arg1, arg2 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(ctx, c, tx, arg1, arg2))
	})
}

// guardPredicate3 wraps a predicate3 in a GuardFunc
func guardPredicate3(description string, p predicate3, arg1, arg2, arg3 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return p(ctx, c, tx, arg1, arg2, arg3)
	})
}

// guardNotPredicate3 inverts a predicate3, and wraps that in a GuardFunc
func guardNotPredicate3(description string, p predicate3, arg1, arg2, arg3 string) GuardFunc {
	return Guard(description, func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(p(ctx, c, tx, arg1, arg2, arg3))
	})
}

// predicate is a function that evaluates based on a string param.
type predicate func(context.Context, *db.Connection, *sql.Tx, string) (bool, error)

// predicate2 is a function that evaluates based on two string params.
type predicate2 func(context.Context, *db.Connection, *sql.Tx, string, string) (bool, error)

// predicate3 is a function that evaluates based on three string params.
type predicate3 func(context.Context, *db.Connection, *sql.Tx, string, string, string) (bool, error)
