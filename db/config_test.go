/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
)

var (
	_ configutil.Resolver = (*Config)(nil)
)

func TestConfigCreateDSN(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{
		Host:             "bar",
		Port:             "1234",
		Username:         "example-string",
		Password:         "dog",
		Database:         "blend",
		Schema:           "primary_schema",
		ApplicationName:  "this-pod-7897744df9-v4bbx",
		SSLMode:          SSLModeVerifyCA,
		LockTimeout:      1704 * time.Millisecond,
		StatementTimeout: 2704 * time.Millisecond,
		ConnectTimeout:   7 * time.Second,
	}

	assert.Equal("postgres://example-string:dog@bar:1234/blend?application_name=this-pod-7897744df9-v4bbx&connect_timeout=7&lock_timeout=1704ms&search_path=primary_schema&sslmode=verify-ca&statement_timeout=2704ms", cfg.CreateDSN())

	cfg = &Config{
		DSN:             "foo",
		Host:            "bar",
		Username:        "example-string",
		Password:        "dog",
		Database:        "blend",
		Schema:          "primary_schema",
		ApplicationName: "this-pod-7897744df9-v4bbx",
		SSLMode:         SSLModeVerifyCA,
	}

	assert.Equal("foo", cfg.CreateDSN())

	cfg = &Config{
		Host:            "bar",
		Port:            "1234",
		Username:        "example-string",
		Password:        "dog",
		Database:        "blend",
		Schema:          "primary_schema",
		ApplicationName: "this-pod-7897744df9-v4bbx",
		SSLMode:         SSLModeVerifyCA,
	}

	assert.Equal("postgres://example-string:dog@bar:1234/blend?application_name=this-pod-7897744df9-v4bbx&search_path=primary_schema&sslmode=verify-ca", cfg.CreateDSN())

	cfg = &Config{
		Host:            "bar",
		Port:            "1234",
		Username:        "example-string",
		Password:        "dog",
		Database:        "blend",
		Schema:          "primary_schema",
		ApplicationName: "this-pod-7897744df9-v4bbx",
	}

	assert.Equal("postgres://example-string:dog@bar:1234/blend?application_name=this-pod-7897744df9-v4bbx&search_path=primary_schema", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "example-string",
		Password: "dog",
		Database: "blend",
	}

	assert.Equal("postgres://example-string:dog@bar:1234/blend", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "example-string",
		Password: "dog",
	}

	assert.Equal("postgres://example-string:dog@bar:1234/postgres", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "example-string",
	}

	assert.Equal("postgres://example-string@bar:1234/postgres", cfg.CreateDSN())

	cfg = &Config{
		Host: "bar",
		Port: "1234",
	}

	assert.Equal("postgres://bar:1234/postgres", cfg.CreateDSN())

	cfg = &Config{
		Host: "bar",
	}

	assert.Equal("postgres://bar:5432/postgres", cfg.CreateDSN())

	cfg = &Config{}
	assert.Equal("postgres://localhost:5432/postgres", cfg.CreateDSN())
}

func TestNewConfigFromDSN(t *testing.T) {
	assert := assert.New(t)

	// Fails to parse URL
	dsn := "a://b"
	parsed, err := NewConfigFromDSN(dsn)
	assert.Equal(Config{}, parsed)
	assert.Equal("invalid connection protocol: a", fmt.Sprintf("%v", err))

	// Success, coverage for lots of fields
	dsn = "postgres://example-string:dog@bar:1234/blend?connect_timeout=5&lock_timeout=4500ms&statement_timeout=5500ms&sslmode=verify-ca"
	parsed, err = NewConfigFromDSN(dsn)
	assert.Nil(err)
	expected := Config{
		Host:             "bar",
		Port:             "1234",
		Database:         "blend",
		Username:         "example-string",
		Password:         "dog",
		ConnectTimeout:   5 * time.Second,
		LockTimeout:      4500 * time.Millisecond,
		StatementTimeout: 5500 * time.Millisecond,
		SSLMode:          SSLModeVerifyCA,
	}
	assert.Equal(expected, parsed)

	// Failure to parse lock timeout
	dsn = "postgres://bar:1234/blend?lock_timeout=1000"
	parsed, err = NewConfigFromDSN(dsn)
	partial := Config{
		Host:     "bar",
		Database: "blend",
	}
	assert.Equal(partial, parsed)
	assert.Matches(`(?m)time: missing unit in duration (")?1000(")?; field: lock_timeout$`, fmt.Sprint(err))

	// Failure to parse statement timeout
	dsn = "postgres://bar:1234/blend?statement_timeout=2000"
	parsed, err = NewConfigFromDSN(dsn)
	partial = Config{
		Host:     "bar",
		Port:     "1234",
		Database: "blend",
	}
	assert.Equal(partial, parsed)
	assert.Matches(`(?m)time: missing unit in duration (")?2000(")?; field: statement_timeout$`, fmt.Sprint(err))
}

func TestNewConfigFromDSNWithSchema(t *testing.T) {
	assert := assert.New(t)

	dsn := "postgres://example-string:dog@bar:1234/blend?connect_timeout=5&sslmode=verify-ca&search_path=primary_schema&application_name=this-pod-7897744df9-v4bbx"

	parsed, err := NewConfigFromDSN(dsn)
	assert.Nil(err)

	expected := Config{
		Host:            "bar",
		Port:            "1234",
		Database:        "blend",
		Schema:          "primary_schema",
		ApplicationName: "this-pod-7897744df9-v4bbx",
		Username:        "example-string",
		Password:        "dog",
		ConnectTimeout:  5 * time.Second,
		SSLMode:         SSLModeVerifyCA,
	}
	assert.Equal(expected, parsed)
}

func TestNewConfigFromDSNConnectTimeoutParseError(t *testing.T) {
	assert := assert.New(t)

	dsn := "postgres://example-string:dog@bar:1234/blend?connect_timeout=abcd&sslmode=verify-ca"

	_, err := NewConfigFromDSN(dsn)
	assert.NotNil(err)
}

func TestConfigValidateProduction(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsUsernameUnset(Config{}.ValidateProduction()))
	assert.True(IsPasswordUnset(Config{Username: "foo"}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModeDisable}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModeAllow}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModePrefer}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: "NOT A REAL MODE"}.ValidateProduction()))
	assert.Nil(Config{Username: "foo", Password: "bar", SSLMode: SSLModeVerifyFull}.ValidateProduction())
	assert.True(IsDurationConversion(Config{Username: "foo", Password: "bar", SSLMode: SSLModeVerifyFull, LockTimeout: time.Nanosecond}.ValidateProduction()))
	assert.True(IsDurationConversion(Config{Username: "foo", Password: "bar", SSLMode: SSLModeVerifyFull, StatementTimeout: time.Nanosecond}.ValidateProduction()))
}

func TestConfigReparse(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{
		Host:            "bar",
		Username:        "example-string",
		Password:        "dog",
		Database:        "blend",
		Schema:          "primary_schema",
		SSLMode:         SSLModeVerifyCA,
		BufferPoolSize:  10,
		MaxConnections:  10,
		IdleConnections: 5,
		MaxLifetime:     100,
	}

	resolved, err := cfg.Reparse()
	assert.Nil(err)
	assert.NotNil(resolved)
	assert.Equal("bar", resolved.Host)

	assert.Equal(resolved.BufferPoolSize, 10)
	assert.Equal(resolved.MaxConnections, 10)
	assert.Equal(resolved.IdleConnections, 5)
	assert.Equal(resolved.MaxLifetime, 100)
}
