package dbutil

import (
	"context"

	"github.com/blend/go-sdk/db"
)

// DatabaseExists returns if a database exists or not.
func DatabaseExists(ctx context.Context, name string, opts ...db.Option) (bool, error) {
	conn, err := OpenManagementConnection(opts...)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return conn.QueryContext(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", name).Any()
}
