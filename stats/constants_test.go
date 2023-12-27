/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestVaultClientBackendKV(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("k:v", Tag("k", "v"))
}
