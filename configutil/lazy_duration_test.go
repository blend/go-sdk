/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestLazyDuration(t *testing.T) {
	its := assert.New(t)

	isNil := LazyDuration(nil)
	var value time.Duration = 0
	hasValue := LazyDuration(&value)
	var value2 time.Duration = 2
	hasValue2 := LazyDuration(&value2)

	var setValue time.Duration
	its.Nil(SetDuration(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(2, setValue)
}
