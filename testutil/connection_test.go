/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/testutil"
	"github.com/blend/go-sdk/uuid"
)

var (
	// actualEnv contains the **actual** / current environment variables when
	// the tests were started. This enables tests to run in parallel without
	// needing to globally lock the result of `env.Env()`.
	actualEnv = env.Env()
)

const (
	dbDefaultUsername = "postgres"
	ipv4Loopback      = "127.0.0.1"
)

func TestResolveDBConfig_InvalidEnv(t *testing.T) {
	it := assert.New(t)

	ev := env.New()
	// Use only the values we set; e.g. in CI `DB_HOST` and `DB_PORT` also get set.
	ev.Set(db.EnvVarDatabaseURL, "postgres://hi:bye@localhost:9999/any")
	ev.Set(db.EnvVarDBPassword, "bye")
	ev.Set(db.EnvVarDBConnectTimeout, "not-int")

	c := db.Config{}
	ctx := env.WithVars(context.TODO(), ev)
	err := testutil.ResolveDBConfig(ctx, &c)
	it.NotNil(err)
	expected := `Failed to read 'DB_*' environment variables. Error:
  time: invalid duration "not-int"

Environment Variables:
- DATABASE_URL="postgres://hi:..password-redacted..@localhost:9999/any"
- DB_PASSWORD="..password-redacted.."
- DB_CONNECT_TIMEOUT="not-int"
`
	it.Equal(expected, fmt.Sprintf("%v", err))
}

