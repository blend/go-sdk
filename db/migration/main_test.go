package migration

import (
	"fmt"
	"os"
	"testing"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"

	"github.com/blend/go-sdk/db"
	_ "github.com/lib/pq"
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

func buildTestSchemaName() string {
	return fmt.Sprintf("test_sch_%s", stringutil.Random(stringutil.LowerLetters, 10))
}
