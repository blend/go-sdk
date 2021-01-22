package testutil

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/dbutil"
	"github.com/blend/go-sdk/uuid"
)

// CreateTestDatabase creates a randomized test database.
func CreateTestDatabase(ctx context.Context, opts ...db.Option) (*db.Connection, error) {
	databaseName := fmt.Sprintf("testdb_%s", uuid.V4().String())
	if err := dbutil.CreateDatabase(ctx, databaseName, opts...); err != nil {
		return nil, err
	}

	defaults := []db.Option{
		db.OptHost("localhost"),
		db.OptSSLMode("disable"),
		db.OptConfigFromEnv(),
		db.OptDatabase(databaseName),
	}
	conn, err := db.New(
		append(defaults, opts...)...,
	)
	if err != nil {
		return nil, err
	}
	err = conn.Open()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
