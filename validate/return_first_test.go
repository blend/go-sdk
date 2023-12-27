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

func TestReturnFirst(t *testing.T) {
	assert := assert.New(t)

	res := ReturnFirst(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)
	assert.Equal(fmt.Errorf("one"), res)
}
