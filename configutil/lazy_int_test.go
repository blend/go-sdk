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

func TestLazyInt(t *testing.T) {
	its := assert.New(t)

	isNil := LazyInt(nil)
	value := 0
	hasValue := LazyInt(&value)
	value2 := 2
	hasValue2 := LazyInt(&value2)

	var setValue int
	its.Nil(SetInt(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(2, setValue)
}
