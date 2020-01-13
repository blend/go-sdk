package dbtrace

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
)

func TestMain(m *testing.M) {
	conn, err := db.New(db.OptConfigFromEnv())
	if err != nil {
		logger.FatalExit(err)
	}
	err = openDefaultDB(conn)
	if err != nil {
		logger.FatalExit(err)
	}
	defer conn.Close()
	os.Exit(m.Run())
}

var (
	defaultConnection *db.Connection
)

func setDefaultDB(conn *db.Connection) {
	defaultConnection = conn
}

func defaultDB() *db.Connection {
	return defaultConnection
}

func openDefaultDB(conn *db.Connection) error {
	err := conn.Open()
	if err != nil {
		return err
	}
	setDefaultDB(conn)
	return nil
}

func createTable(tx *sql.Tx) error {
	createSQL := `CREATE TABLE IF NOT EXISTS test_table (
		id serial not null primary key
	);`
	return db.IgnoreExecResult(defaultDB().Invoke(db.OptTx(tx)).Exec(createSQL))
}
