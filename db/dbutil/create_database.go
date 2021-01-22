package dbutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
)

// CreateDatabase creates a database with a given name.
//
// Note: the `name` parameter is passed to the statement directly (not via. a parameter).
// You should use extreme care to not pass user submitted inputs to this function.
func CreateDatabase(ctx context.Context, name string, opts ...db.Option) error {
	conn, err := OpenManagementConnection(opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = ValidateDatabaseName(name); err != nil {
		return err
	}

	statement := fmt.Sprintf("CREATE DATABASE %s", name)
	_, err = conn.ExecContext(ctx, statement)
	return err
}
