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

func TestLazyFloat64(t *testing.T) {
	its := assert.New(t)

	isNil := LazyFloat64(nil)
	var value float64 = 0
	hasValue := LazyFloat64(&value)
	var value2 float64 = 2
	hasValue2 := LazyFloat64(&value2)

	var setValue float64
	its.Nil(SetFloat64(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(2, setValue)
}
