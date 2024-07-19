/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestErrorJSON(t *testing.T) {
	// assert that the error can be serialized as json.
	assert := assert.New(t)

	testErr := Error("this is only a test")

	contents, err := json.Marshal(testErr)
	assert.Nil(err)
	assert.Equal("\"this is only a test\"", string(contents))
}
