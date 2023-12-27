/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package httpstats

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sanitize"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.FilterListenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	options := stats.NewAddListenerOptions(opts...)

	requestSanitizer := sanitize.NewRequestSanitizer(options.RequestSanitizeDefaults...)

	log.Filter(webutil.FlagHTTPRequest,
		stats.FilterNameSanitization,
		webutil.NewHTTPRequestEventFilter(func(_ context.Context, wre webutil.HTTPRequestEvent) (webutil.HTTPRequestEvent, bool) {
			wre.Request = requestSanitizer.Sanitize(wre.Request)
			return wre, false
		}),
	)

	log.Listen(webutil.FlagHTTPRequest, stats.ListenerNameStats,
		webutil.NewHTTPRequestEventListener(func(ctx context.Context, wre webutil.HTTPRequestEvent) {
			var route string
			if len(wre.Route) > 0 {
				route = stats.Tag(TagRoute, wre.Route)
			} else {
				route = stats.Tag(TagRoute, RouteNotFound)
			}

			proto := stats.Tag(TagProto, wre.Request.Proto)
			method := stats.Tag(TagMethod, wre.Request.Method)
			status := stats.Tag(TagStatus, strconv.Itoa(wre.StatusCode))
			tags := []string{
				proto, route, method, status,
			}
			tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)
			tags = append(tags, options.GetDefaultsAsTags()...)

			_ = collector.Increment(MetricNameHTTPRequest, tags...)
			_ = collector.Gauge(MetricNameHTTPRequestSize, float64(wre.ContentLength), tags...)
			_ = collector.Histogram(MetricNameHTTPRequestElapsed, timeutil.Milliseconds(wre.Elapsed), tags...)
			_ = collector.Gauge(MetricNameHTTPRequestElapsedLast, timeutil.Milliseconds(wre.Elapsed), tags...)
		}),
	)
}
