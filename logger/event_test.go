/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type timestampProvider time.Time

func (tsp timestampProvider) GetFlag() string { return "timestamp_provider" }

func (tsp timestampProvider) GetTimestamp() time.Time {
	return time.Time(tsp)
}

var (
	_ Event             = (*timestampProvider)(nil)
	_ TimestampProvider = (*timestampProvider)(nil)
)

func TestGetEventTimestamp(t *testing.T) {
	assert := assert.New(t)

	ts1 := time.Date(2019, 8, 21, 12, 11, 10, 9, time.UTC)
	ts2 := time.Date(2019, 8, 20, 12, 11, 10, 9, time.UTC)
	ts3 := time.Date(2019, 8, 21, 12, 11, 10, 9, time.UTC)

	tsp := timestampProvider(ts1)

	tsc := WithTimestamp(context.Background(), ts2)
	ttsc := WithTriggerTimestamp(context.Background(), ts3)

	comboctx := WithTimestamp(WithTriggerTimestamp(context.Background(), ts3), ts2)
	comboctx2 := WithTriggerTimestamp(WithTimestamp(context.Background(), ts2), ts3)

	// test the timestamp provider takes precedence
	assert.Equal(ts1, GetEventTimestamp(context.Background(), tsp))
	assert.Equal(ts1, GetEventTimestamp(tsc, tsp))
	assert.Equal(ts1, GetEventTimestamp(ttsc, tsp))
	assert.Equal(ts1, GetEventTimestamp(comboctx, tsp))
	assert.Equal(ts1, GetEventTimestamp(comboctx2, tsp))

	me := NewMessageEvent(Info, "this is just a test")

	assert.False(GetEventTimestamp(context.Background(), me).IsZero())
	assert.Equal(ts2, GetEventTimestamp(tsc, me))
	assert.Equal(ts3, GetEventTimestamp(ttsc, me))
	assert.Equal(ts2, GetEventTimestamp(comboctx, me))
	assert.Equal(ts2, GetEventTimestamp(comboctx2, me))
}
