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

// AddListenerOptions are options for adding listeners.
type AddListenerOptions struct {
	RequestSanitizeDefaults []sanitize.RequestOption
}

// AddListenerOption mutates AddListenerOptions
type AddListenerOption func(*AddListenerOptions)

// AddListeners adds web listeners.
func AddListeners(log logger.FilterListenable, collector stats.Collector, opts ...AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	var options AddListenerOptions
	for _, opt := range opts {
		opt(&options)
	}

	log.Filter(r2.Flag,
		stats.FilterNameSanitization,
		r2.NewEventFilter(func(_ context.Context, r2e r2.Event) (r2.Event, bool) {
			r2e.Request = sanitize.Request(r2e.Request, options.RequestSanitizeDefaults...)
			return r2e, false
		}),
	)
	log.Filter(r2.FlagResponse,
		stats.FilterNameSanitization,
		r2.NewEventFilter(func(_ context.Context, r2e r2.Event) (r2.Event, bool) {
			r2e.Request = sanitize.Request(r2e.Request, options.RequestSanitizeDefaults...)
			return r2e, false
		}),
	)

	log.Listen(r2.FlagResponse, stats.ListenerNameStats,
		r2.NewEventListener(func(_ context.Context, r2e r2.Event) {

			hostname := stats.Tag(TagHostname, r2e.Request.URL.Hostname())
			method := stats.Tag(TagMethod, r2e.Request.Method)
			status := stats.Tag(TagStatus, strconv.Itoa(r2e.Response.StatusCode))
			tags := []string{
				hostname, method, status,
			}

			_ = collector.Increment(MetricNameHTTPClientRequest, tags...)
			_ = collector.Gauge(MetricNameHTTPClientRequestElapsed, timeutil.Milliseconds(r2e.Elapsed), tags...)
			_ = collector.TimeInMilliseconds(MetricNameHTTPClientRequestElapsed, r2e.Elapsed, tags...)
		}),
	)
}
