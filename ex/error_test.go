/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ex

import (
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestClassMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := Class("this is only a test")
	contents, marshalErr := json.Marshal(err)
	assert.Nil(marshalErr)
	assert.Equal(`"this is only a test"`, string(contents))
}
