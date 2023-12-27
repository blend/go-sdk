/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestOptConfig(t *testing.T) {
	assert := assert.New(t)

	log := None()
	assert.Nil(OptConfig(Config{
		Flags:		[]string{"foo", "bar"},
		Writable:	[]string{"foo"},
		Format:		FormatJSON,
	})(log))

	assert.Any(log.Flags.Flags(), func(v interface{}) bool { return v.(string) == "foo" })
	assert.Any(log.Flags.Flags(), func(v interface{}) bool { return v.(string) == "bar" })
	assert.Any(log.Writable.Flags(), func(v interface{}) bool { return v.(string) == "foo" })
	assert.None(log.Writable.Flags(), func(v interface{}) bool { return v.(string) == "bar" })
}

func TestOptConfigFromEnv(t *testing.T) {
	assert := assert.New(t)

	defer env.Restore()
	env.Env().Set("LOG_FLAGS", "foo,bar")
	env.Env().Set("LOG_WRITABLE", "foo")
	env.Env().Set("LOG_HIDE_TIMESTAMP", "true")
	env.Env().Set("LOG_HIDE_FIELDS", "true")
	env.Env().Set("LOG_TIME_FORMAT", time.Kitchen)
	env.Env().Set("NO_COLOR", "true")

	log := None()
	assert.Nil(OptConfigFromEnv()(log))

	assert.Any(log.Flags.Flags(), func(v interface{}) bool { return v.(string) == "foo" })
	assert.Any(log.Flags.Flags(), func(v interface{}) bool { return v.(string) == "bar" })
	assert.Any(log.Writable.Flags(), func(v interface{}) bool { return v.(string) == "foo" })
	assert.None(log.Writable.Flags(), func(v interface{}) bool { return v.(string) == "bar" })
	assert.True(log.Formatter.(*TextOutputFormatter).HideTimestamp)
	assert.True(log.Formatter.(*TextOutputFormatter).HideFields)
	assert.True(log.Formatter.(*TextOutputFormatter).NoColor)
	assert.Equal(time.Kitchen, log.Formatter.(*TextOutputFormatter).TimeFormat)
}

func TestOptOutput(t *testing.T) {
	assert := assert.New(t)

	log := None()

	buf := new(bytes.Buffer)
	assert.Nil(OptOutput(buf)(log))

	typed, ok := log.Output.(*InterlockedWriter)
	assert.True(ok)
	assert.NotNil(typed.Output)
}

func TestOptions(t *testing.T) {
	assert := assert.New(t)

	log := None()

	assert.Nil(log.Output)
	assert.Nil(OptOutput(new(bytes.Buffer))(log))
	assert.NotNil(log.Output)

	assert.Nil(log.Formatter)
	assert.Nil(OptText(OptTextNoColor())(log))
	assert.NotNil(log.Formatter)
	assert.True(log.Formatter.(*TextOutputFormatter).NoColor)

	assert.Nil(OptJSON(OptJSONPretty())(log))
	assert.NotNil(log.Formatter)
	assert.True(log.Formatter.(*JSONOutputFormatter).Pretty)

	assert.True(log.Flags.None())
	assert.Nil(OptFlags(NewFlags("test1", "test2"))(log))
	assert.False(log.Flags.None())
	assert.True(log.Flags.IsEnabled("test1"))
	assert.True(log.Flags.IsEnabled("test2"))

	assert.False(log.Flags.IsEnabled("foo"))
	assert.Nil(OptEnabled("foo")(log))
	assert.True(log.Flags.IsEnabled("foo"))
	assert.Nil(OptDisabled("foo")(log))
	assert.False(log.Flags.IsEnabled("foo"))

	assert.False(log.Flags.All())
	assert.Nil(OptAll()(log))
	assert.True(log.Flags.All())
}
