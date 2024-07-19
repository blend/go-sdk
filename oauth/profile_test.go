/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestProfileUsername(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(Profile{}.Username())
	assert.Equal("foo", Profile{Email: "foo"}.Username())

	profile := Profile{
		Email: "test@blend.com",
	}

	assert.Equal("test@blend.com", profile.Username())

	profile = Profile{
		Email: "test2@blendlabs.com",
	}
	assert.Equal("test2@blendlabs.com", profile.Username())

	profile = Profile{
		Email: "test2+why@blendlabs.com",
	}
	assert.Equal("test2+why@blendlabs.com", profile.Username())
}
