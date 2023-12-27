/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTraceConfig(t *testing.T) {
	a := assert.New(t)
	var empty SecretTraceConfig
	config := SecretTraceConfig{
		KeyName:	"A_KEY",
		VaultOperation:	"k1.put",
	}

	err := OptTraceConfig(config)(&empty)
	a.Nil(err)
	a.Equal("A_KEY", empty.KeyName)
	a.Equal("k1.put", empty.VaultOperation)
}

func TestOptTraceKeyName(t *testing.T) {
	a := assert.New(t)
	var ptr *SecretTraceConfig

	err := OptTraceKeyName("A_KEY")(ptr)
	a.NotNil(err)

	ptr = &SecretTraceConfig{}

	err = OptTraceKeyName("A_KEY")(ptr)
	a.Nil(err)
	a.Equal("A_KEY", ptr.KeyName)
}

func TestOptTraceVaultOperation(t *testing.T) {
	a := assert.New(t)
	var ptr *SecretTraceConfig

	err := OptTraceVaultOperation("k1.put")(ptr)
	a.NotNil(err)

	ptr = &SecretTraceConfig{}

	err = OptTraceVaultOperation("k1.put")(ptr)
	a.Nil(err)
	a.Equal("k1.put", ptr.VaultOperation)
}
