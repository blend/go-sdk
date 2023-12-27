/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

import "time"

// Collector is a stats collector.
type Collector interface {
	Taggable
	Count(name string, value int64, tags ...string) error
	Increment(name string, tags ...string) error
	Gauge(name string, value float64, tags ...string) error
	Histogram(name string, value float64, tags ...string) error
	Distribution(name string, value float64, tags ...string) error
	TimeInMilliseconds(name string, value time.Duration, tags ...string) error
	Flush() error
	Close() error
}
