package stats

import (
	"context"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// AddQueryListeners adds db listeners.
func AddQueryListeners(log logger.Listenable, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	log.Listen(logger.Query, ListenerNameStats, db.NewQueryEventListener(func(_ context.Context, qe db.QueryEvent) {
		engine := Tag(TagEngine, qe.Engine)
		database := Tag(TagDatabase, qe.Database)

		tags := []string{
			engine, database,
		}
		if len(qe.QueryLabel) > 0 {
			tags = append(tags, Tag(TagQuery, qe.QueryLabel))
		}
		if qe.Err != nil {
			if ex := ex.As(qe.Err); ex != nil && ex.Class != nil {
				tags = append(tags, Tag(TagClass, ex.Class.Error()))
			}
			tags = append(tags, TagError)
		}

		stats.Increment(MetricNameDBQuery, tags...)
		stats.TimeInMilliseconds(MetricNameDBQueryElapsed, qe.Elapsed, tags...)
	}))
}
