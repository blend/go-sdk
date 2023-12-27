/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLocalIP(t *testing.T) {
	assert := assert.New(t)

	assert.NotEmpty(LocalIP())
}
