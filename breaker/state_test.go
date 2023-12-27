/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package breaker

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStateConstants(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(State(0), StateClosed)
	assert.Equal(State(1), StateHalfOpen)
	assert.Equal(State(2), StateOpen)

	assert.Equal(StateClosed.String(), "closed")
	assert.Equal(StateHalfOpen.String(), "half-open")
	assert.Equal(StateOpen.String(), "open")
	assert.Equal(State(100).String(), "unknown state: 100")
}
