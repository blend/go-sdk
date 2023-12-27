/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
