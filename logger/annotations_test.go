/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombineAnnotations(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(CombineAnnotations(nil, nil, nil))
	combined := CombineAnnotations(Annotations{"foo": "bar"}, nil, Annotations{"moo": "loo"})
	assert.Equal("bar", combined["foo"])
	assert.Equal("loo", combined["moo"])
}
