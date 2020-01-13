package httpmetrics

import (
	"context"
	"strconv"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/webutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(webutil.HTTPRequest, stats.ListenerNameStats, webutil.NewHTTPRequestEventListener(func(_ context.Context, wre webutil.HTTPRequestEvent) {
		var route string
		if len(wre.Route) > 0 {
			route = stats.Tag(TagRoute, wre.Route)
		} else {
			route = stats.Tag(TagRoute, RouteNotFound)
		}

		method := stats.Tag(TagMethod, wre.Request.Method)
		tags := []string{
			route, method,
		}
		collector.Increment(MetricNameHTTPRequest, tags...)
		collector.Gauge(MetricNameHTTPRequestSize, float64(wre.Request.ContentLength), tags...)
	}))

	log.Listen(webutil.HTTPResponse, stats.ListenerNameStats, webutil.NewHTTPResponseEventListener(func(_ context.Context, wre webutil.HTTPResponseEvent) {
		var route string
		if len(wre.Route) > 0 {
			route = stats.Tag(TagRoute, wre.Route)
		} else {
			route = stats.Tag(TagRoute, RouteNotFound)
		}

		method := stats.Tag(TagMethod, wre.Request.Method)
		status := stats.Tag(TagStatus, strconv.Itoa(wre.StatusCode))
		tags := []string{
			route, method, status,
		}

		collector.Increment(MetricNameHTTPResponse, tags...)
		collector.Gauge(MetricNameHTTPResponseSize, float64(wre.ContentLength), tags...)
		collector.TimeInMilliseconds(MetricNameHTTPResponseElapsed, wre.Elapsed, tags...)
	}))
}
