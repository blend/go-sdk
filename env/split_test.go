/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Split(t *testing.T) {
	its := assert.New(t)

	var key, value string
	key, value = Split("")
	its.Empty(key)
	its.Empty(value)

	key, value = Split("FOO")
	its.Empty(key)
	its.Empty(value)

	key, value = Split("FOO=bar")
	its.Equal("FOO", key)
	its.Equal("bar", value)

	key, value = Split("FOO=bar=baz")
	its.Equal("FOO", key)
	its.Equal("bar=baz", value)
}
