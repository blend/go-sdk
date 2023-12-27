/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
func AddListeners(log logger.Listenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	options := stats.NewAddListenerOptions(opts...)

	log.Listen(vault.Flag, stats.ListenerNameStats, vault.NewEventListener(func(ctx context.Context, ve vault.Event) {
		tags := []string{
			stats.Tag("method", ve.Method),
			stats.Tag("status", strconv.Itoa(ve.StatusCode)),
			stats.Tag("path", ve.Path),
		}
		tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)
		_ = collector.Increment("vault.request", tags...)
		_ = collector.TimeInMilliseconds("vault.request.elapsed", ve.Elapsed, tags...)
		_ = collector.Histogram("vault.request.elapsed", timeutil.Milliseconds(ve.Elapsed), tags...)
	}))
}
