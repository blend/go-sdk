/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPortFromBindAddr(t *testing.T) {
	assert := assert.New(t)

	testCases := map[string]int32{
		"":			0,
		"2":			2,
		":2":			2,
		"127.0.0.1:1234":	1234,
		":8080":		8080,
	}

	for input, expected := range testCases {
		assert.Equal(expected, PortFromBindAddr(input))
	}
}
