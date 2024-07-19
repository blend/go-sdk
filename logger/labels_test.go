/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombineLabels(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(CombineLabels(nil, nil, nil))
	combined := CombineLabels(Labels{"foo": "bar"}, nil, Labels{"moo": "loo"})
	assert.Equal("bar", combined["foo"])
	assert.Equal("loo", combined["moo"])
}