func TestValidatePool_NotRunning(t *testing.T) {
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "pgx"

	// Use a random port and ":fingers_crossed:" it isn't open on the machine.
	port := getRandomPort()
	c.Port = fmt.Sprintf("%d", port)

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `Network error:
  Could not connect to database.

advice on restart

Connection String:
  "postgres://%[1]s:..password-redacted..@%[2]s:%[3]d/%[4]s?connect_timeout=5&search_path=%[5]s&sslmode=disable"
`
	it.Equal(formatExpectedPort(expected, port, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_NoSSL(t *testing.T) {
	t.Skip("it is not a safe assumption that all local installations of postgres disable ssl")
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "pgx"
	c.SSLMode = db.SSLModeRequire

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	expected := `PostgreSQL error when connecting to the database:
  server refused TLS connection

advice on restart

Connection String:
  "postgres://%[1]s:..password-redacted..@%[2]s:%[3]d/%[4]s?connect_timeout=5&search_path=%[5]s&sslmode=require"
`
	it.Equal(formatExpected(expected, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_MissingPassword(t *testing.T) {
	t.Skip("it is not a safe assumption that all local installations of postgres require password authentication")
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "pgx"
	c.Password = ""

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `PostgreSQL error when connecting to the database:
  password authentication failed for user "%[1]s"

advice on restart

Connection String:
  "postgres://%[1]s@%[2]s:%[3]d/%[4]s?connect_timeout=5&search_path=%[5]s&sslmode=disable"
`
	it.Equal(formatExpected(expected, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_WrongPassword(t *testing.T) {
	t.Skip("it is not a safe assumption that all local installations of postgres require password authentication")
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "pgx"
	c.Password = uuid.V4().String()

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `PostgreSQL error when connecting to the database:
  password authentication failed for user "%[1]s"

advice on restart

Connection String:
  "postgres://%[1]s:..password-redacted..@%[2]s:%[3]d/%[4]s?connect_timeout=5&search_path=%[5]s&sslmode=disable"
`
	it.Equal(formatExpected(expected, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_UnsupportedDriver(t *testing.T) {
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "not-pgx"
	c.Password = uuid.V4().ToShortString()

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `Error from 'sql' package:
  unknown driver "not-pgx" (forgotten import?)
Database Engine:
  not-pgx

advice on restart

Connection String:
  "postgres://%[1]s:..password-redacted..@%[2]s:%[3]d/%[4]s?connect_timeout=5&search_path=%[5]s&sslmode=disable"
`
	it.Equal(formatExpected(expected, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_NotAccepting(t *testing.T) {
	it := assert.New(t)

	c := defaultConfig(it, actualEnv)
	c.Engine = "pgx"
	// NOTE: This makes two strong assumptions that are somewhat specific to
	//       the default `postgres` Docker container:
	//       - the `template0` DB exists in the running PostgreSQL instance
	//       - the `template0` DB is not accepting connections.
	c.Database = "template0"

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `PostgreSQL error when connecting to the database:
  database "template0" is not currently accepting connections

advice on restart

Connection String:
  "postgres://%[1]s:..password-redacted..@%[2]s:%[3]d/template0?connect_timeout=5&search_path=%[5]s&sslmode=disable"
`
	it.Equal(formatExpected(expected, actualEnv), fmt.Sprintf("%v", err))
}

func TestValidatePool_DoesNotExist(t *testing.T) {
	it := assert.New(t)

	database := fmt.Sprintf("testdb_%s", uuid.V4().String())
	overrides := env.New()
	overrides.Set(db.EnvVarDBName, database)
	combined := env.Merge(actualEnv, overrides)

	c := defaultConfig(it, combined)
	c.Engine = "pgx"
	c.Database = database

	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)
	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expectedTemplate := `PostgreSQL error when connecting to the database:
  database "%s" does not exist

advice on restart

Connection String:
  "postgres://%%[1]s:..password-redacted..@%%[2]s:%%[3]d/%%[4]s?connect_timeout=5&search_path=%%[5]s&sslmode=disable"
`
	expected := fmt.Sprintf(expectedTemplate, c.Database)
	it.Equal(formatExpected(expected, combined), fmt.Sprintf("%v", err))
}

func TestValidatePool_InvalidProtocolInDSN(t *testing.T) {
	it := assert.New(t)

	c := db.Config{Engine: "pgx", DSN: "x://y"}
	pool, err := db.New(db.OptConfig(c))
	it.Nil(err)

	err = testutil.ValidatePool(context.TODO(), pool, "\nadvice on restart\n")
	it.NotNil(err)
	expected := `Unexpected Open() failure:
  invalid connection protocol: x

advice on restart

Connection String:
  "Failed to parse DSN: see DATABASE_URL environment variable"
`
	it.Equal(expected, fmt.Sprintf("%v", err))
}

func getRandomPort() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return 2048 + r.Intn(4096)
}

// defaultConfig returns a database config that sets a few defaults and then
// uses environment variables to populate the rest.
func defaultConfig(it *assert.Assertions, ev env.Vars) db.Config {
	c := db.Config{
		Schema:   db.DefaultSchema,
		Username: dbDefaultUsername,
		Password: "s33kr1t",
		SSLMode:  db.SSLModeDisable,
	}
	ctx := env.WithVars(context.TODO(), ev)
	err := c.Resolve(ctx)
	it.Nil(err)
	// Ensure IPv4 address is used; in some environments where IPv6 may be
	// a problem this can result in unexpected failures of the form
	// > dial tcp [::1]:5432: connect: cannot assign requested address
	if c.Host == db.DefaultHost {
		c.Host = ipv4Loopback
	}
	return c
}

func defaultHost(ev env.Vars) string {
	host, ok := ev[db.EnvVarDBHost]
	if !ok {
		// Ensure IPv4 address is used
		return ipv4Loopback
	}
	// Ensure IPv4 address is used
	if host == db.DefaultHost {
		return ipv4Loopback
	}

	return host
}

func defaultUsername(ev env.Vars) string {
	username, ok := ev[db.EnvVarDBUser]
	if !ok {
		return dbDefaultUsername
	}

	return username
}

func defaultDatabase(ev env.Vars) string {
	name, ok := ev[db.EnvVarDBName]
	if !ok {
		return db.DefaultDatabase
	}

	return name
}

func defaultSchema(ev env.Vars) string {
	schema, ok := ev[db.EnvVarDBSchema]
	if !ok {
		return db.DefaultSchema
	}

	return schema
}

// formatExpectedPort determines the DSN parameters to be used in a template string
// to populate an expected error template and then populates the template. s
func formatExpectedPort(template string, port int, ev env.Vars) string {
	host := defaultHost(ev)
	username := defaultUsername(ev)
	database := defaultDatabase(ev)
	schema := defaultSchema(ev)
	return fmt.Sprintf(template, username, host, port, database, schema)
}

func defaultPort(ev env.Vars) int {
	port, err := strconv.Atoi(ev.String(db.EnvVarDBPort))
	if err != nil {
		return 5432
	}

	return port
}

// formatExpected determines the DSN parameters to be used in a template string
// to populate an expected error template and then populates the template.
func formatExpected(template string, ev env.Vars) string {
	port := defaultPort(ev)
	return formatExpectedPort(template, port, ev)
}
