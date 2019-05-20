package main

import (
	"context"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/logger"
)

func main() {
	suite := migration.New(
		migration.Group(
			migration.Step(
				migration.Always(),
				migration.Statements(
					"select pg_sleep(10);",
				),
				db.OptTimeout(500*time.Millisecond),
			),
		),
	)

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
