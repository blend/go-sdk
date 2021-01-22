package dbutil

import "github.com/blend/go-sdk/db"

// OpenManagementConnection creates a database connection to the default database (typically postgres).
func OpenManagementConnection(options ...db.Option) (*db.Connection, error) {
	defaults := []db.Option{
		db.OptHost("localhost"),
		db.OptSSLMode("disable"),
		db.OptConfigFromEnv(),
		db.OptDatabase("postgres"),
	}
	conn, err := db.New(
		append(defaults, append(options, db.OptDatabase("postgres"))...)...,
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
