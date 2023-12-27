/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.bar.local", OptTimeout(time.Second))
	assert.Equal(time.Second, r.Client.Timeout)
}
