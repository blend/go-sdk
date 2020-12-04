package db

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

// NewConfigFromDSN creates a new config from a DSN.
// Errors can be produced by parsing the DSN.
func NewConfigFromDSN(dsn string) (config Config, err error) {
	parsed, parseErr := ParseURL(dsn)
	if parseErr != nil {
		err = ex.New(parseErr)
		return
	}

	pieces := stringutil.SplitSpace(parsed)
	for _, piece := range pieces {
		if strings.HasPrefix(piece, "host=") {
			config.Host = strings.TrimPrefix(piece, "host=")
		} else if strings.HasPrefix(piece, "port=") {
			config.Port = strings.TrimPrefix(piece, "port=")
		} else if strings.HasPrefix(piece, "dbname=") {
			config.Database = strings.TrimPrefix(piece, "dbname=")
		} else if strings.HasPrefix(piece, "user=") {
			config.Username = strings.TrimPrefix(piece, "user=")
		} else if strings.HasPrefix(piece, "password=") {
			config.Password = strings.TrimPrefix(piece, "password=")
		} else if strings.HasPrefix(piece, "sslmode=") {
			config.SSLMode = strings.TrimPrefix(piece, "sslmode=")
		} else if strings.HasPrefix(piece, "search_path=") {
			config.Schema = strings.TrimPrefix(piece, "search_path=")
		} else if strings.HasPrefix(piece, "application_name=") {
			config.ApplicationName = strings.TrimPrefix(piece, "application_name=")
		} else if strings.HasPrefix(piece, "connect_timeout=") {
			timeout, parseErr := strconv.Atoi(strings.TrimPrefix(piece, "connect_timeout="))
			if parseErr != nil {
				err = ex.New(parseErr, ex.OptMessage("field: connect_timeout"))
				return
			}
			config.ConnectTimeout = time.Second * time.Duration(timeout)
		} else if strings.HasPrefix(piece, "lock_timeout=") {
			config.LockTimeout, parseErr = time.ParseDuration(strings.TrimPrefix(piece, "lock_timeout="))
			if parseErr != nil {
				err = ex.New(parseErr, ex.OptMessage("field: lock_timeout"))
				return
			}
		} else if strings.HasPrefix(piece, "statement_timeout=") {
			config.StatementTimeout, parseErr = time.ParseDuration(strings.TrimPrefix(piece, "statement_timeout="))
			if parseErr != nil {
				err = ex.New(parseErr, ex.OptMessage("field: statement_timeout"))
				return
			}
		}
	}

	return
}

// NewConfigFromEnv returns a new config from the environment.
// The environment variable mappings are as follows:
//  -  DB_ENGINE            = Engine
//  -  DATABASE_URL         = DSN     //note that this has precedence over other vars (!!)
//  -  DB_HOST              = Host
//  -  DB_PORT              = Port
//  -  DB_NAME              = Database
//  -  DB_SCHEMA            = Schema
//  -  DB_APPLICATION_NAME  = ApplicationName
//  -  DB_USER              = Username
//  -  DB_PASSWORD          = Password
//  -  DB_CONNECT_TIMEOUT   = ConnectTimeout
//  -  DB_LOCK_TIMEOUT      = LockTimeout
//  -  DB_STATEMENT_TIMEOUT = StatementTimeout
//  -  DB_SSLMODE           = SSLMode
//  -  DB_IDLE_CONNECTIONS  = IdleConnections
//  -  DB_MAX_CONNECTIONS   = MaxConnections
//  -  DB_MAX_LIFETIME      = MaxLifetime
//  -  DB_BUFFER_POOL_SIZE  = BufferPoolSize
func NewConfigFromEnv() (config Config, err error) {
	if err = (&config).Resolve(env.WithVars(context.Background(), env.Env())); err != nil {
		return
	}
	return
}

