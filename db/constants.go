/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"time"
)

const (
	// DefaultEngine is the default database engine.
	DefaultEngine	= "pgx"	// "postgres"

	// EnvVarDBEngine is the environment variable used to set the Go `sql` driver.
	EnvVarDBEngine	= "DB_ENGINE"
	// EnvVarDatabaseURL is the environment variable used to set the entire
	// database connection string.
	EnvVarDatabaseURL	= "DATABASE_URL"
	// EnvVarDBHost is the environment variable used to set the host in a
	// database connection string.
	EnvVarDBHost	= "DB_HOST"
	// EnvVarDBPort is the environment variable used to set the port in a
	// database connection string.
	EnvVarDBPort	= "DB_PORT"
	// EnvVarDBName is the environment variable used to set the database name
	// in a database connection string.
	EnvVarDBName	= "DB_NAME"
	// EnvVarDBSchema is the environment variable used to set the database
	// schema in a database connection string.
	EnvVarDBSchema	= "DB_SCHEMA"
	// EnvVarDBApplicationName is the environment variable used to set the
	// `application_name` configuration parameter in a `lib/pq` connection
	// string.
	//
	// See: https://www.postgresql.org/docs/12/runtime-config-logging.html#GUC-APPLICATION-NAME
	EnvVarDBApplicationName	= "DB_APPLICATION_NAME"
	// EnvVarDBUser is the environment variable used to set the user in a
	// database connection string.
	EnvVarDBUser	= "DB_USER"
	// EnvVarDBPassword is the environment variable used to set the password
	// in a database connection string.
	EnvVarDBPassword	= "DB_PASSWORD"
	// EnvVarDBConnectTimeout is is the environment variable used to set the
	// connect timeout in a database connection string.
	EnvVarDBConnectTimeout	= "DB_CONNECT_TIMEOUT"
	// EnvVarDBLockTimeout is is the environment variable used to set the lock
	// timeout on a database config.
	EnvVarDBLockTimeout	= "DB_LOCK_TIMEOUT"
	// EnvVarDBStatementTimeout is is the environment variable used to set the
	// statement timeout on a database config.
	EnvVarDBStatementTimeout	= "DB_STATEMENT_TIMEOUT"
	// EnvVarDBSSLMode is the environment variable used to set the SSL mode in
	// a database connection string.
	EnvVarDBSSLMode	= "DB_SSLMODE"
	// EnvVarDBIdleConnections is the environment variable used to set the
	// maximum number of idle connections allowed in a connection pool.
	EnvVarDBIdleConnections	= "DB_IDLE_CONNECTIONS"
	// EnvVarDBMaxConnections is the environment variable used to set the
	// maximum number of connections allowed in a connection pool.
	EnvVarDBMaxConnections	= "DB_MAX_CONNECTIONS"
	// EnvVarDBMaxLifetime is the environment variable used to set the maximum
	// lifetime of a connection in a connection pool.
	EnvVarDBMaxLifetime	= "DB_MAX_LIFETIME"
	// EnvVarDBMaxIdleTime is the environment variable used to set the maximum
	// time a connection can be idle.
	EnvVarDBMaxIdleTime	= "DB_MAX_IDLE_TIME"
	// EnvVarDBBufferPoolSize is the environment variable used to set the buffer
	// pool size on a connection in a connection pool.
	EnvVarDBBufferPoolSize	= "DB_BUFFER_POOL_SIZE"
	// EnvVarDBDialect is the environment variable used to set the dialect
	// on a connection configuration (e.g. `postgres` or `cockroachdb`).
	EnvVarDBDialect	= "DB_DIALECT"

	// DefaultHost is the default database hostname, typically used
	// when developing locally.
	DefaultHost	= "localhost"
	// DefaultPort is the default postgres port.
	DefaultPort	= "5432"
	// DefaultDatabase is the default database to connect to, we use
	// `postgres` to not pollute the template databases.
	DefaultDatabase	= "postgres"

	// DefaultSchema is the default schema to connect to
	DefaultSchema	= "public"

	// SSLModeDisable is an ssl mode.
	// Postgres Docs: "I don't care about security, and I don't want to pay the overhead of encryption."
	SSLModeDisable	= "disable"
	// SSLModeAllow is an ssl mode.
	// Postgres Docs: "I don't care about security, but I will pay the overhead of encryption if the server insists on it."
	SSLModeAllow	= "allow"
	// SSLModePrefer is an ssl mode.
	// Postgres Docs: "I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it"
	SSLModePrefer	= "prefer"
	// SSLModeRequire is an ssl mode.
	// Postgres Docs: "I want my data to be encrypted, and I accept the overhead. I trust that the network will make sure I always connect to the server I want."
	SSLModeRequire	= "require"
	// SSLModeVerifyCA is an ssl mode.
	// Postgres Docs: "I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server that I trust."
	SSLModeVerifyCA	= "verify-ca"
	// SSLModeVerifyFull is an ssl mode.
	// Postgres Docs: "I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server I trust, and that it's the one I specify."
	SSLModeVerifyFull	= "verify-full"

	// DefaultConnectTimeout is the default connect timeout.
	DefaultConnectTimeout	= 5 * time.Second

	// DefaultIdleConnections is the default number of idle connections.
	DefaultIdleConnections	= 16
	// DefaultMaxConnections is the default maximum number of connections.
	DefaultMaxConnections	= 32
	// DefaultMaxLifetime is the default maximum lifetime of driver connections.
	DefaultMaxLifetime	= time.Duration(0)
	// DefaultMaxIdleTime is the default maximum idle time of driver connections.
	DefaultMaxIdleTime	= time.Duration(0)
	// DefaultBufferPoolSize is the default number of buffer pool entries to maintain.
	DefaultBufferPoolSize	= 1024
)
