/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCheckLabels(t *testing.T) {
	assert := assert.New(t)

	goodLabels := Labels{"foo": "bar", "foo.com/bar": "baz"}
	assert.Nil(CheckLabels(goodLabels))
	badLabels := Labels{"foo": "bar", "_foo.com/bar": "baz"}
	assert.NotNil(CheckLabels(badLabels))
}
