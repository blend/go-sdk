/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSession(t *testing.T) {
	assert := assert.New(t)

	session := &Session{}
	session.WithBaseURL("https://foo.com/bar")
	assert.Equal("https://foo.com/bar", session.BaseURL)
	session.WithUserAgent("example-string")
	assert.Equal("example-string", session.UserAgent)
	session.WithRemoteAddr("10.10.32.1")
	assert.Equal("10.10.32.1", session.RemoteAddr)
}
