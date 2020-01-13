package db

import (
	"database/sql"
	"github.com/blend/go-sdk/logger"
)

// Option is an option for database connections.
type Option func(c *Connection) error

// OptConnection sets the underlying driver connection.
func OptConnection(conn *sql.DB) Option {
	return func(c *Connection) error {
		c.Connection = conn
		return nil
	}
}

// OptLog sets the tracer on the connection.
func OptLog(log logger.Log) Option {
	return func(c *Connection) error {
		c.Log = log
		return nil
	}
}

// OptTracer sets the tracer on the connection.
func OptTracer(tracer Tracer) Option {
	return func(c *Connection) error {
		c.Tracer = tracer
		return nil
	}
}

// OptStatementInterceptor sets the statement interceptor on the connection.
func OptStatementInterceptor(interceptor StatementInterceptor) Option {
	return func(c *Connection) error {
		c.StatementInterceptor = interceptor
		return nil
	}
}

// OptConfig sets the config on a connection.
func OptConfig(cfg Config) Option {
	return func(c *Connection) error {
		c.Config = cfg
		return nil
	}
}

// OptConfigFromEnv sets the config on a connection from the environment.
func OptConfigFromEnv() Option {
	return func(c *Connection) error {
		cfg, err := NewConfigFromEnv()
		if err != nil {
			return err
		}
		c.Config = cfg
		return nil
	}
}
