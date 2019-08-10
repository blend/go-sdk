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

// NewGroupWithActions is a helper for "migrations.NewGroup(migrations.OptActions(actions ...migration.Actionable))"
func NewGroupWithActions(actions ...Actionable) *Group {
	return NewGroup(OptActions(actions...))
}

// NewGroupWithAction returns a new group with a single action.
func NewGroupWithAction(guard GuardFunc, action Action) *Group {
	return NewGroup(OptActions(NewStep(guard, action)))
}

// Group is an series of migration actions.
// It uses normally transactions to apply these actions as an atomic unit, but this transaction can be bypassed by
// setting the SkipTransaction flag to true. This allows the use of CONCURRENT index creation and other operations that
// postgres will not allow within a transaction.
type Group struct {
	Actions         []Actionable
	SkipTransaction bool
}

// Action runs the groups actions within a transaction.
func (ga *Group) Action(ctx context.Context, c *db.Connection) (err error) {
	var tx *sql.Tx
	if !ga.SkipTransaction {
		tx, err = c.Begin()
		if err != nil {
			return
		}
	}

	// commit or rollback the transaction if this is a transactional operation
	if tx != nil {
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
