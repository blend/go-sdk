/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cronstats

import (
	"context"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds web listeners.
func AddListeners(log logger.Listenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	flags := []string{
		cron.FlagBegin,
		cron.FlagComplete,
		cron.FlagCanceled,
		cron.FlagErrored,
		cron.FlagSuccess,
		cron.FlagBroken,
		cron.FlagFixed,
	}

	options := stats.NewAddListenerOptions(opts...)

	for _, flag := range flags {
		log.Listen(flag, stats.ListenerNameStats,
			cron.NewEventListener(func(ctx context.Context, ce cron.Event) {
				var tags []string
				tags = append(tags, stats.Tag(TagJob, ce.JobName))
				tags = append(tags, stats.Tag(TagJobStatus, ce.Flag))

				tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)

				_ = collector.Increment(MetricNameCron, tags...)
				if ce.Elapsed > 0 {
					_ = collector.Gauge(MetricNameCronElapsedLast, timeutil.Milliseconds(ce.Elapsed), tags...)
					_ = collector.Histogram(MetricNameCronElapsed, timeutil.Milliseconds(ce.Elapsed), tags...)
				}
			}),
		)
	}
}
