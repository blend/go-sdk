/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_Recover(t *testing.T) {
	assert := assert.New(t)

	errors := make(chan error, 1)
	Recover(func() error {
		return fmt.Errorf("test")
	}, errors)

	assert.NotEmpty(errors)
	assert.Equal(fmt.Errorf("test"), <-errors)

	errors = make(chan error, 1)
	Recover(func() error {
		panic("test")
	}, errors)

	assert.NotEmpty(errors)
	assert.Equal("test", ex.ErrClass(<-errors))
}
