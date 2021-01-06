package httpstats

import "github.com/blend/go-sdk/webutil"

// HTTP stats constants
const (
	MetricNameHTTPRequest        string = string(webutil.FlagHTTPRequest)
	MetricNameHTTPRequestSize    string = MetricNameHTTPRequest + ".size"
	MetricNameHTTPRequestElapsed string = MetricNameHTTPRequest + ".elapsed"

	TagRoute  string = "route"
	TagMethod string = "method"
	TagStatus string = "status"

	RouteNotFound string = "not_found"
)
