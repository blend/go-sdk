package httpmetrics

import "github.com/blend/go-sdk/webutil"

// HTTP stats constants
const (
	MetricNameHTTPRequest         string = string(webutil.HTTPRequest)
	MetricNameHTTPRequestSize     string = string(webutil.HTTPRequest) + ".size"
	MetricNameHTTPResponse        string = string(webutil.HTTPResponse)
	MetricNameHTTPResponseSize    string = string(webutil.HTTPResponse) + ".size"
	MetricNameHTTPResponseElapsed string = MetricNameHTTPResponse + ".elapsed"

	TagRoute  string = "route"
	TagMethod string = "method"
	TagStatus string = "status"

	RouteNotFound string = "not_found"
)
