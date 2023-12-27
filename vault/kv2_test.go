/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFixSecretDataPrefix(t *testing.T) {
	assert := assert.New(t)

	kv2 := &KV2{}

	assert.Equal(kv2.fixSecretDataPrefix("secret/foo/bar"), "secret/data/foo/bar")
	assert.Equal(kv2.fixSecretDataPrefix("secret/data/foo/bar"), "secret/data/foo/bar")
	assert.Equal(kv2.fixSecretDataPrefix("secret/datav2/foo/bar"), "secret/data/datav2/foo/bar")

	assert.Equal(kv2.fixSecretDataPrefix("secretv2/foo/bar"), "secretv2/foo/bar")
	assert.Equal(kv2.fixSecretDataPrefix("secretv2/datav2/foo/bar"), "secretv2/datav2/foo/bar")
}
