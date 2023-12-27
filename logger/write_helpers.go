/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blend/go-sdk/ansi"
)

// FormatAnnotations formats the output of annotations as a string.
// Field keys will be printed in alphabetic order.
func FormatAnnotations(tf TextFormatter, keyColor ansi.Color, annotations Annotations) string {
	var keys []string
	for key := range annotations {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%s", tf.Colorize(key, keyColor), annotations[key]))
	}
	return strings.Join(values, " ")
}

// FormatLabels formats the output of labels as a string.
// Field keys will be printed in alphabetic order.
func FormatLabels(tf TextFormatter, keyColor ansi.Color, labels Labels) string {
	var keys []string
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%s", tf.Colorize(key, keyColor), labels[key]))
	}
	return strings.Join(values, " ")
}
