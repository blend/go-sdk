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

func TestReturnAll(t *testing.T) {
	assert := assert.New(t)

	res := ReturnAll(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)
	assert.Len(res, 2)
}
