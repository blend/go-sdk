/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

// FormatHeaders formats headers for output.
// Header keys will be printed in alphabetic order.
func FormatHeaders(tf logger.TextFormatter, keyColor ansi.Color, header http.Header) string {
	var keys []string
	for key := range header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s:%s", tf.Colorize(key, keyColor), header.Get(key)))
	}
	return "{ " + strings.Join(values, " ") + " }"
}

// ColorizeByStatusCode returns a value colored by an http status code.
func ColorizeByStatusCode(statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 {	//the http 2xx range is ok
		return ansi.ColorGreen.Apply(value)
	} else if statusCode == http.StatusInternalServerError {
		return ansi.ColorRed.Apply(value)
	}
	return ansi.ColorYellow.Apply(value)
}

// ColorizeByStatusCodeWithFormatter returns a value colored by an http status code with a given formatter.
func ColorizeByStatusCodeWithFormatter(tf logger.TextFormatter, statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 {	//the http 2xx range is ok
		return tf.Colorize(value, ansi.ColorGreen)
	} else if statusCode == http.StatusInternalServerError {
		return tf.Colorize(value, ansi.ColorRed)
	}
	return tf.Colorize(value, ansi.ColorYellow)
}

// ColorizeStatusCode colorizes a status code.
func ColorizeStatusCode(statusCode int) string {
	return ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode))
}

// ColorizeStatusCodeWithFormatter colorizes a status code with a given formatter.
func ColorizeStatusCodeWithFormatter(tf logger.TextFormatter, statusCode int) string {
	return ColorizeByStatusCodeWithFormatter(tf, statusCode, strconv.Itoa(statusCode))
}