// MustNewConfigFromEnv returns a new config from the environment,
// it will panic if there is an error.
func MustNewConfigFromEnv() Config {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Config is a set of connection config options.
type Config struct {
	// Engine is the database engine.
	Engine string `json:"engine,omitempty" yaml:"engine,omitempty" env:"DB_ENGINE"`
	// DSN is a fully formed DSN (this skips DSN formation from all other variables outside `schema`).
	DSN string `json:"dsn,omitempty" yaml:"dsn,omitempty" env:"DATABASE_URL"`
	// Host is the server to connect to.
	Host string `json:"host,omitempty" yaml:"host,omitempty" env:"DB_HOST"`
	// Port is the port to connect to.
	Port string `json:"port,omitempty" yaml:"port,omitempty" env:"DB_PORT"`
	// DBName is the database name
	Database string `json:"database,omitempty" yaml:"database,omitempty" env:"DB_NAME"`
	// Schema is the application schema within the database, defaults to `public`. This schema is used to set the
	// Postgres "search_path" If you want to reference tables in other schemas, you'll need to specify those schemas
	// in your queries e.g. "SELECT * FROM schema_two.table_one..."
	// Using the public schema in a production application is considered bad practice as newly created roles will have
	// visibility into this data by default. We strongly recommend specifying this option and using a schema that is
	// owned by your service's role
	// We recommend against setting a multi-schema search_path, but if you really want to, you provide multiple comma-
	// separated schema names as the value for this config, or you can dbc.Invoke().Exec a SET statement on a newly
	// opened connection such as "SET search_path = 'schema_one,schema_two';" Again, we recommend against this practice
	// and encourage you to specify schema names beyond the first in your queries.
	Schema string `json:"schema,omitempty" yaml:"schema,omitempty" env:"DB_SCHEMA"`
	// ApplicationName is the name set by an application connection to a database
	// server, intended to be transmitted in the connection string. It can be
	// used to uniquely identify an application and will be included in the
	// `pg_stat_activity` view.
	//
	// See: https://www.postgresql.org/docs/12/runtime-config-logging.html#GUC-APPLICATION-NAME
	ApplicationName string `json:"applicationName,omitempty" yaml:"applicationName,omitempty" env:"DB_APPLICATION_NAME"`
	// Username is the username for the connection via password auth.
	Username string `json:"username,omitempty" yaml:"username,omitempty" env:"DB_USER"`
	// Password is the password for the connection via password auth.
	Password string `json:"password,omitempty" yaml:"password,omitempty" env:"DB_PASSWORD"`
	// ConnectTimeout determines the maximum wait for connection. The minimum
	// allowed timeout is 2 seconds, so anything below is treated the same
	// as unset. PostgreSQL will only accept second precision so this value will be
	// rounded to the nearest second before being set on a connection string.
	// Use `Validate()` to confirm that `ConnectTimeout` is exact to second
	// precision.
	//
	// See: https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-CONNECT-CONNECT-TIMEOUT
	ConnectTimeout time.Duration `json:"connectTimeout" yaml:"connectTimeout" env:"DB_CONNECT_TIMEOUT"`
	// LockTimeout is the timeout to use when attempting to acquire a lock.
	// PostgreSQL will only accept millisecond precision so this value will be
	// rounded to the nearest millisecond before being set on a connection string.
	// Use `Validate()` to confirm that `LockTimeout` is exact to millisecond
	// precision.
	//
	// See: https://www.postgresql.org/docs/current/runtime-config-client.html#GUC-LOCK-TIMEOUT
	LockTimeout time.Duration `json:"lockTimeout" yaml:"lockTimeout" env:"DB_LOCK_TIMEOUT"`
	// StatementTimeout is the timeout to use when invoking a SQL statement.
	// PostgreSQL will only accept millisecond precision so this value will be
	// rounded to the nearest millisecond before being set on a connection string.
	// Use `Validate()` to confirm that `StatementTimeout` is exact to millisecond
	// precision.
	//
	// See: https://www.postgresql.org/docs/current/runtime-config-client.html#GUC-STATEMENT-TIMEOUT
	StatementTimeout time.Duration `json:"statementTimeout" yaml:"statementTimeout" env:"DB_STATEMENT_TIMEOUT"`
	// SSLMode is the sslmode for the connection.
	SSLMode string `json:"sslMode,omitempty" yaml:"sslMode,omitempty" env:"DB_SSLMODE"`
	// IdleConnections is the number of idle connections.
	IdleConnections int `json:"idleConnections,omitempty" yaml:"idleConnections,omitempty" env:"DB_IDLE_CONNECTIONS"`
	// MaxConnections is the maximum number of connections.
	MaxConnections int `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty" env:"DB_MAX_CONNECTIONS"`
	// MaxLifetime is the maximum time a connection can be open.
	MaxLifetime time.Duration `json:"maxLifetime,omitempty" yaml:"maxLifetime,omitempty" env:"DB_MAX_LIFETIME"`
	// BufferPoolSize is the number of query composition buffers to maintain.
	BufferPoolSize int `json:"bufferPoolSize,omitempty" yaml:"bufferPoolSize,omitempty" env:"DB_BUFFER_POOL_SIZE"`
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return c.DSN == "" && c.Host == "" && c.Port == "" && c.Database == "" && c.Schema == "" && c.Username == "" && c.Password == "" && c.SSLMode == ""
}

// Resolve applies any external data sources to the config.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Engine, configutil.Env("DB_ENGINE"), configutil.String(c.Engine), configutil.String(DefaultEngine)),
		configutil.SetString(&c.DSN, configutil.Env("DATABASE_URL"), configutil.String(c.DSN)),
		configutil.SetString(&c.Host, configutil.Env("DB_HOST"), configutil.String(c.Host), configutil.String(DefaultHost)),
		configutil.SetString(&c.Port, configutil.Env("DB_PORT"), configutil.String(c.Port), configutil.String(DefaultPort)),
		configutil.SetString(&c.Database, configutil.Env("DB_NAME"), configutil.String(c.Database), configutil.String(DefaultDatabase)),
		configutil.SetString(&c.Schema, configutil.Env("DB_SCHEMA"), configutil.String(c.Schema)),
		configutil.SetString(&c.ApplicationName, configutil.Env(EnvVarDBApplicationName), configutil.String(c.ApplicationName)),
		configutil.SetString(&c.Username, configutil.Env("DB_USER"), configutil.String(c.Username), configutil.Env("USER")),
		configutil.SetString(&c.Password, configutil.Env("DB_PASSWORD"), configutil.String(c.Password)),
		configutil.SetDuration(&c.ConnectTimeout, configutil.Env("DB_CONNECT_TIMEOUT"), configutil.Duration(c.ConnectTimeout), configutil.Duration(DefaultConnectTimeout)),
		configutil.SetDuration(&c.LockTimeout, configutil.Env("DB_LOCK_TIMEOUT"), configutil.Duration(c.LockTimeout)),
		configutil.SetDuration(&c.StatementTimeout, configutil.Env("DB_STATEMENT_TIMEOUT"), configutil.Duration(c.StatementTimeout)),
		configutil.SetString(&c.SSLMode, configutil.Env("DB_SSLMODE"), configutil.String(c.SSLMode)),
		configutil.SetInt(&c.IdleConnections, configutil.Env("DB_IDLE_CONNECTIONS"), configutil.Int(c.IdleConnections), configutil.Int(DefaultIdleConnections)),
		configutil.SetInt(&c.MaxConnections, configutil.Env("DB_MAX_CONNECTIONS"), configutil.Int(c.MaxConnections), configutil.Int(DefaultMaxConnections)),
		configutil.SetDuration(&c.MaxLifetime, configutil.Env("DB_MAX_LIFETIME"), configutil.Duration(c.MaxLifetime), configutil.Duration(DefaultMaxLifetime)),
		configutil.SetInt(&c.BufferPoolSize, configutil.Env("DB_BUFFER_POOL_SIZE"), configutil.Int(c.BufferPoolSize), configutil.Int(DefaultBufferPoolSize)),
	)
}

