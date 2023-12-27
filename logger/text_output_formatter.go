/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/bufferutil"
)

var (
	_ WriteFormatter = (*TextOutputFormatter)(nil)
)

// NewTextOutputFormatter returns a new text writer for a given output.
func NewTextOutputFormatter(options ...TextOutputFormatterOption) *TextOutputFormatter {
	tf := &TextOutputFormatter{
		BufferPool:	bufferutil.NewPool(DefaultBufferPoolSize),
		TimeFormat:	DefaultTextTimeFormat,
	}

	for _, option := range options {
		option(tf)
	}

	return tf
}

// TextOutputFormatterOption is an option for text formatters.
type TextOutputFormatterOption func(*TextOutputFormatter)

// OptTextConfig sets the text formatter config.
func OptTextConfig(cfg TextConfig) TextOutputFormatterOption {
	return func(tf *TextOutputFormatter) {
		tf.HideTimestamp = cfg.HideTimestamp
		tf.HideFields = cfg.HideFields
		tf.NoColor = cfg.NoColor
		tf.TimeFormat = cfg.TimeFormatOrDefault()
	}
}

// OptTextTimeFormat sets the timestamp format.
func OptTextTimeFormat(format string) TextOutputFormatterOption {
	return func(tf *TextOutputFormatter) { tf.TimeFormat = format }
}

// OptTextHideTimestamp hides the timestamp in output.
func OptTextHideTimestamp() TextOutputFormatterOption {
	return func(tf *TextOutputFormatter) { tf.HideTimestamp = true }
}

// OptTextHideFields hides the fields in output.
func OptTextHideFields() TextOutputFormatterOption {
	return func(tf *TextOutputFormatter) { tf.HideFields = true }
}

// OptTextNoColor disables colorizing text output.
func OptTextNoColor() TextOutputFormatterOption {
	return func(tf *TextOutputFormatter) { tf.NoColor = true }
}

// TextOutputFormatter handles formatting messages as text.
type TextOutputFormatter struct {
	HideTimestamp	bool
	HideFields	bool
	NoColor		bool
	TimeFormat	string

	BufferPool	*bufferutil.Pool
}

// TimeFormatOrDefault returns the time format or a default
func (tf TextOutputFormatter) TimeFormatOrDefault() string {
	if len(tf.TimeFormat) > 0 {
		return tf.TimeFormat
	}
	return DefaultTextTimeFormat
}

// Colorize (optionally) applies a color to a string.
func (tf TextOutputFormatter) Colorize(value string, color ansi.Color) string {
	if tf.NoColor {
		return value
	}
	return color.Apply(value)
}

// FormatFlag formats the flag portion of the message.
func (tf TextOutputFormatter) FormatFlag(flag string, color ansi.Color) string {
	return fmt.Sprintf("[%s]", tf.Colorize(string(flag), color))
}

// FormatTimestamp returns a new timestamp string.
func (tf TextOutputFormatter) FormatTimestamp(ts time.Time) string {
	value := ts.Format(tf.TimeFormatOrDefault())
	return tf.Colorize(fmt.Sprintf("%-30s", value), ansi.ColorLightBlack)
}

// FormatPath returns the sub-context path section of the message as a string.
func (tf TextOutputFormatter) FormatPath(path ...string) string {
	if len(path) == 0 {
		return ""
	}
	if len(path) == 1 {
		return fmt.Sprintf("[%s]", tf.Colorize(path[0], ansi.ColorBlue))
	}
	if !tf.NoColor {
		var colorized []string
		for _, subPath := range path {
			colorized = append(colorized, tf.Colorize(subPath, ansi.ColorBlue))
		}
		path = colorized
	}
	return fmt.Sprintf("[%s]", strings.Join(path, " > "))
}

// FormatAnnotations returns the scope annotations section of the message as a string.
func (tf TextOutputFormatter) FormatAnnotations(annotations Annotations) string {
	return FormatAnnotations(tf, ansi.ColorGreen, annotations)
}

// FormatLabels returns the scope labels section of the message as a string.
func (tf TextOutputFormatter) FormatLabels(labels Labels) string {
	return FormatLabels(tf, ansi.ColorBlue, labels)
}

// WriteFormat implements write formatter.
func (tf TextOutputFormatter) WriteFormat(ctx context.Context, output io.Writer, e Event) error {
	buffer := tf.BufferPool.Get()
	defer tf.BufferPool.Put(buffer)

	if !tf.HideTimestamp {
		buffer.WriteString(tf.FormatTimestamp(GetEventTimestamp(ctx, e)))
		buffer.WriteString(Space)
	}

	scopePath := GetPath(ctx)
	if scopePath != nil {
		buffer.WriteString(tf.FormatPath(scopePath...))
		buffer.WriteString(Space)
	}

	buffer.WriteString(tf.FormatFlag(e.GetFlag(), FlagTextColor(e.GetFlag())))
	buffer.WriteString(Space)

	if typed, ok := e.(TextWritable); ok {
		typed.WriteText(tf, buffer)
	} else if stringer, ok := e.(fmt.Stringer); ok {
		buffer.WriteString(stringer.String())
	}

	if !tf.HideFields {
		labels := GetLabels(ctx)
		if len(labels) > 0 {
			buffer.WriteString("\t")
			buffer.WriteString(tf.FormatLabels(labels))
		}

		annotations := GetAnnotations(ctx)
		if len(annotations) > 0 {
			buffer.WriteString("\t")
			buffer.WriteString(tf.FormatAnnotations(annotations))
		}
	}

	buffer.WriteString(Newline)
	_, err := io.Copy(output, buffer)
	return err
}
