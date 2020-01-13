package configutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

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