// Reparse creates a DSN and reparses it, in case some values need to be coalesced.
func (c Config) Reparse() (Config, error) {
	return NewConfigFromDSN(c.CreateDSN())
}

// MustReparse creates a DSN and reparses it, in case some values need to be coalesced,
// and panics if there is an error.
func (c Config) MustReparse() Config {
	cfg, err := NewConfigFromDSN(c.CreateDSN())
	if err != nil {
		panic(err)
	}
	return cfg
}

// EngineOrDefault returns the database engine.
func (c Config) EngineOrDefault() string {
	if c.Engine != "" {
		return c.Engine
	}
	return DefaultEngine
}

// HostOrDefault returns the postgres host for the connection or a default.
func (c Config) HostOrDefault() string {
	if c.Host != "" {
		return c.Host
	}
	return DefaultHost
}

// PortOrDefault returns the port for a connection if it is not the standard postgres port.
func (c Config) PortOrDefault() string {
	if c.Port != "" {
		return c.Port
	}
	return DefaultPort
}

// DatabaseOrDefault returns the connection database or a default.
func (c Config) DatabaseOrDefault() string {
	if c.Database != "" {
		return c.Database
	}
	return DefaultDatabase
}

// SchemaOrDefault returns the schema on the search_path or the default ("public"). It's considered bad practice to
// use the public schema in production
func (c Config) SchemaOrDefault() string {
	if c.Schema != "" {
		return c.Schema
	}
	return DefaultSchema
}

