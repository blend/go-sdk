package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestConfigCreateDSN(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{
		DSN:      "foo",
		Host:     "bar",
		Username: "bailey",
		Password: "dog",
		Database: "blend",
		Schema:   "mortgages",
		SSLMode:  SSLModeVerifyCA,
	}

	assert.Equal("foo", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "bailey",
		Password: "dog",
		Database: "blend",
		Schema:   "mortgages",
		SSLMode:  SSLModeVerifyCA,
	}

	assert.Equal("postgres://bailey:dog@bar:1234/blend?search_path=mortgages&sslmode=verify-ca", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "bailey",
		Password: "dog",
		Database: "blend",
		Schema:   "mortgages",
	}

	assert.Equal("postgres://bailey:dog@bar:1234/blend?search_path=mortgages", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "bailey",
		Password: "dog",
		Database: "blend",
	}

	assert.Equal("postgres://bailey:dog@bar:1234/blend", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "bailey",
		Password: "dog",
	}

	assert.Equal("postgres://bailey:dog@bar:1234/postgres", cfg.CreateDSN())

	cfg = &Config{
		Host:     "bar",
		Port:     "1234",
		Username: "bailey",
	}

	assert.Equal("postgres://bailey@bar:1234/postgres", cfg.CreateDSN())

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

func TestConfigValidateProduction(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsUsernameUnset(Config{}.ValidateProduction()))
	assert.True(IsPasswordUnset(Config{Username: "foo"}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModeDisable}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModeAllow}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: SSLModePrefer}.ValidateProduction()))
	assert.True(IsUnsafeSSLMode(Config{Username: "foo", Password: "bar", SSLMode: "NOT A REAL MODE"}.ValidateProduction()))
}

func TestConfigResolve(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{
		Host:     "bar",
		Username: "bailey",
		Password: "dog",
		Database: "blend",
		Schema:   "mortgages",
		SSLMode:  SSLModeVerifyCA,
	}

	resolved, err := cfg.Resolve()
	assert.Nil(err)
	assert.NotNil(resolved)
	assert.Equal("bar", resolved.Host)
}
