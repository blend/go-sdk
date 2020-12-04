package dbutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
)

// DatabaseExists returns if a database exists or not.
func DatabaseExists(ctx context.Context, name string, opts ...db.Option) error {
	conn, err := OpenManagementConnection(opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	any, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", name)).Any()
	if err != nil {
		return err
	}
	if !any {
		return fmt.Errorf("database does not exist: %s", name)
	}
	return nil
}
