/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package testutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jackc/pgconn"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
)

const (
	passwordText		= "..password-redacted.."
	requireDBErrorTemplate	= `%s
%s
Connection String:
  %q
`
)

func redactEnvironmentVariable(key, value string) string {
	if key == db.EnvVarDBPassword {
		return passwordText
	}

	if key == db.EnvVarDatabaseURL {
		return createLoggingDSN(db.Config{DSN: value})
	}

	return value
}

// allDBEnvironmentVariables returns a slice of all the environment variables
// used in `sdk/db/config.go::Config.Resolve()`. The accepted list of
// environment variables may change there over time so this hardcoded list
// may drift.
func allDBEnvironmentVariables() []string {
	return []string{
		db.EnvVarDBEngine,
		db.EnvVarDatabaseURL,
		db.EnvVarDBHost,
		db.EnvVarDBPort,
		db.EnvVarDBName,
		db.EnvVarDBSchema,
		db.EnvVarDBApplicationName,
		db.EnvVarDBUser,
		db.EnvVarDBPassword,
		db.EnvVarDBConnectTimeout,
		db.EnvVarDBLockTimeout,
		db.EnvVarDBStatementTimeout,
		db.EnvVarDBSSLMode,
		db.EnvVarDBIdleConnections,
		db.EnvVarDBMaxConnections,
		db.EnvVarDBMaxLifetime,
		db.EnvVarDBBufferPoolSize,
		db.EnvVarDBDialect,
	}
}

func envResolveError(err error, ev env.Vars) error {
	found := []string{}
	for _, name := range allDBEnvironmentVariables() {
		value, ok := ev[name]
		if ok {
			line := fmt.Sprintf("- %s=%q", name, redactEnvironmentVariable(name, value))
			found = append(found, line)
		}
	}

	lines := []string{
		"Failed to read 'DB_*' environment variables. Error:",
		"%s",
	}
	if len(found) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Environment Variables:")
		lines = append(lines, found...)
	}

	lines = append(lines, "")	// Trailing newline
	template := strings.Join(lines, "\n")
	errString := fmt.Sprintf("%v", err)
	return ex.Class(fmt.Sprintf(template, indentTwo(errString)))
}

// ResolveDBConfig is intended to be used to help debug issues resolving
// a `db.Config` from the environment.
//
// In the case of failure, this wraps the `Resolve()` error with a helpful
// message and a list of all relevant environment variables.
func ResolveDBConfig(ctx context.Context, c *db.Config) error {
	ev := env.GetVars(ctx)
	err := c.Resolve(ctx)
	if err == nil {
		return nil
	}

	return envResolveError(err, ev)
}

func indentTwo(s string) string {
	lines := strings.Split(s, "\n")
	indented := make([]string, len(lines))
	for i, line := range lines {
		indented[i] = "  " + line
	}
	return strings.Join(indented, "\n")
}

func getSQLErrorMessage(err error) *string {
	errString := err.Error()
	// NOTE: The string-munging is partially because `errors.errorString` is
	//       not exported. We could instead get around this by using `reflect`
	//       to get the underlying package and type name. Additionally, these
	//       errors may be wrapped in an `ex.Ex` as `Class`.
	if strings.HasPrefix(errString, "sql: ") {
		withoutPrefix := strings.TrimPrefix(errString, "sql: ")
		return &withoutPrefix
	}

	return nil
}

// ValidatePool validates that
// - the connection string is valid
// - the selected `sql` driver can be used
// - a simple ping can be sent over the connection (is the DB reachable?)
//
// In the case of failure, this tries to diagnose the connection error and
// produce helpful tips on how to resolve.
func ValidatePool(ctx context.Context, pool *db.Connection, hints string) error {
	if pool == nil {
		return ex.New("Cannot validate a nil connection pool")
	}

	err := poolOpen(pool, hints)
	if err != nil {
		return err
	}

	return verifyConnect(ctx, pool, hints)
}

func formatKnownError(header, hints, dsn string) error {
	return ex.Class(fmt.Sprintf(requireDBErrorTemplate, header, hints, dsn))
}

