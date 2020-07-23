package httpstats

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(webutil.FlagHTTPRequest, stats.ListenerNameStats,
		webutil.NewHTTPRequestEventListener(func(_ context.Context, wre webutil.HTTPRequestEvent) {
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

			_ = collector.Increment(MetricNameHTTPRequest, tags...)
			_ = collector.Gauge(MetricNameHTTPRequestSize, float64(wre.ContentLength), tags...)
			_ = collector.Gauge(MetricNameHTTPRequestElapsed, timeutil.Milliseconds(wre.Elapsed), tags...)
			_ = collector.TimeInMilliseconds(MetricNameHTTPRequestElapsed, wre.Elapsed, tags...)
		}),
	)
}
