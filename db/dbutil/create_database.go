package dbutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
)

// CreateDatabase creates a database with a given name.
func CreateDatabase(ctx context.Context, name string, opts ...db.Option) error {
	conn, err := OpenManagementConnection(opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s ", name))
	return err
}
