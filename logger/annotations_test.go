/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
