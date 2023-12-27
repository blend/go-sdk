/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
