package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/logger"
)

func main() {
	suite := migration.New(migration.OptGroups(
		migration.NewGroup(migration.OptActions(
			migration.NewStep(
				migration.Always(),
				func(ctx context.Context, connection *db.Connection, tx *sql.Tx) error {
					return connection.Invoke(db.OptTimeout(500 * time.Millisecond)).Exec("select pg_sleep(10);")
				},
			),
		),
		)))

	suite.Log = logger.Prod()

	conn, err := db.Open(db.New(db.OptConfigFromEnv()))
	if err != nil {
		logger.FatalExit(err)
	}
	suite.Log.Info("starting migrations")
	outerTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := suite.Apply(outerTimeout, conn); err != nil {
		logger.FatalExit(err)
	}
}
