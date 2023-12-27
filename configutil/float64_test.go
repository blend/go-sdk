/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFloat64(t *testing.T) {
	assert := assert.New(t)

	floatValue := Float64(0)
	ptr, err := floatValue.Float64(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	floatValue = Float64(3.14)
	ptr, err = floatValue.Float64(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(3.14, *ptr)
}
