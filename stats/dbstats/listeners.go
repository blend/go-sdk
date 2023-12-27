/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package dbstats

import (
	"context"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds db listeners.
func AddListeners(log logger.Listenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	options := stats.NewAddListenerOptions(opts...)

	log.Listen(db.QueryFlag, stats.ListenerNameStats, db.NewQueryEventListener(func(ctx context.Context, qe db.QueryEvent) {
		engine := stats.Tag(TagEngine, qe.Engine)
		database := stats.Tag(TagDatabase, qe.Database)

		tags := []string{
			engine, database,
		}
		if len(qe.Label) > 0 {
			tags = append(tags, stats.Tag(TagQuery, qe.Label))
		}
		if qe.Err != nil {
			tags = append(tags, stats.TagError)
		}

		tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)

		_ = collector.Increment(MetricNameDBQuery, tags...)
		_ = collector.Gauge(MetricNameDBQueryElapsedLast, timeutil.Milliseconds(qe.Elapsed), tags...)
		_ = collector.Histogram(MetricNameDBQueryElapsed, timeutil.Milliseconds(qe.Elapsed), tags...)
	}))
}
