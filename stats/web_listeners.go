package stats

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// AddWebListeners adds web listeners.
func AddWebListeners(log logger.Listenable, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(webutil.HTTPResponse, ListenerNameStats, webutil.NewHTTPResponseEventListener(func(_ context.Context, wre webutil.HTTPResponseEvent) {
		var route string
		if len(wre.Route) > 0 {
			route = Tag(TagRoute, wre.Route)
		} else {
			route = Tag(TagRoute, RouteNotFound)
		}

		method := Tag(TagMethod, wre.Request.Method)
		status := Tag(TagStatus, strconv.Itoa(wre.StatusCode))
		tags := []string{
			route, method, status,
		}

		stats.Increment(MetricNameHTTPRequest, tags...)
		stats.TimeInMilliseconds(MetricNameHTTPRequestElapsed, wre.Elapsed, tags...)
	}))
}
