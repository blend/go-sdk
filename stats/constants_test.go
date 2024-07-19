/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
