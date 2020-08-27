package main

import (
	"database/sql"
	"log"

	"github.com/blend/go-sdk/db"
)

func ignoreConnDone(err error) error {
	if err == sql.ErrConnDone {
		return nil
	}
	return err
}

func cleanUp(pool *db.Connection) {
	err := ignoreConnDone(pool.Close())
	if err != nil {
		log.Fatal(err)
	}
}
