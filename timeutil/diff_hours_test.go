/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestDiffHours(t *testing.T) {
	assert := assert.New(t)

	t1 := time.Date(2017, 02, 27, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2017, 02, 24, 16, 0, 0, 0, time.UTC)
	t3 := time.Date(2017, 02, 28, 12, 0, 0, 0, time.UTC)

	assert.Equal(68, DiffHours(t2, t1))
	assert.Equal(24, DiffHours(t1, t3))
}
