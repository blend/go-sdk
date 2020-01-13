package stats

import (
	"context"
	"fmt"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// AddErrorListeners adds error listeners.
func AddErrorListeners(log logger.Listenable, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	listener := logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		stats.Increment(MetricNameError,
			Tag(TagSeverity, string(ee.GetFlag())),
			Tag(TagClass, fmt.Sprintf("%v", ex.ErrClass(ee.Err))),
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}
