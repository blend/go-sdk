/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package protoutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Duration(t *testing.T) {
	its := assert.New(t)

	its.Nil(Duration(0))
	its.Equal(500*time.Millisecond, FromDuration(Duration(500*time.Millisecond)))

	// from duration handles nil
	its.Zero(FromDuration(nil))
}
