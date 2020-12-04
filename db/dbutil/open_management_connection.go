package dbutil

import "github.com/blend/go-sdk/db"

// OpenManagementConnection creates a database connection to the default database (typically postgres).
func OpenManagementConnection(options ...db.Option) (*db.Connection, error) {
	defaults := []db.Option{
		db.OptHost("localhost"),
		db.OptDatabase("postgres"),
		db.OptSSLMode("disable"),
		db.OptConfigFromEnv(),
	}
	conn, err := db.New(
		append(defaults, options...)...,
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
