package dbutil

import (
	"context"

	"github.com/blend/go-sdk/db"
)

// CloseAllConnections closes all other connections to a database.
func CloseAllConnections(ctx context.Context, conn *db.Connection, databaseName string) error {
	_, err := conn.Invoke(db.OptContext(ctx)).Exec(`select pg_terminate_backend(pid) from pg_stat_activity where datname = $1;`, databaseName)
	return err
}
