package dbstats

import (
	"context"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds db listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(db.QueryFlag, stats.ListenerNameStats, db.NewQueryEventListener(func(_ context.Context, qe db.QueryEvent) {
		engine := stats.Tag(TagEngine, qe.Engine)
		database := stats.Tag(TagDatabase, qe.Database)

		tags := []string{
			engine, database,
		}
		if len(qe.Label) > 0 {
			tags = append(tags, stats.Tag(TagQuery, qe.Label))
		}
		if qe.Err != nil {
			if ex := ex.As(qe.Err); ex != nil && ex.Class != nil {
				tags = append(tags, stats.Tag(stats.TagClass, ex.Class.Error()))
			}
			tags = append(tags, stats.TagError)
		}

		_ = collector.Increment(MetricNameDBQuery, tags...)
		_ = collector.Gauge(MetricNameDBQueryElapsed, timeutil.Milliseconds(qe.Elapsed), tags...)
		_ = collector.TimeInMilliseconds(MetricNameDBQueryElapsed, qe.Elapsed, tags...)
	}))
}