// IdleConnectionsOrDefault returns the number of idle connections or a default.
func (c Config) IdleConnectionsOrDefault() int {
	if c.IdleConnections > 0 {
		return c.IdleConnections
	}
	return DefaultIdleConnections
}

// MaxConnectionsOrDefault returns the maximum number of connections or a default.
func (c Config) MaxConnectionsOrDefault() int {
	if c.MaxConnections > 0 {
		return c.MaxConnections
	}
	return DefaultMaxConnections
}

// MaxLifetimeOrDefault returns the maximum lifetime of a driver connection.
func (c Config) MaxLifetimeOrDefault() time.Duration {
	if c.MaxLifetime > 0 {
		return c.MaxLifetime
	}
	return DefaultMaxLifetime
}

// BufferPoolSizeOrDefault returns the number of query buffers to maintain or a default.
func (c Config) BufferPoolSizeOrDefault() int {
	if c.BufferPoolSize > 0 {
		return c.BufferPoolSize
	}
	return DefaultBufferPoolSize
}

// CreateDSN creates a postgres connection string from the config.
func (c Config) CreateDSN() string {
	if c.DSN != "" {
		return c.DSN
	}

	host := c.HostOrDefault()
	if c.PortOrDefault() != "" {
		host = host + ":" + c.PortOrDefault()
	}

	dsn := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.DatabaseOrDefault(),
	}

	if len(c.Username) > 0 {
		if len(c.Password) > 0 {
			dsn.User = url.UserPassword(c.Username, c.Password)
		} else {
			dsn.User = url.User(c.Username)
		}
	}

	queryArgs := url.Values{}
	if len(c.SSLMode) > 0 {
		queryArgs.Add("sslmode", c.SSLMode)
	}
	if c.ConnectTimeout > 0 {
		setTimeoutSeconds(queryArgs, "connect_timeout", c.ConnectTimeout)
	}
	if c.LockTimeout > 0 {
		setTimeoutMilliseconds(queryArgs, "lock_timeout", c.LockTimeout)
	}
	if c.StatementTimeout > 0 {
		setTimeoutMilliseconds(queryArgs, "statement_timeout", c.StatementTimeout)
	}
	if c.Schema != "" {
		queryArgs.Add("search_path", c.Schema)
	}
	if c.ApplicationName != "" {
		queryArgs.Add("application_name", c.ApplicationName)
	}

	dsn.RawQuery = queryArgs.Encode()
	return dsn.String()
}

// CreateLoggingDSN creates a postgres connection string from the config suitable for logging.
// It will not include the password.
func (c Config) CreateLoggingDSN() string {
	host := c.HostOrDefault()
	if c.PortOrDefault() != "" {
		host = host + ":" + c.PortOrDefault()
	}

	dsn := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.DatabaseOrDefault(),
	}

	if len(c.Username) > 0 {
		dsn.User = url.User(c.Username)
	}

	queryArgs := url.Values{}
	if len(c.SSLMode) > 0 {
		queryArgs.Add("sslmode", c.SSLMode)
	}
	if c.ConnectTimeout > 0 {
		setTimeoutSeconds(queryArgs, "connect_timeout", c.ConnectTimeout)
	}
	if c.LockTimeout > 0 {
		setTimeoutMilliseconds(queryArgs, "lock_timeout", c.LockTimeout)
	}
	if c.StatementTimeout > 0 {
		setTimeoutMilliseconds(queryArgs, "statement_timeout", c.StatementTimeout)
	}
	if c.Schema != "" {
		queryArgs.Add("search_path", c.Schema)
	}

	dsn.RawQuery = queryArgs.Encode()
	return dsn.String()
}

