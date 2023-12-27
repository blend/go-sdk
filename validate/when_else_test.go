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

func TestWhenElse(t *testing.T) {
	assert := assert.New(t)

	var toggle bool
	when := WhenElse(func() bool { return toggle }, func() error { return fmt.Errorf("passes") }, func() error { return fmt.Errorf("fails") })

	err := when()
	assert.Equal(fmt.Errorf("fails"), err)

	toggle = true

	err = when()
	assert.Equal(fmt.Errorf("passes"), err)
}
