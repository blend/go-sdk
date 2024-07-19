/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSyncState(t *testing.T) {
	assert := assert.New(t)

	state := &SyncState{
		Values: map[string]interface{}{
			"foo":  "bar",
			"buzz": "fuzz",
		},
	}

	assert.Len(state.Keys(), 2)
	assert.Equal("bar", state.Get("foo"))
	assert.Equal("fuzz", state.Get("buzz"))

	state.Set("bar", "foo")
	assert.Equal("foo", state.Get("bar"))
	state.Remove("bar")
	assert.Nil(state.Get("bar"))
}