// Validate validates that user-provided values are valid, e.g. that timeouts
// can be exactly rounded into a multiple of a given base value.
func (c Config) Validate() error {
	if c.ConnectTimeout.Round(time.Second) != c.ConnectTimeout {
		return ex.New(ErrDurationConversion, ex.OptMessagef("connect_timeout=%s", c.ConnectTimeout))
	}
	if c.LockTimeout.Round(time.Millisecond) != c.LockTimeout {
		return ex.New(ErrDurationConversion, ex.OptMessagef("lock_timeout=%s", c.LockTimeout))
	}
	if c.StatementTimeout.Round(time.Millisecond) != c.StatementTimeout {
		return ex.New(ErrDurationConversion, ex.OptMessagef("statement_timeout=%s", c.StatementTimeout))
	}

	return nil
}

// ValidateProduction validates production configuration for the config.
func (c Config) ValidateProduction() error {
	if !(len(c.SSLMode) == 0 ||
		stringutil.EqualsCaseless(c.SSLMode, SSLModeRequire) ||
		stringutil.EqualsCaseless(c.SSLMode, SSLModeVerifyCA) ||
		stringutil.EqualsCaseless(c.SSLMode, SSLModeVerifyFull)) {
		return ex.New(ErrUnsafeSSLMode, ex.OptMessagef("sslmode: %s", c.SSLMode))
	}
	if len(c.Username) == 0 {
		return ex.New(ErrUsernameUnset)
	}
	if len(c.Password) == 0 {
		return ex.New(ErrPasswordUnset)
	}
	return c.Validate()
}

// setTimeoutMilliseconds sets a timeout value in connection string query parameters.
//
// Valid units for this parameter in PostgresSQL are "ms", "s", "min", "h"
// and "d" and the value should be between 0 and 2147483647ms. We explicitly
// cast to milliseconds but leave validation on the value to PostgreSQL.
//
//   blend=> BEGIN;
//   BEGIN
//   blend=> SET LOCAL lock_timeout TO '4000ms';
//   SET
//   blend=> SHOW lock_timeout;
//    lock_timeout
//   --------------
//    4s
//   (1 row)
//   --
//   blend=> SET LOCAL lock_timeout TO '4500ms';
//   SET
//   blend=> SHOW lock_timeout;
//    lock_timeout
//   --------------
//    4500ms
//   (1 row)
//   --
//   blend=> SET LOCAL lock_timeout = 'go';
//   ERROR:  invalid value for parameter "lock_timeout": "go"
//   blend=> SET LOCAL lock_timeout = '1ns';
//   ERROR:  invalid value for parameter "lock_timeout": "1ns"
//   HINT:  Valid units for this parameter are "ms", "s", "min", "h", and "d".
//   blend=> SET LOCAL lock_timeout = '-1ms';
//   ERROR:  -1 is outside the valid range for parameter "lock_timeout" (0 .. 2147483647)
//   --
//   blend=> COMMIT;
//   COMMIT
//
// See:
// - https://www.postgresql.org/docs/current/runtime-config-client.html#GUC-LOCK-TIMEOUT
// - https://www.postgresql.org/docs/current/runtime-config-client.html#GUC-STATEMENT-TIMEOUT
func setTimeoutMilliseconds(q url.Values, name string, d time.Duration) {
	ms := d.Round(time.Millisecond) / time.Millisecond
	q.Add(name, fmt.Sprintf("%dms", ms))
}

// setTimeoutSeconds sets a timeout value in connection string query parameters.
//
// This timeout is expected to be an exact number of seconds (as an integer)
// so we convert `d` to an integer first and set the value as a query parameter
// without units.
//
// See:
// - https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-CONNECT-CONNECT-TIMEOUT
func setTimeoutSeconds(q url.Values, name string, d time.Duration) {
	s := d.Round(time.Second) / time.Second
	q.Add(name, fmt.Sprintf("%d", s))
}
