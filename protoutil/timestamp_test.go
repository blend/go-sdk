/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Timestamp(t *testing.T) {
	its := assert.New(t)

	timestamp := time.Now().UTC()

	its.Equal(timestamp, FromTimestamp(Timestamp(timestamp)))

	// timestamp handles zero inputs
	its.Nil(Timestamp(time.Time{}))

	// from timestamp handles nil
	its.True(FromTimestamp(nil).IsZero())
}
