package cronstats

import (
	"context"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(cron.FlagComplete, stats.ListenerNameStats,
		cron.NewEventListener(func(_ context.Context, ce cron.Event) {
			var tags []string
			tags = append(tags, stats.Tag(TagJob, ce.JobName))
			tags = append(tags, stats.Tag(TagJobStatus, ce.Flag))

			_ = collector.Increment(MetricNameCron, tags...)
			_ = collector.Gauge(MetricNameCronElapsed, timeutil.Milliseconds(ce.Elapsed), tags...)
			_ = collector.TimeInMilliseconds(MetricNameCronElapsed, ce.Elapsed, tags...)
		}),
	)
}
