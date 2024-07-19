/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2stats

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/sanitize"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.FilterListenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	options := stats.NewAddListenerOptions(opts...)

	requestSanitizer := sanitize.NewRequestSanitizer(options.RequestSanitizeDefaults...)

	log.Filter(r2.Flag,
		stats.FilterNameSanitization,
		r2.NewEventFilter(func(_ context.Context, r2e r2.Event) (r2.Event, bool) {
			r2e.Request = requestSanitizer.Sanitize(r2e.Request)
			return r2e, false
		}),
	)
	log.Filter(r2.FlagResponse,
		stats.FilterNameSanitization,
		r2.NewEventFilter(func(_ context.Context, r2e r2.Event) (r2.Event, bool) {
			r2e.Request = requestSanitizer.Sanitize(r2e.Request)
			return r2e, false
		}),
	)

	log.Listen(r2.FlagResponse, stats.ListenerNameStats,
		r2.NewEventListener(func(ctx context.Context, r2e r2.Event) {
			hostname := stats.Tag(TagHostname, r2e.Request.URL.Hostname())
			target := stats.Tag(TagTarget, r2e.Request.URL.Hostname())
			method := stats.Tag(TagMethod, r2e.Request.Method)
			status := stats.Tag(TagStatus, strconv.Itoa(r2e.Response.StatusCode))
			tags := []string{
				hostname, target, method, status,
			}
			tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)
			_ = collector.Increment(MetricNameHTTPClientRequest, tags...)
			_ = collector.Gauge(MetricNameHTTPClientRequestElapsedLast, timeutil.Milliseconds(r2e.Elapsed), tags...)
			_ = collector.Histogram(MetricNameHTTPClientRequestElapsed, timeutil.Milliseconds(r2e.Elapsed), tags...)
		}),
	)
}