func formatUnknownError(header, hints, dsn string) error {
	return ex.New(fmt.Sprintf(requireDBErrorTemplate, header, hints, dsn))
}

// poolOpen calls `Open()` to verify the connection string is valid and
// that the selected `sql` driver can be used.
func poolOpen(pool *db.Connection, hints string) error {
	// Early exit if the connection is already open.
	if pool.Connection != nil {
		return nil
	}

	err := pool.Open()
	if err == nil {
		return nil
	}

	dsn := createLoggingDSN(pool.Config)
	sqlErrorMessage := getSQLErrorMessage(err)
	if sqlErrorMessage != nil {
		header := fmt.Sprintf(
			"Error from 'sql' package:\n  %s\nDatabase Engine:\n  %s",
			*sqlErrorMessage, pool.Config.EngineOrDefault(),
		)
		return formatKnownError(header, hints, dsn)
	}

	errString := fmt.Sprintf("%+v", err)
	header := fmt.Sprintf("Unexpected Open() failure:\n%s", indentTwo(errString))
	return formatUnknownError(header, hints, dsn)
}

func unwrapNetOpError(err error) *net.OpError {
	noe, ok := err.(*net.OpError)
	if ok {
		return noe
	}

	ue := errors.Unwrap(err)
	noe, ok = ue.(*net.OpError)
	if ok {
		return noe
	}

	return nil
}

func isConnectionRefused(err error) bool {
	noe := unwrapNetOpError(err)
	if noe == nil {
		return false
	}

	// NOTE: We could go deeper in here by type asserting `noe.Err` as an
	//       `*os.SyscallError` and checking for `syscall.ECONNREFUSED`.
	//       The string `connect: connection refused` has been verified in
	//       Go 1.12, 1.13, 1.14, 1.15 on macOS and Alpine Linux but may change
	//       in future releases.
	return noe.Err.Error() == "connect: connection refused"
}

func getPGXErrorMessage(err error) *string {
	pe, ok := err.(*pgconn.PgError)
	if ok {
		return &pe.Message
	}

	ue := errors.Unwrap(err)
	pe, ok = ue.(*pgconn.PgError)
	if ok {
		return &pe.Message
	}

	errString := err.Error()
	// NOTE: The string-munging is partially because `pgconn.connectError` is
	//       not exported.
	if strings.HasPrefix(errString, "failed to connect to `host=") {
		wrappedErrString := ue.Error()
		return &wrappedErrString
	}

	return nil
}

// verifyConnect verifies that the target database is actually running and the
// connection pool can actually connect.
func verifyConnect(ctx context.Context, pool *db.Connection, hints string) error {
	err := pool.Connection.PingContext(ctx)
	if err == nil {
		return nil
	}

	dsn := createLoggingDSN(pool.Config)
	if isConnectionRefused(err) {
		header := "Network error:\n  Could not connect to database."
		return formatKnownError(header, hints, dsn)
	}

	pgxErrorMessage := getPGXErrorMessage(err)
	if pgxErrorMessage != nil {
		header := fmt.Sprintf("PostgreSQL error when connecting to the database:\n  %s", *pgxErrorMessage)
		return formatKnownError(header, hints, dsn)
	}

	errString := fmt.Sprintf("%+v", err)
	header := fmt.Sprintf("Unexpected PingContext() failure:\n%s", indentTwo(errString))
	return formatUnknownError(header, hints, dsn)
}

func createLoggingDSN(c db.Config) string {
	if c.DSN != "" {
		nc, err := db.NewConfigFromDSN(c.DSN)
		if err != nil {
			return "Failed to parse DSN: see DATABASE_URL environment variable"
		}
		return createLoggingDSN(nc)
	}

	dsn := c.CreateLoggingDSN()
	if c.Username == "" || c.Password == "" {
		return dsn
	}

	parts := strings.SplitN(dsn, "@", 2)
	if len(parts) != 2 {
		return dsn
	}

	return fmt.Sprintf("%s:%s@%s", parts[0], passwordText, parts[1])
}
