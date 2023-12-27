/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestTryReadYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	paths, err := Read(&cfg, OptPaths("testdata/config.yaml"))
	assert.Nil(err)
	assert.Len(paths, 1)
	assert.Equal("testdata/config.yaml", paths[0])
	assert.Equal("test_yaml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	paths, err := Read(&cfg, OptPaths("testdata/config.yml"))
	assert.Nil(err)
	assert.Len(paths, 1)
	assert.Equal("testdata/config.yml", paths[0])
	assert.Equal("test_yml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestTryReadJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	paths, err := Read(&cfg, OptPaths("testdata/config.json"))
	assert.Nil(err)
	assert.Len(paths, 1)
	assert.Equal("testdata/config.json", paths[0])
	assert.Equal("test_json", cfg.Environment)
	assert.Equal("moo", cfg.Other)
}

func TestReadUnset(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	paths, err := Read(&cfg, OptPaths(""))
	assert.Nil(err)
	assert.Empty(paths)
	assert.Empty(paths)
	assert.NotEqual("dev", cfg.Environment)
}

func TestReadMany(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	paths, err := Read(&cfg, OptPaths("testdata/project.yml", "testdata/config.yml"))
	assert.Nil(err)
	assert.Equal([]string{"testdata/project.yml", "testdata/config.yml"}, paths)
	assert.Equal("test_yml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
	assert.Equal("project-base", cfg.Base)
}

func TestReadPathNotFound(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	_, err := Read(&cfg, OptPaths(filepath.Join("testdata", uuid.V4().String())))
	assert.Nil(err)
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

func TestReadResolver(t *testing.T) {
	assert := assert.New(t)

	var cfg resolvedConfig
	path, err := Read(&cfg,
		OptPaths(""),
		OptEnv(env.Vars{"ENVIRONMENT": "resolved"}),
	)
	assert.Nil(err)
	assert.Empty(path)
	assert.Equal("resolved", cfg.Environment)
}

func TestRead_multiple(t *testing.T) {
	assert := assert.New(t)

	contents0 := `
serviceName: "serviceName-contents0"

field0: "field0-contents0"
field2: "field2-contents0"
field3: "field3-contents0"

child:
  field0: "child-field0-contents0"
  field2: "child-field2-contents0"
  field3: "child-field3-contents0"
`

	contents1 := `
serviceEnv: "serviceEnv-contents1"

field1: "field1-contents1"
field3: "field3-contents1"

child:
  field1: "child-field1-contents1"
  field3: "child-field3-contents1"
`

	contents2 := `
version: "version-contents2"

field2: "field2-contents2"

child:
  field2: "child-field2-contents2"
`

	var cfg fullConfig
	path, err := Read(&cfg,
		OptPaths(""),
		OptAddContentString("yml", contents0),
		OptAddContentString("yml", contents1),
		OptAddContentString("yml", contents2),
		OptEnv(env.Vars{"SERVICE_ENV": "env-resolved"}),
	)
	assert.Nil(err)
	assert.Empty(path)
	assert.Equal("env-resolved", cfg.ServiceEnv)

	assert.Equal("serviceName-contents0", cfg.ServiceName)
	assert.Equal("version-contents2", cfg.Version)

	assert.Equal("field0-contents0", cfg.Field0)
	assert.Equal("field1-contents1", cfg.Field1)
	assert.Equal("field2-contents2", cfg.Field2)
	assert.Equal("field3-contents1", cfg.Field3)

	assert.Equal("child-field0-contents0", cfg.Child.Field0)
	assert.Equal("child-field1-contents1", cfg.Child.Field1)
	assert.Equal("child-field2-contents2", cfg.Child.Field2)
	assert.Equal("child-field3-contents1", cfg.Child.Field3)
}
