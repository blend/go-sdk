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
	return migration.Guard(fmt.Sprintf("create user `%s`", username), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return c.Invoke(context.TODO(), tx).Query(`SELECT 1 FROM users WHERE username = $1`, strings.ToLower(username)).None()
	})
}

func main() {
	suite := migration.New(
		migration.Group(
			migration.Step(
				migration.TableNotExists("users"),
				migration.Statements(
					"CREATE TABLE users (username varchar(255) primary key);",
				),
			),
		),
		migration.Group(
			migration.Step(
				UserNotExists("bailey"),
				migration.Exec("INSERT INTO users (username) VALUES ($1)", "bailey"),
			),
		),
		migration.Group(
			migration.Step(
				UserNotExists("bailey"),
				migration.Exec("INSERT INTO users (username) VALUES ($1)", "bailey"),
			),
		),
		migration.Group(
			migration.Step(
				migration.TableExists("users"),
				migration.Statements(
					"DROP TABLE users;",
				),
			),
		),
	).WithLogger(logger.All())

	conn, err := db.Open(db.NewFromConfig(db.MustNewConfigFromEnv()))
	if err != nil {
		logger.FatalExit(err)
	}

	if err := suite.Apply(conn); err != nil {
		logger.FatalExit(err)
	}
}
