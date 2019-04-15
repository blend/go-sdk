package configutil

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

type config struct {
	Environment string `json:"env" yaml:"env" env:"SERVICE_ENV"`
	Other       string `json:"other" yaml:"other" env:"OTHER"`
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

func TestTryReadYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptPaths("testdata/config.yaml"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.yaml")
	assert.Equal("test_yaml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptPaths("testdata/config.yml"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.yml")
	assert.Equal("test_yml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptPaths("testdata/config.json"))
	assert.Nil(err)
	assert.Equal(path, "testdata/config.json")
	assert.Equal("test_json", cfg.Environment)
	assert.Equal("moo", cfg.Other)
}

func TestReadUnset(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	path, err := Read(&cfg, OptPaths(""))
	assert.Nil(err)
	assert.Empty(path)
	assert.NotEqual("dev", cfg.Environment)
}

func TestReadPathNotFound(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	_, err := Read(&cfg, OptPaths(filepath.Join("testdata", uuid.V4().String())))
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
