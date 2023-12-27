/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWhen(t *testing.T) {
	assert := assert.New(t)

	var toggle bool
	when := When(func() bool { return toggle }, func() error { return fmt.Errorf("passes") })
	assert.Nil(when())
	toggle = true
	assert.Equal(fmt.Errorf("passes"), when())
}
