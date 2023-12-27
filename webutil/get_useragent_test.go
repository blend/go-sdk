/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetUseragent(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("go-sdk test", GetUserAgent(NewMockRequest("GET", "/")))
}
