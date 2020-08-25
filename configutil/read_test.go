package configutil

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

type config struct {
	Environment string `json:"env" yaml:"env" env:"SERVICE_ENV"`
	Other       string `json:"other" yaml:"other" env:"OTHER"`
}

type bareResolvedConfig struct {
	config
}

// Resolve implements configutil.BareResolver.
func (br *bareResolvedConfig) Resolve() error {
	br.Environment = "bare resolved"
	return nil
}

type resolvedConfig struct {
	config
}

// Resolve implements configutil.BareResolver.
func (r *resolvedConfig) Resolve(ctx context.Context) error {
	r.Environment = env.GetVars(ctx).String("ENVIRONMENT")
	return nil
}

func TestTryReadYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptFilePaths("testdata/config.yaml"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.yaml")
	assert.Equal("test_yaml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptFilePaths("testdata/config.yml"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.yml")
	assert.Equal("test_yml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptFilePaths("testdata/config.json"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.json")
	assert.Equal("test_json", cfg.Environment)
	assert.Equal("moo", cfg.Other)
}

func TestReadUnset(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptFilePaths(""))
	assert.Nil(err)
	assert.Empty(path)
	assert.NotEqual("dev", cfg.Environment)
}

func TestReadPathNotFound(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	_, err := Read(&cfg, OptFilePaths(filepath.Join("testdata", uuid.V4().String())))
	assert.True(IsNotExist(err))
}

func TestIsUnset(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsConfigPathUnset(ex.New(ErrConfigPathUnset)))
	assert.False(IsConfigPathUnset(ex.New(uuid.V4().String())))
}

func TestIsIgnored(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsIgnored(nil))
	assert.True(IsIgnored(ex.New(ErrConfigPathUnset)))
	assert.True(IsIgnored(ex.New(ErrInvalidConfigExtension)))
}

func TestDeserializeInvalid(t *testing.T) {
	assert := assert.New(t)

	err := deserialize(".???", nil, nil)
	assert.NotNil(err)
	assert.True(IsInvalidConfigExtension(err))
}

func TestDeserializeYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(deserialize(ExtensionYAML, bytes.NewBuffer([]byte("env: test\nother: foo")), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestDeserializeYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(deserialize(ExtensionYML, bytes.NewBuffer([]byte("env: test\nother: foo")), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestDeserializeAddsPrefix(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(deserialize("yml", bytes.NewBuffer([]byte("env: test\nother: foo")), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestDeserializeJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(deserialize(ExtensionJSON, bytes.NewBuffer([]byte(`{"env": "test", "other": "foo"}`)), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestReadBareResolve(t *testing.T) {
	assert := assert.New(t)

	var cfg bareResolvedConfig
	path, err := Read(&cfg, OptFilePaths(""))
	assert.Nil(err)
	assert.Empty(path)
	assert.Equal("bare resolved", cfg.Environment)
}

func TestReadResolver(t *testing.T) {
	assert := assert.New(t)

	var cfg resolvedConfig
	path, err := Read(&cfg,
		OptFilePaths(""),
		OptEnv(env.Vars{"ENVIRONMENT": "resolved"}),
	)
	assert.Nil(err)
	assert.Empty(path)
	assert.Equal("resolved", cfg.Environment)
}
