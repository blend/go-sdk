package httpmetrics

import "github.com/blend/go-sdk/webutil"

// HTTP stats constants
const (
	MetricNameHTTPRequest         string = string(webutil.HTTPRequest)
	MetricNameHTTPRequestSize     string = MetricNameHTTPRequest + ".size"
	MetricNameHTTPResponse        string = string(webutil.HTTPResponse)
	MetricNameHTTPResponseSize    string = MetricNameHTTPResponse + ".size"
	MetricNameHTTPResponseElapsed string = MetricNameHTTPResponse + ".elapsed"

	TagRoute  string = "route"
	TagMethod string = "method"
	TagStatus string = "status"

	RouteNotFound string = "not_found"
)
