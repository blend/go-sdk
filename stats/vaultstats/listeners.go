package vaultstats

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/vault"
)

// AddListeners adds web listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(vault.Flag, stats.ListenerNameStats, vault.NewEventListener(func(_ context.Context, ve vault.Event) {
		tags := []string{
			stats.Tag("method", ve.Method),
			stats.Tag("status", strconv.Itoa(ve.StatusCode)),
			stats.Tag("path", ve.Path),
		}
		_ = collector.Increment("vault.request", tags...)
		_ = collector.Gauge("vault.request.elapsed", timeutil.Milliseconds(ve.Elapsed), tags...)
		_ = collector.TimeInMilliseconds("vault.request.elapsed", ve.Elapsed, tags...)
		_ = collector.Distribution("vault.request.elapsed", timeutil.Milliseconds(ve.Elapsed), tags...)
	}))
}
