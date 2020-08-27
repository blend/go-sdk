package main

import (
	"database/sql"
	"log"

	"github.com/blend/go-sdk/db"
)

func ignoreTxDone(err error) error {
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

func ignoreConnDone(err error) error {
	if err == sql.ErrConnDone {
		return nil
	}
	return err
}

func txFinalize(tx *sql.Tx, err error) error {
	if tx == nil {
		return err
	}

	rollbackErr := ignoreTxDone(tx.Rollback())
	return nest(err, rollbackErr)
}

func cleanUp(pool *db.Connection) {
	err := ignoreConnDone(pool.Close())
	if err != nil {
		log.Fatal(err)
	}
}
