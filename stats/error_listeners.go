/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// AddErrorListeners adds error listeners.
func AddErrorListeners(log logger.Listenable, stats Collector, opts ...AddListenerOption) {
	if log == nil || stats == nil {
		return
	}

	options := NewAddListenerOptions(opts...)

	listener := logger.NewErrorEventListener(func(ctx context.Context, ee logger.ErrorEvent) {
		tags := []string{
			Tag(TagSeverity, string(ee.GetFlag())),
		}
		tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)
		_ = stats.Increment(MetricNameError,
			tags...,
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
func AddErrorListenersByClass(log logger.Listenable, stats Collector, opts ...AddListenerOption) {
	if log == nil || stats == nil {
		return
	}

	options := NewAddListenerOptions(opts...)

	listener := logger.NewErrorEventListener(func(ctx context.Context, ee logger.ErrorEvent) {
		tags := []string{
			Tag(TagSeverity, string(ee.GetFlag())),
			Tag(TagClass, fmt.Sprintf("%v", ex.ErrClass(ee.Err))),
		}
		tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)
		_ = stats.Increment(MetricNameError,
			tags...,
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}
