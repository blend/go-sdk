package logger

import (
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
