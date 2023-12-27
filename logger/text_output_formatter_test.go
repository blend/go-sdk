/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestTextOutputFormatter(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()
	assert.False(tf.HideTimestamp)
	assert.False(tf.HideFields)
	assert.False(tf.NoColor)
	assert.Equal(DefaultTextTimeFormat, tf.TimeFormatOrDefault())

	tf = NewTextOutputFormatter(
		OptTextTimeFormat(time.RFC3339),
		OptTextHideTimestamp(),
		OptTextHideFields(),
		OptTextNoColor(),
	)

	assert.True(tf.HideTimestamp)
	assert.True(tf.HideFields)
	assert.True(tf.NoColor)
	assert.Equal(time.RFC3339, tf.TimeFormatOrDefault())

	tf = NewTextOutputFormatter(OptTextConfig(TextConfig{
		HideTimestamp: true,
		HideFields:    true,
		NoColor:       true,
		TimeFormat:    time.Kitchen,
	}))

	assert.True(tf.HideTimestamp)
	assert.True(tf.HideFields)
	assert.True(tf.NoColor)
	assert.Equal(time.Kitchen, tf.TimeFormatOrDefault())
}

func TestTextOutputFormatterWriteFormat(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()
	paths := []string{"foo", "bar"}
	ctxPaths := make([]string, len(paths))
	copy(ctxPaths, paths) // make sure that the original paths are not modified
	ctx := WithPath(context.Background(), ctxPaths...)

	message := "this is a test"
	buf := new(bytes.Buffer)

	assert.Nil(tf.WriteFormat(ctx, buf, NewMessageEvent(Info, message)))
	assert.Contains(buf.String(), message)
	assert.Equal(paths, GetPath(ctx)) // test that the context paths are not modified
}

func TestTextOutputFormatterColorize(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()
	assert.Equal(ansi.ColorRed.Apply("foo"), tf.Colorize("foo", ansi.ColorRed))
	tf.NoColor = true
	assert.Equal("foo", tf.Colorize("foo", ansi.ColorRed))
}

func TestTextOutputFormatterFormatFlag(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()
	assert.Equal("["+ansi.ColorRed.Apply("flag")+"]", tf.FormatFlag("flag", ansi.ColorRed))
}

func TestTextOutputFormatterFormatFlagNoColor(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter(OptTextNoColor())
	assert.Equal("[flag]", tf.FormatFlag("flag", ansi.ColorRed))
}

func TestTextOutputFormatterFormatTimestamp(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()

	actual := tf.FormatTimestamp(time.Date(2006, 01, 02, 03, 04, 05, 06, time.UTC))
	assert.Equal(ansi.ColorLightBlack.Apply("2006-01-02T03:04:05.000000006Z"), actual)

	tf.TimeFormat = time.Kitchen
	actual = tf.FormatTimestamp(time.Date(2006, 01, 02, 03, 04, 05, 06, time.UTC))
	assert.Equal(ansi.ColorLightBlack.Apply(fmt.Sprintf("%-30s", "3:04AM")), actual)
}

func TestTextOutputFormatterFormatTimestampNoColor(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter(OptTextNoColor())

	actual := tf.FormatTimestamp(time.Date(2006, 01, 02, 03, 04, 05, 06, time.UTC))
	assert.Equal("2006-01-02T03:04:05.000000006Z", actual)

	tf.TimeFormat = time.Kitchen
	actual = tf.FormatTimestamp(time.Date(2006, 01, 02, 03, 04, 05, 06, time.UTC))
	assert.Equal(fmt.Sprintf("%-30s", "3:04AM"), actual)
}

func TestTextOutputFormatterFormatFields(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter()
	actual := tf.FormatLabels(Labels{"foo": "bar", "buzz": "one"})

	expected := fmt.Sprintf("%s=%v %s=%v", ansi.ColorBlue.Apply("buzz"), "one", ansi.ColorBlue.Apply("foo"), "bar")
	assert.Equal(expected, actual)
}
