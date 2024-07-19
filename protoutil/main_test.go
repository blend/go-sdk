/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import (
	"math"
	"time"

	"github.com/blend/go-sdk/protoutil/testdata"
	"github.com/blend/go-sdk/uuid"
)

// newTestMessage creates a new test message.
func newTestMessage() *testdata.Message {
	return &testdata.Message{
		Uid:           uuid.V4().String(),
		TimestampUtc:  Timestamp(time.Now().UTC()),
		Elapsed:       Duration(500 * time.Millisecond),
		StatusCode:    200,
		ContentLength: math.MaxInt32 + 1,
		Value:         3.14,
		Error:         "this is just a test",
	}
}
