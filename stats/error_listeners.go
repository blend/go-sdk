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
		_ = stats.Increment(MetricNameError,
			Tag(TagSeverity, string(ee.GetFlag())),
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}

// AddErrorListenersByClass adds error listeners that add an exception class tag.
//
// NOTE: this will create many tag values if you do not use exceptions correctly,
// that is, if you put variable data in the exception class.
// If there is any doubt which of these to use (AddErrorListeners or AddErrorListenersByClass)
// use the version that does not add the class information (AddErrorListeners).
func AddErrorListenersByClass(log logger.Listenable, stats Collector) {
	if log == nil || stats == nil {
		return
	}

	listener := logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		_ = stats.Increment(MetricNameError,
			Tag(TagSeverity, string(ee.GetFlag())),
			Tag(TagClass, fmt.Sprintf("%v", ex.ErrClass(ee.Err))),
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}
