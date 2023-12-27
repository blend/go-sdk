/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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

func Test_Duration(t *testing.T) {
	its := assert.New(t)

	its.Nil(Duration(0))
	its.Equal(500*time.Millisecond, FromDuration(Duration(500*time.Millisecond)))

	// from duration handles nil
	its.Zero(FromDuration(nil))
}
