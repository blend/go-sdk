/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/logger"
)

// UserNotExists creates a index on the given connection if it does not exist.
func UserNotExists(username string) migration.GuardFunc {
	return migration.Guard(fmt.Sprintf("create user `%s`", username), func(ctx context.Context, c *db.Connection, tx *sql.Tx) (bool, error) {
		return c.Invoke(db.OptContext(ctx), db.OptTx(tx)).Query(`SELECT 1 FROM users WHERE username = $1`, strings.ToLower(username)).None()
	})
}

func main() {
	suite := migration.New(
		migration.OptGroups(
			migration.NewGroupWithAction(
				migration.TableNotExists("users"),
				migration.Statements(
					"CREATE TABLE users (username varchar(255) primary key);",
				),
			),
			migration.NewGroupWithAction(
				UserNotExists("example-string"),
				migration.Exec("INSERT INTO users (username) VALUES ($1)", "example-string"),
			),
			migration.NewGroupWithAction(
				UserNotExists("example-string"),
				migration.Exec("INSERT INTO users (username) VALUES ($1)", "example-string"),
			),
			migration.NewGroupWithAction(
				migration.TableExists("users"),
				migration.Statements(
					"DROP TABLE users;",
				),
			),
		),
	)
	suite.Log = logger.All()

	conn, err := db.Open(db.New(db.OptConfigFromEnv()))
	if err != nil {
		logger.FatalExit(err)
	}

	if err := suite.Apply(context.Background(), conn); err != nil {
		logger.FatalExit(err)
	}
}
