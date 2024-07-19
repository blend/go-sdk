/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package autoflush

import "time"

// Defaults
const (
	DefaultMaxFlushes          = 128
	DefaultMaxLen              = 512
	DefaultFlushInterval       = 500 * time.Millisecond
	DefaultShutdownGracePeriod = 10 * time.Second
)

// Metric names
const (
	MetricFlush               string = "autoflush.flush"
	MetricFlushItemCount      string = "autoflush.flush.item_count"
	MetricFlushEnqueueElapsed string = "autoflush.flush.enqueue.elapsed"
	MetricFlushHandler        string = "autoflush.flush.handler"
	MetricFlushHandlerElapsed string = "autoflush.flush.handler.elapsed"
	MetricFlushQueueLength    string = "autoflush.flush.queue_length"
	MetricBufferLength        string = "autoflush.buffer.length"
	MetricAdd                 string = "autoflush.add"
	MetricAddElapsed          string = "autoflush.add.elapsed"
	MetricAddMany             string = "autoflush.add_many"
	MetricAddManyItemCount    string = "autoflush.add_many.item_count"
	MetricAddManyElapsed      string = "autoflush.add_many.elapsed"
)
