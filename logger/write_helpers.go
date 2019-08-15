package logger

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/webutil"
)

// WriteHTTPRequest is a helper method to write request start events to a writer.
func WriteHTTPRequest(tf TextFormatter, wr io.Writer, req *http.Request) {
	if ip := webutil.GetRemoteAddr(req); len(ip) > 0 {
		io.WriteString(wr, ip)
		io.WriteString(wr, Space)
	}
	io.WriteString(wr, tf.Colorize(req.Method, ansi.ColorBlue))
	if req.URL != nil {
		io.WriteString(wr, Space)
		io.WriteString(wr, req.URL.String())
	}
}

// WriteHTTPResponse is a helper method to write request complete events to a writer.
func WriteHTTPResponse(tf TextFormatter, wr io.Writer, req *http.Request, statusCode, contentLength int, contentType string, elapsed time.Duration) {
	if ip := webutil.GetRemoteAddr(req); len(ip) > 0 {
		io.WriteString(wr, ip)
		io.WriteString(wr, Space)
	}
	io.WriteString(wr, tf.Colorize(req.Method, ansi.ColorBlue))
	io.WriteString(wr, Space)
	io.WriteString(wr, req.URL.String())
	io.WriteString(wr, Space)
	io.WriteString(wr, ColorizeStatusCodeWithFormatter(tf, statusCode))
	io.WriteString(wr, Space)
	io.WriteString(wr, elapsed.String())
	if len(contentType) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, contentType)
	}
	io.WriteString(wr, Space)
	io.WriteString(wr, stringutil.FileSize(contentLength))
}

// FormatHeaders formats headers for output.
// Header keys will be printed in alphabetic order.
func FormatHeaders(tf TextFormatter, keyColor ansi.Color, header http.Header) string {
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

// FormatFields formats the output of fields.
// Field keys will be printed in alphabetic order.
func FormatFields(tf TextFormatter, keyColor ansi.Color, fields Fields) string {
	var keys []string
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var values []string
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%v", tf.Colorize(key, keyColor), fields[key]))
	}
	return strings.Join(values, " ")
}

// MergeDecomposed merges sets of decomposed data.
func MergeDecomposed(sets ...map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for _, set := range sets {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}
