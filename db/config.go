package db

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
)

// NewConfig creates a new config.
func NewConfig() *Config {
	return &Config{}
}

// NewConfigFromDSN creates a new config from a dsn.
func NewConfigFromDSN(dsn string) (*Config, error) {
	parsed, err := ParseURL(dsn)
	if err != nil {
		return nil, exception.New(err)
	}

	var config Config
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
		}
	}
	return &config, nil
}

// NewConfigFromEnv returns a new config from the environment.
// The environment variable mappings are as follows:
//	-	DATABSE_URL 	= DSN 	//note that this has precedence over other vars (!!)
// 	-	DB_HOST 		= Host
//	-	DB_PORT 		= Port
//	- 	DB_NAME 		= Database
//	-	DB_SCHEMA		= Schema
//	-	DB_USER 		= Username
//	-	DB_PASSWORD 	= Password
//	-	DB_SSLMODE 		= SSLMode
func NewConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.Env().ReadInto(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// MustNewConfigFromEnv returns a new config from the environment,
// it will panic if there is an error.
func MustNewConfigFromEnv() *Config {
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
	// Schema is the application schema within the database, defaults to `public`.
	Schema string `json:"schema,omitempty" yaml:"schema,omitempty" env:"DB_SCHEMA"`
	// Username is the username for the connection via password auth.
	Username string `json:"username,omitempty" yaml:"username,omitempty" env:"DB_USER"`
	// Password is the password for the connection via password auth.
	Password string `json:"password,omitempty" yaml:"password,omitempty" env:"DB_PASSWORD"`
	// SSLMode is the sslmode for the connection.
	SSLMode string `json:"sslMode,omitempty" yaml:"sslMode,omitempty" env:"DB_SSLMODE"`
	// PlanCacheDisabled indicates if we should use the prepared statement plan cache.
	PlanCacheDisabled *bool `json:"planCacheDisabled,omitempty" yaml:"planCacheDisabled,omitempty" env:"DB_DISABLE_PLAN_CACHE"`
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

// WithEngine sets the databse engine.
func (c *Config) WithEngine(engine string) *Config {
	c.Engine = engine
	return c
}

// WithDSN sets the config dsn and returns a reference to the config.
func (c *Config) WithDSN(dsn string) *Config {
	c.DSN = dsn
	return c
}

// WithHost sets the config host and returns a reference to the config.
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

// WithPort sets the config host and returns a reference to the config.
func (c *Config) WithPort(port string) *Config {
	c.Port = port
	return c
}

// WithDatabase sets the config database and returns a reference to the config.
func (c *Config) WithDatabase(database string) *Config {
	c.Database = database
	return c
}

// WithSchema sets the config schema and returns a reference to the config.
func (c *Config) WithSchema(schema string) *Config {
	c.Schema = schema
	return c
}

// WithUsername sets the config username and returns a reference to the config.
func (c *Config) WithUsername(username string) *Config {
	c.Username = username
	return c
}

// WithPassword sets the config password and returns a reference to the config.
func (c *Config) WithPassword(password string) *Config {
	c.Password = password
	return c
}

// WithSSLMode sets the config sslMode and returns a reference to the config.
func (c *Config) WithSSLMode(sslMode string) *Config {
	c.SSLMode = sslMode
	return c
}

// GetEngine returns the database engine.
func (c Config) GetEngine(inherited ...string) string {
	return configutil.CoalesceString(c.Engine, DefaultEngine, inherited...)
}

// GetDSN returns the postgres dsn (fully quallified url) for the config.
// If unset, it's generated from the host, port and database.
func (c Config) GetDSN(inherited ...string) string {
	return configutil.CoalesceString(c.DSN, "", inherited...)
}

// GetHost returns the postgres host for the connection or a default.
func (c Config) GetHost(inherited ...string) string {
	return configutil.CoalesceString(c.Host, DefaultHost, inherited...)
}

// GetPort returns the port for a connection if it is not the standard postgres port.
func (c Config) GetPort(inherited ...string) string {
	return configutil.CoalesceString(c.Port, DefaultPort, inherited...)
}

// GetDatabase returns the connection database or a default.
func (c Config) GetDatabase(inherited ...string) string {
	return configutil.CoalesceString(c.Database, DefaultDatabase, inherited...)
}

// GetSchema returns the connection schema or a default.
func (c Config) GetSchema(inherited ...string) string {
	return configutil.CoalesceString(c.Schema, "", inherited...)
}

// GetUsername returns the connection username or a default.
func (c Config) GetUsername(inherited ...string) string {
	return configutil.CoalesceString(c.Username, "", inherited...)
}

// GetPassword returns the connection password or a default.
func (c Config) GetPassword(inherited ...string) string {
	return configutil.CoalesceString(c.Password, "", inherited...)
}

// GetSSLMode returns the connection ssl mode.
// It defaults to unset, which will then use the driver defaults.
func (c Config) GetSSLMode(inherited ...string) string {
	return configutil.CoalesceString(c.SSLMode, "", inherited...)
}

// GetPlanCacheDisabled returns if we should disable the statement plan cache or a default.
func (c Config) GetPlanCacheDisabled(inherited ...bool) bool {
	return configutil.CoalesceBool(c.PlanCacheDisabled, DefaultPlanCacheDisabled, inherited...)
}

// GetIdleConnections returns the number of idle connections or a default.
func (c Config) GetIdleConnections(inherited ...int) int {
	return configutil.CoalesceInt(c.IdleConnections, DefaultIdleConnections, inherited...)
}

// GetMaxConnections returns the maximum number of connections or a default.
func (c Config) GetMaxConnections(inherited ...int) int {
	return configutil.CoalesceInt(c.MaxConnections, DefaultMaxConnections, inherited...)
}

// GetMaxLifetime returns the maximum lifetime of a driver connection.
func (c Config) GetMaxLifetime(inherited ...time.Duration) time.Duration {
	return configutil.CoalesceDuration(c.MaxLifetime, DefaultMaxLifetime, inherited...)
}

// GetBufferPoolSize returns the number of query buffers to maintain or a default.
func (c Config) GetBufferPoolSize(inherited ...int) int {
	return configutil.CoalesceInt(c.BufferPoolSize, DefaultBufferPoolSize, inherited...)
}

// CreateDSN creates a postgres connection string from the config.
func (c Config) CreateDSN() string {
	if len(c.GetDSN()) > 0 {
		return c.GetDSN()
	}

	var host string
	if len(c.GetPort()) > 0 {
		host = fmt.Sprintf("%s:%s", c.GetHost(), c.GetPort())
	} else {
		host = c.GetHost()
	}

	dsn := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   c.GetDatabase(),
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
	if len(c.Schema) > 0 {
		queryArgs.Add("search_path", c.Schema)
	}

	dsn.RawQuery = queryArgs.Encode()
	return dsn.String()
}

// Resolve creates a DSN and reparses it, in case some values need to be coalesced.
func (c Config) Resolve() (*Config, error) {
	return NewConfigFromDSN(c.CreateDSN())
}

// MustResolve creates a DSN and reparses it, in case some values need to be coalesced,
// and panics if there is an error.
func (c Config) MustResolve() *Config {
	cfg, err := NewConfigFromDSN(c.CreateDSN())
	if err != nil {
		panic(err)
	}
	return cfg
}

// ValidateProduction validates production configuration for the config.
func (c Config) ValidateProduction() error {
	if !(len(c.GetSSLMode()) == 0 ||
		stringutil.EqualsCaseless(c.GetSSLMode(), SSLModeRequire) ||
		stringutil.EqualsCaseless(c.GetSSLMode(), SSLModeVerifyCA) ||
		stringutil.EqualsCaseless(c.GetSSLMode(), SSLModeVerifyFull)) {
		return exception.New(ErrUnsafeSSLMode).WithMessagef("sslmode: %s", c.GetSSLMode())
	}

	if len(c.GetUsername()) == 0 {
		return exception.New(ErrUsernameUnset)
	}
	if len(c.GetPassword()) == 0 {
		return exception.New(ErrPasswordUnset)
	}
	return nil
}
